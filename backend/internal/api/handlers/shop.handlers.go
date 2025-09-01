package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/db/errors"
	"go.uber.org/zap"

	"github.com/petrejonn/naytife/internal/observability"
)

type ShopStatus string

const (
	DRAFT     ShopStatus = "DRAFT"
	PUBLISHED ShopStatus = "PUBLISHED"
	ARCHIVED  ShopStatus = "ARCHIVED"
	SUSPENDED ShopStatus = "SUSPENDED"
)

// CreateShop creates a shop
// @Summary      Create a shop
// @Description
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param        shop body models.ShopCreateParams true "Shop object that needs to be created"
// @Success      200  {object}   models.SuccessResponse{data=models.Shop} "Shop created successfully"
// @Failure      400  {object}   models.ErrorResponse "Invalid request body"
// @Failure      409  {object}   models.ErrorResponse "Conflict"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops [post]
func (h *Handler) CreateShop(c *fiber.Ctx) error {
	userSub, _ := c.Locals("user_id").(string)

	// Parse input
	var shop models.ShopCreateParams
	if err := c.BodyParser(&shop); err != nil {
		zap.L().Error("CreateShop: failed to parse request body", zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Validate user
	user, err := h.Repository.GetUserBySub(c.Context(), &userSub)
	if err != nil {
		zap.L().Warn("CreateShop: failed to get user profile", zap.String("user_sub", userSub), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusUnauthorized, "Failed to get profile", nil)
	}

	// Defaults
	shop.Status = string(db.ProductStatusDRAFT)
	shop.CurrencyCode = "NGN"
	shop.Status = string(PUBLISHED)

	// Validate slug + struct
	if !slug.IsSlug(shop.Subdomain) {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid subdomain format", nil)
	}
	validator := &models.XValidator{}
	if errs := validator.Validate(&shop); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: errMsgs}
	}

	// Create shop in DB
	param := db.CreateShopParams{
		OwnerID:         user.UserID,
		Title:           shop.Title,
		Subdomain:       shop.Subdomain,
		CurrencyCode:    shop.CurrencyCode,
		Status:          shop.Status,
		CurrentTemplate: &shop.Template,
	}
	objDB, err := h.Repository.CreateShop(c.Context(), param)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == errors.UniqueViolation {
			zap.L().Info("CreateShop: shop already exists", zap.String("subdomain", shop.Subdomain), zap.String("user_id", user.UserID.String()))
			return api.ErrorResponse(c, fiber.StatusConflict, "Shop already exists", nil)
		}
		zap.L().Error("CreateShop: failed to create shop", zap.Error(err), zap.String("subdomain", shop.Subdomain))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create shop", nil)
	}

	// Response
	resp := models.Shop{
		ID:                  objDB.ShopID,
		Title:               objDB.Title,
		Subdomain:           objDB.Subdomain,
		CurrencyCode:        objDB.CurrencyCode,
		Status:              objDB.Status,
		CreatedAt:           objDB.CreatedAt,
		UpdatedAt:           objDB.UpdatedAt,
		Email:               objDB.Email,
		About:               objDB.About,
		Address:             objDB.Address,
		PhoneNumber:         objDB.PhoneNumber,
		WhatsappPhoneNumber: objDB.WhatsappPhoneNumber,
		WhatsappLink:        objDB.WhatsappLink,
		FacebookLink:        objDB.FacebookLink,
		InstagramLink:       objDB.InstagramLink,
		SeoDescription:      objDB.SeoDescription,
		SeoKeywords:         objDB.SeoKeywords,
		SeoTitle:            objDB.SeoTitle,
		CurrentTemplate:     objDB.CurrentTemplate,
	}

	// Auto-trigger deployment for new shops (DB record + StoreDeployerClient)
	go func(shopID int64, subdomain, templateName string) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		ctx, finish := observability.StartSpan(ctx, "autoDeployNewShop", "store-deployer", "POST", "deploy")
		defer finish(0, nil)

		// 1) Create deployment record (deploying)
		startedAt := pgtype.Timestamptz{Time: time.Now(), Valid: true}
		deployment, derr := h.Repository.CreateDeployment(ctx, db.CreateDeploymentParams{
			ShopID:          shopID,
			TemplateName:    templateName,
			TemplateVersion: "latest",
			Status:          "deploying",
			DeploymentType:  "full",
			Message:         nil,
			StartedAt:       startedAt,
		})
		if derr != nil {
			zap.L().Error("autoDeployNewShop: failed to create deployment record",
				zap.Int64("shop_id", shopID),
				zap.String("subdomain", subdomain),
				zap.String("template", templateName),
				zap.Error(derr))
			return
		}

		// 2) Call store-deployer via the client
		if err := h.StoreDeployerClient.Deploy(ctx, shopID, subdomain, templateName); err != nil {
			errMsg := err.Error()
			_ = h.Repository.UpdateDeploymentStatus(ctx, db.UpdateDeploymentStatusParams{
				DeploymentID: deployment.DeploymentID,
				Status:       "failed",
				Message:      &errMsg,
			})

			zap.L().Warn("autoDeployNewShop: auto-deploy failed",
				zap.Int64("shop_id", shopID),
				zap.String("subdomain", subdomain),
				zap.String("template", templateName),
				zap.Int64("deployment_id", deployment.DeploymentID),
				zap.Error(err))
			return
		}

		// Optional: if you have a distinct "queued"/"requested" state, set it here.
		// If you want to keep legacy semantics (leave "deploying"), do nothing.
		// _ = h.Repository.UpdateDeploymentStatus(ctx, db.UpdateDeploymentStatusParams{
		// 	DeploymentID: deployment.DeploymentID,
		// 	Status:       "deploying",
		// 	Message:      nil,
		// })

		zap.L().Info("autoDeployNewShop: auto-deployment triggered",
			zap.Int64("shop_id", shopID),
			zap.String("subdomain", subdomain),
			zap.String("template", templateName))
	}(objDB.ShopID, objDB.Subdomain, shop.Template)

	zap.L().Info("CreateShop: shop created successfully",
		zap.Int64("shop_id", objDB.ShopID),
		zap.String("subdomain", objDB.Subdomain),
		zap.String("owner_id", user.UserID.String()))

	return api.SuccessResponse(c, fiber.StatusCreated, resp, "Shop created")
}

// GetShops fetches auth user shop
// @Summary      Fetch all shops
// @Description
// @Tags         shops
// @Produce      json
// @Success      200  {object}  models.SuccessResponse{data=[]models.Shop}  "Shops fetched successfully"
// @Failure      401  {object}   models.ErrorResponse "Unauthorized"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops [get]
func (h *Handler) GetShops(c *fiber.Ctx) error {
	userSub, _ := c.Locals("user_id").(string)

	var user db.User
	var err error
	user, err = h.Repository.GetUserBySub(c.Context(), &userSub)
	if err != nil {
		zap.L().Warn("GetShops: failed to get user profile", zap.String("user_sub", userSub), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusUnauthorized, "Failed to get profile", nil)
	}
	objsDB, err := h.Repository.GetShopsByOwner(c.Context(), user.UserID)
	if err != nil {
		zap.L().Error("GetShops: failed to get shops for owner", zap.Error(err), zap.String("user_id", user.UserID.String()))
		api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get shops", nil)
	}
	var shops []models.Shop
	for _, obj := range objsDB {
		shop := models.Shop{
			ID:                  obj.ShopID,
			Title:               obj.Title,
			Subdomain:           obj.Subdomain,
			CurrencyCode:        obj.CurrencyCode,
			Status:              obj.Status,
			CreatedAt:           obj.CreatedAt,
			UpdatedAt:           obj.UpdatedAt,
			Email:               obj.Email,
			About:               obj.About,
			Address:             obj.Address,
			PhoneNumber:         obj.PhoneNumber,
			WhatsappPhoneNumber: obj.WhatsappPhoneNumber,
			WhatsappLink:        obj.WhatsappLink,
			FacebookLink:        obj.FacebookLink,
			InstagramLink:       obj.InstagramLink,
			SeoDescription:      obj.SeoDescription,
			SeoKeywords:         obj.SeoKeywords,
			SeoTitle:            obj.SeoTitle,
			CurrentTemplate:     obj.CurrentTemplate,
		}

		// Fetch shop images
		shopImages, imgErr := h.Repository.GetShopImages(c.Context(), obj.ShopID)
		if imgErr == nil {
			// Only add images if successfully retrieved
			shop.Images = &models.ShopImagesData{
				ID:                shopImages.ShopImageID,
				FaviconUrl:        shopImages.FaviconUrl,
				LogoUrl:           shopImages.LogoUrl,
				LogoUrlDark:       shopImages.LogoUrlDark,
				BannerUrl:         shopImages.BannerUrl,
				BannerUrlDark:     shopImages.BannerUrlDark,
				CoverImageUrl:     shopImages.CoverImageUrl,
				CoverImageUrlDark: shopImages.CoverImageUrlDark,
			}
		}

		shops = append(shops, shop)
	}

	return api.SuccessResponse(c, fiber.StatusOK, shops, "Shops retrieved")
}

// DeleteShop deletes a shop
// @Summary      Delete a shop
// @Description
// @Tags         shops
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Success      200  {object}   models.SuccessResponse "Shop deleted successfully"
// @Failure      400  {object}   models.ErrorResponse "Invalid request body"
// @Failure      401  {object}   models.ErrorResponse "Unauthorized"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id} [delete]
func (h *Handler) DeleteShop(c *fiber.Ctx) error {
	shopID := c.Params("shop_id", "0")
	shopIDInt, _ := strconv.ParseInt(shopID, 10, 64)

	// Fetch shop info before deletion (needed for cleanup)
	shop, err := h.Repository.GetShop(c.Context(), shopIDInt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
		}
		zap.L().Error("DeleteShop: failed to retrieve shop before deletion",
			zap.Int64("shop_id", shopIDInt),
			zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve shop", nil)
	}

	// Trigger cleanup asynchronously in store-deployer
	go func(shopID int64, subdomain string) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		ctx, finish := observability.StartSpan(ctx, "cleanupStoreFiles", "store-deployer", "DELETE", "cleanup")
		defer finish(0, nil)

		if err := h.StoreDeployerClient.Cleanup(ctx, subdomain, shopID); err != nil {
			zap.L().Warn("DeleteShop: failed to cleanup store files",
				zap.Int64("shop_id", shopID),
				zap.String("subdomain", subdomain),
				zap.Error(err))
		}
	}(shopIDInt, shop.Subdomain)

	// Proceed with database deletion
	err = h.Repository.DeleteShop(c.Context(), shopIDInt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
		}
		zap.L().Error("DeleteShop: failed to delete shop",
			zap.Int64("shop_id", shopIDInt),
			zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete shop", nil)
	}

	return api.SuccessResponse(c, fiber.StatusOK, nil, fmt.Sprintf("Deleted Shop {%s}", shopID))
}

// GetShop fetches a shop
// @Summary      Fetch a shop
// @Description
// @Tags         shops
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Success      200  {object}   models.SuccessResponse{data=models.Shop} "Shop fetched successfully"
// @Failure      404  {object}   models.ErrorResponse "Shop not found"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id} [get]
func (h *Handler) GetShop(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	objDB, err := h.Repository.GetShop(c.Context(), shopID)
	if err != nil {
		if err == pgx.ErrNoRows {
			zap.L().Warn("GetShop: shop not found", zap.Int64("shop_id", shopID))
			return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
		}
		zap.L().Error("GetShop: failed to fetch shop", zap.Int64("shop_id", shopID), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch shop", nil)
	}
	resp := models.Shop{
		ID:                  objDB.ShopID,
		Title:               objDB.Title,
		Subdomain:           objDB.Subdomain,
		CurrencyCode:        objDB.CurrencyCode,
		Status:              objDB.Status,
		CreatedAt:           objDB.CreatedAt,
		UpdatedAt:           objDB.UpdatedAt,
		Email:               objDB.Email,
		About:               objDB.About,
		Address:             objDB.Address,
		PhoneNumber:         objDB.PhoneNumber,
		WhatsappPhoneNumber: objDB.WhatsappPhoneNumber,
		WhatsappLink:        objDB.WhatsappLink,
		FacebookLink:        objDB.FacebookLink,
		InstagramLink:       objDB.InstagramLink,
		SeoDescription:      objDB.SeoDescription,
		SeoKeywords:         objDB.SeoKeywords,
		SeoTitle:            objDB.SeoTitle,
		CurrentTemplate:     objDB.CurrentTemplate,
	}

	// Fetch shop images
	shopImages, err := h.Repository.GetShopImages(c.Context(), shopID)
	if err == nil {
		// Only add images if successfully retrieved
		resp.Images = &models.ShopImagesData{
			ID:                shopImages.ShopImageID,
			FaviconUrl:        shopImages.FaviconUrl,
			LogoUrl:           shopImages.LogoUrl,
			LogoUrlDark:       shopImages.LogoUrlDark,
			BannerUrl:         shopImages.BannerUrl,
			BannerUrlDark:     shopImages.BannerUrlDark,
			CoverImageUrl:     shopImages.CoverImageUrl,
			CoverImageUrlDark: shopImages.CoverImageUrlDark,
		}
	}

	return api.SuccessResponse(c, fiber.StatusOK, resp, "Shop retrieved")
}

// GetShopBySubDomain fetches a shop by subdomain
// @Summary      Fetch a shop by subdomain
// @Description
// @Tags         shops
// @Produce      json
// @Param        subdomain path string true "Shop Subdomain"
// @Success      200  {object}   models.SuccessResponse{data=models.Shop} "Shop fetched successfully"
// @Failure      404  {object}   models.ErrorResponse "Shop not found"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /subdomains/{subdomain} [get]
func (h *Handler) GetShopBySubDomain(c *fiber.Ctx) error {
	subdomain := c.Params("subdomain", "")
	if subdomain == "" {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Subdomain is required", nil)
	}

	objDB, err := h.Repository.GetShopBySubDomain(c.Context(), subdomain)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch shop", nil)
	}

	resp := models.Shop{
		ID:                  objDB.ShopID,
		Title:               objDB.Title,
		Subdomain:           objDB.Subdomain,
		CurrencyCode:        objDB.CurrencyCode,
		Status:              objDB.Status,
		CreatedAt:           objDB.CreatedAt,
		UpdatedAt:           objDB.UpdatedAt,
		Email:               objDB.Email,
		About:               objDB.About,
		Address:             objDB.Address,
		PhoneNumber:         objDB.PhoneNumber,
		WhatsappPhoneNumber: objDB.WhatsappPhoneNumber,
		WhatsappLink:        objDB.WhatsappLink,
		FacebookLink:        objDB.FacebookLink,
		InstagramLink:       objDB.InstagramLink,
		SeoDescription:      objDB.SeoDescription,
		SeoKeywords:         objDB.SeoKeywords,
		SeoTitle:            objDB.SeoTitle,
		CurrentTemplate:     objDB.CurrentTemplate,
	}

	// Fetch shop images
	shopImages, err := h.Repository.GetShopImages(c.Context(), objDB.ShopID)
	if err == nil {
		// Only add images if successfully retrieved
		resp.Images = &models.ShopImagesData{
			ID:                shopImages.ShopImageID,
			FaviconUrl:        shopImages.FaviconUrl,
			LogoUrl:           shopImages.LogoUrl,
			LogoUrlDark:       shopImages.LogoUrlDark,
			BannerUrl:         shopImages.BannerUrl,
			BannerUrlDark:     shopImages.BannerUrlDark,
			CoverImageUrl:     shopImages.CoverImageUrl,
			CoverImageUrlDark: shopImages.CoverImageUrlDark,
		}
	}

	return api.SuccessResponse(c, fiber.StatusOK, resp, "Shop retrieved")
}

// CheckSubdomainAvailability checks if a subdomain is available
// @Summary      Check subdomain availability
// @Description  Check if a subdomain is available for creating a new shop
// @Tags         shops
// @Produce      json
// @Param        subdomain path string true "Subdomain to check"
// @Success      200  {object}   models.SuccessResponse{data=models.SubdomainAvailabilityResponse} "Subdomain availability checked"
// @Failure      400  {object}   models.ErrorResponse "Invalid subdomain format"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /subdomains/{subdomain}/check [get]
func (h *Handler) CheckSubdomainAvailability(c *fiber.Ctx) error {
	subdomain := c.Params("subdomain", "")
	if subdomain == "" {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Subdomain is required", nil)
	}

	// Validate subdomain format
	if !slug.IsSlug(subdomain) {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid subdomain format", nil)
	}

	// Check if subdomain already exists
	_, err := h.Repository.GetShopBySubDomain(c.Context(), subdomain)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Subdomain is available
			response := models.SubdomainAvailabilityResponse{
				Subdomain: subdomain,
				Available: true,
				Message:   "Subdomain is available",
			}
			return api.SuccessResponse(c, fiber.StatusOK, response, "Subdomain is available")
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check subdomain availability", nil)
	}

	// Subdomain already exists
	response := models.SubdomainAvailabilityResponse{
		Subdomain: subdomain,
		Available: false,
		Message:   "Subdomain is already taken",
	}
	return api.SuccessResponse(c, fiber.StatusOK, response, "Subdomain is not available")
}

// UpdateShop updates a shop
// @Summary      Update a shop
// @Description
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        shop body models.ShopUpdateParams true "Shop update parameters"
// @Success      200  {object}   models.SuccessResponse{data=models.Shop} "Shop updated successfully"
// @Failure      400  {object}   models.ErrorResponse "Invalid request body"
// @Failure      404  {object}   models.ErrorResponse "Shop not found"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id} [put]
func (h *Handler) UpdateShop(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)

	var shop models.ShopUpdateParams
	if err := c.BodyParser(&shop); err != nil {
		h.Logger.Warn("invalid request body", zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validator := &models.XValidator{}
	if errs := validator.Validate(&shop); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		h.Logger.Warn("validation failed", zap.String("errors", errMsgs))
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: errMsgs}
	}

	param := db.UpdateShopParams{
		ShopID:              shopID,
		Title:               shop.Title,
		CurrencyCode:        shop.CurrencyCode,
		About:               shop.About,
		Status:              shop.Status,
		PhoneNumber:         shop.PhoneNumber,
		WhatsappLink:        shop.WhatsappLink,
		WhatsappPhoneNumber: shop.WhatsappPhoneNumber,
		FacebookLink:        shop.FacebookLink,
		InstagramLink:       shop.InstagramLink,
		SeoDescription:      shop.SeoDescription,
		SeoKeywords:         shop.SeoKeywords,
		SeoTitle:            shop.SeoTitle,
		Address:             shop.Address,
		Email:               shop.Email,
	}

	objDB, err := h.Repository.UpdateShop(c.Context(), param)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
		}
		h.Logger.Error("failed to update shop",
			zap.Int64("shop_id", shopID),
			zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update shop", nil)
	}

	resp := models.Shop{
		ID:                  objDB.ShopID,
		Title:               objDB.Title,
		Subdomain:           objDB.Subdomain,
		CurrencyCode:        objDB.CurrencyCode,
		Status:              objDB.Status,
		CreatedAt:           objDB.CreatedAt,
		UpdatedAt:           objDB.UpdatedAt,
		Email:               objDB.Email,
		About:               objDB.About,
		Address:             objDB.Address,
		PhoneNumber:         objDB.PhoneNumber,
		WhatsappPhoneNumber: objDB.WhatsappPhoneNumber,
		WhatsappLink:        objDB.WhatsappLink,
		FacebookLink:        objDB.FacebookLink,
		InstagramLink:       objDB.InstagramLink,
		SeoDescription:      objDB.SeoDescription,
		SeoKeywords:         objDB.SeoKeywords,
		SeoTitle:            objDB.SeoTitle,
		CurrentTemplate:     objDB.CurrentTemplate,
	}

	// Auto-publish asynchronously via StoreDeployerClient
	go func(shopID int64, subdomain string) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		ctx, finish := observability.StartSpan(ctx, "autoPublishShopUpdate", "store-deployer", "POST", "update-data")
		defer finish(0, nil)

		if err := h.StoreDeployerClient.UpdateData(ctx, subdomain, shopID, "shop"); err != nil {
			h.Logger.Warn("auto-publish failed",
				zap.Int64("shop_id", shopID),
				zap.Error(err))
		}
	}(shopID, objDB.Subdomain)

	return api.SuccessResponse(c, fiber.StatusOK, resp, "Shop updated successfully")
}

// UpdateShopImages updates the shop images
// @Summary      Update shop images
// @Description  Update shop images (logo, favicon, banner, etc)
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        images body models.ShopImagesUpdateParams true "Shop images object"
// @Success      200  {object}   models.SuccessResponse{data=models.ShopImagesData} "Shop images updated successfully"
// @Failure      400  {object}   models.ErrorResponse "Bad request"
// @Failure      404  {object}   models.ErrorResponse "Shop not found"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/images [put]
func (h *Handler) UpdateShopImages(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		h.Logger.Warn("invalid shop id", zap.String("shop_id", shopIDStr), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	// Verify shop exists
	shop, err := h.Repository.GetShop(c.Context(), shopID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
		}
		h.Logger.Error("failed to get shop", zap.Int64("shop_id", shopID), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to verify shop", nil)
	}

	// Parse body
	var imageParams models.ShopImagesUpdateParams
	if err := c.BodyParser(&imageParams); err != nil {
		h.Logger.Warn("invalid request body", zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Update or create shop images
	var updatedImages db.ShopImage
	err = h.Repository.WithTx(c.Context(), func(q *db.Queries) error {
		_, err := q.GetShopImages(c.Context(), shopID)
		if err != nil {
			if err == pgx.ErrNoRows {
				// Create new record
				createdImgs, err := q.CreateShopImages(c.Context(), db.CreateShopImagesParams{
					ShopID:            shopID,
					FaviconUrl:        imageParams.FaviconUrl,
					LogoUrl:           imageParams.LogoUrl,
					LogoUrlDark:       imageParams.LogoUrlDark,
					BannerUrl:         imageParams.BannerUrl,
					BannerUrlDark:     imageParams.BannerUrlDark,
					CoverImageUrl:     imageParams.CoverImageUrl,
					CoverImageUrlDark: imageParams.CoverImageUrlDark,
				})
				if err != nil {
					return err
				}
				updatedImages = createdImgs
				return nil
			}
			return err
		}

		// Update existing record
		updatedImgs, err := q.UpdateShopImages(c.Context(), db.UpdateShopImagesParams{
			ShopID:            shopID,
			FaviconUrl:        imageParams.FaviconUrl,
			LogoUrl:           imageParams.LogoUrl,
			LogoUrlDark:       imageParams.LogoUrlDark,
			BannerUrl:         imageParams.BannerUrl,
			BannerUrlDark:     imageParams.BannerUrlDark,
			CoverImageUrl:     imageParams.CoverImageUrl,
			CoverImageUrlDark: imageParams.CoverImageUrlDark,
		})
		if err != nil {
			return err
		}
		updatedImages = updatedImgs
		return nil
	})
	if err != nil {
		h.Logger.Error("failed to update shop images", zap.Int64("shop_id", shopID), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update shop images", nil)
	}

	response := models.ShopImagesData{
		ID:                updatedImages.ShopImageID,
		FaviconUrl:        updatedImages.FaviconUrl,
		LogoUrl:           updatedImages.LogoUrl,
		LogoUrlDark:       updatedImages.LogoUrlDark,
		BannerUrl:         updatedImages.BannerUrl,
		BannerUrlDark:     updatedImages.BannerUrlDark,
		CoverImageUrl:     updatedImages.CoverImageUrl,
		CoverImageUrlDark: updatedImages.CoverImageUrlDark,
	}

	// Auto-publish asynchronously
	go func(shopID int64, subdomain string) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		ctx, finish := observability.StartSpan(ctx, "autoPublishShopImages", "store-deployer", "POST", "update-data")
		defer finish(0, nil)

		if err := h.StoreDeployerClient.UpdateData(ctx, subdomain, shopID, "shop"); err != nil {
			h.Logger.Warn("auto-publish failed",
				zap.Int64("shop_id", shopID),
				zap.Error(err))
		}
	}(shopID, shop.Subdomain)

	return api.SuccessResponse(c, fiber.StatusOK, response, "Shop images updated successfully")
}
