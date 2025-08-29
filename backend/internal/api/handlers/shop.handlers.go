package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/db/errors"
	"go.uber.org/zap"

	ic "github.com/petrejonn/naytife/internal/httpclient"

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
		h.Logger.Error("failed to parse request body", zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Validate user
	user, err := h.Repository.GetUserBySub(c.Context(), &userSub)
	if err != nil {
		h.Logger.Warn("failed to get user profile", zap.String("user_sub", userSub), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusUnauthorized, "Failed to get profile", nil)
	}

	// Default values
	shop.Status = string(db.ProductStatusDRAFT)
	shop.CurrencyCode = "NGN"

	// Validate slug
	if !slug.IsSlug(shop.Subdomain) {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid subdomain format", nil)
	}

	// Validate struct fields
	validator := &models.XValidator{}
	if errs := validator.Validate(&shop); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
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
			h.Logger.Info("shop already exists",
				zap.String("subdomain", shop.Subdomain),
				zap.String("user_id", user.UserID.String()))
			return api.ErrorResponse(c, fiber.StatusConflict, "Shop already exists", nil)
		}
		h.Logger.Error("failed to create shop", zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create shop", nil)
	}

	// Prepare response
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

	// Auto-trigger deployment for new shops
	go func(shopID int64, subdomain, templateName string) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		ctx, finish := observability.StartSpan(ctx, "autoDeployNewShop", "store-deployer", "POST", "deploy")
		defer finish(0, nil)

		h.autoDeployNewShopWithCtx(ctx, shopID, subdomain, templateName)
	}(objDB.ShopID, objDB.Subdomain, shop.Template)

	h.Logger.Info("shop created successfully",
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

	// Development mode bypass
	devMode := os.Getenv("DEV_MODE")
	if devMode == "true" {
		// Use the test user we created
		testEmail := "test@example.com"
		user, err = h.Repository.GetUser(c.Context(), &testEmail)
	} else {
		user, err = h.Repository.GetUserBySub(c.Context(), &userSub)
	}

	if err != nil {
		return api.ErrorResponse(c, fiber.StatusUnauthorized, "Failed to get profile", nil)
	}
	objsDB, err := h.Repository.GetShopsByOwner(c.Context(), user.UserID)
	if err != nil {
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

	// First, get the shop data before deletion to extract subdomain for cleanup
	shop, err := h.Repository.GetShop(c.Context(), shopIDInt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve shop", nil)
	}

	// Call store-deployer cleanup endpoint to remove R2 files
	// This is done before database deletion to ensure we have the shop data
	// Pass request context so cancellation and trace propagation work correctly
	err = h.cleanupStoreFiles(c.Context(), shopIDInt, shop.Subdomain)
	if err != nil {
		// Log the error but don't fail the deletion - database cleanup should proceed
		fmt.Printf("Warning: Failed to cleanup R2 files for shop %d: %v\n", shopIDInt, err)
	}

	// Proceed with database deletion
	err = h.Repository.DeleteShop(c.Context(), shopIDInt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
		}
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
			return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
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
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validator := &models.XValidator{}
	if errs := validator.Validate(&shop); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
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

	// Auto-publish if publish handler is available
	// Use cancellable worker context for auto publish
	go func(shopID int64, changeType, entity, description string) {
		// derive worker context from request context
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		ctx, finish := observability.StartSpan(ctx, "autoPublishShopChanges", "store-deployer", "POST", "update-data")
		defer finish(0, nil)
		h.autoPublishShopChangesWithCtx(ctx, shopID, changeType, entity, description)
	}(shopID, "shop_update", fmt.Sprintf("shop:%d", shopID), "Shop details updated")

	return api.SuccessResponse(c, fiber.StatusOK, resp, "Shop updated")
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
	// Parse shop ID from params
	shopIDStr := c.Params("shop_id", "0")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	// Verify shop exists
	_, err = h.Repository.GetShop(c.Context(), shopID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to verify shop", nil)
	}

	// Parse request body
	var imageParams models.ShopImagesUpdateParams
	if err := c.BodyParser(&imageParams); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Check if shop images record exists
	var updatedImages db.ShopImage
	var exists bool

	err = h.Repository.WithTx(c.Context(), func(q *db.Queries) error {
		// Try to get existing shop images
		_, err := q.GetShopImages(c.Context(), shopID)
		if err != nil {
			if err == pgx.ErrNoRows {
				// No existing record, we'll create a new one
				exists = false
				return nil
			}
			return err
		}

		exists = true
		return nil
	})

	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check shop images", nil)
	}

	// Update or create shop images
	err = h.Repository.WithTx(c.Context(), func(q *db.Queries) error {
		if exists {
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
		} else {
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
		}
		return nil
	})

	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update shop images", nil)
	}

	// Return success response
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

	// Auto-publish if publish handler is available
	go func(shopID int64, changeType, entity, description string) {
		// derive worker context from request context
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		ctx, finish := observability.StartSpan(ctx, "autoPublishShopChanges", "store-deployer", "POST", "update-data")
		defer finish(0, nil)
		h.autoPublishShopChangesWithCtx(ctx, shopID, changeType, entity, description)
	}(shopID, "shop_update", fmt.Sprintf("shop:%d", shopID), "Shop images updated")

	return api.SuccessResponse(c, fiber.StatusOK, response, "Shop images updated successfully")
}

// cleanupStoreFiles calls the store-deployer service to cleanup R2 files for a deleted shop
func (h *Handler) cleanupStoreFiles(ctx context.Context, shopID int64, subdomain string) error {
	storeDeployerURL := os.Getenv("STORE_DEPLOYER_URL")
	if storeDeployerURL == "" {
		storeDeployerURL = "http://store-deployer:9003"
	}

	// Prepare the cleanup request
	cleanupReq := map[string]interface{}{
		"shop_id": fmt.Sprintf("%d", shopID),
	}

	reqBody, err := json.Marshal(cleanupReq)
	if err != nil {
		return fmt.Errorf("failed to marshal cleanup request: %v", err)
	}

	// Call the store-deployer cleanup endpoint
	url := fmt.Sprintf("%s/cleanup/%s", storeDeployerURL, subdomain)

	// use a short cancellable context for the cleanup call
	// TODO: callers should pass a context; we enforce a short timeout here as a safeguard.
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create cleanup request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Observability and shared client
	observability.InjectTraceHeaders(ctx, req)
	observability.EnsureRequestID(req)
	ctx, finish := observability.StartSpan(ctx, "cleanupStoreFiles", "store-deployer", "DELETE", url)
	defer finish(0, nil)
	start := time.Now()
	resp, err := ic.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call cleanup endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Read response for logging
	responseBody, _ := io.ReadAll(resp.Body)
	observability.RecordServiceRequest("store-deployer", "DELETE", url, resp.StatusCode, time.Since(start))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cleanup endpoint returned status %d: %s", resp.StatusCode, string(responseBody))
	}

	fmt.Printf("Successfully triggered cleanup for shop %d (subdomain: %s)\n", shopID, subdomain)
	return nil
}

// autoDeployNewShop automatically triggers deployment for newly created shops
func (h *Handler) autoDeployNewShop(shopID int64, subdomain string, templateName string) {
	// delegate to context-aware implementation with a short cancellable worker context
	// TODO: callers should pass a request context; this fallback uses TODO to signal a needed improvement
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	h.autoDeployNewShopWithCtx(ctx, shopID, subdomain, templateName)
}

// autoDeployNewShopWithCtx performs the same work but uses the provided context for cancellation and tracing
func (h *Handler) autoDeployNewShopWithCtx(ctx context.Context, shopID int64, subdomain string, templateName string) {
	select {
	case <-time.After(2 * time.Second):
	case <-ctx.Done():
		return
	}

	startedAt := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	deployment, err := h.Repository.CreateDeployment(ctx, db.CreateDeploymentParams{
		ShopID:          shopID,
		TemplateName:    templateName,
		TemplateVersion: "latest",
		Status:          "deploying",
		DeploymentType:  "full",
		Message:         nil,
		StartedAt:       startedAt,
	})
	if err != nil {
		h.Logger.Error("failed to create deployment record",
			zap.Int64("shop_id", shopID), zap.Error(err))
		return
	}

	deploymentReq := map[string]interface{}{
		"shop_id":       fmt.Sprintf("%d", shopID),
		"subdomain":     subdomain,
		"template_name": templateName,
		"version":       "",
		"data_override": map[string]string{},
	}
	reqBody, _ := json.Marshal(deploymentReq)

	storeDeployerURL := os.Getenv("STORE_DEPLOYER_URL")
	if storeDeployerURL == "" {
		storeDeployerURL = "http://store-deployer:9003"
	}
	url := storeDeployerURL + "/deploy"

	req, err := retryablehttp.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		h.Logger.Error("failed to create retryable request",
			zap.Int64("shop_id", shopID), zap.Error(err))
		return
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	observability.InjectTraceHeaders(ctx, req.Request)
	observability.EnsureRequestID(req.Request)

	start := time.Now()
	resp, err := h.RetryClient.Do(req)
	if err != nil {
		errMsg := err.Error()
		h.Logger.Error("deployment request failed after retries",
			zap.Int64("shop_id", shopID),
			zap.Int64("deployment_id", deployment.DeploymentID),
			zap.Error(err))

		_ = h.Repository.UpdateDeploymentStatus(ctx, db.UpdateDeploymentStatusParams{
			DeploymentID: deployment.DeploymentID,
			Status:       "failed",
			Message:      &errMsg,
		})
		return
	}
	defer resp.Body.Close()

	observability.RecordServiceRequest("store-deployer", "POST", url, resp.StatusCode, time.Since(start))

	h.Logger.Info("auto-deployment triggered",
		zap.Int64("shop_id", shopID),
		zap.String("subdomain", subdomain),
		zap.Int64("deployment_id", deployment.DeploymentID),
		zap.Int("status_code", resp.StatusCode))
}

func (h *Handler) autoPublishShopChangesWithCtx(ctx context.Context, shopID int64, changeType, entity, description string) {
	// Get shop details for subdomain
	shop, err := h.Repository.GetShop(ctx, shopID)
	if err != nil {
		return // Silently fail for auto-publish
	}

	// Call store-deployer to update shop data only
	if err := h.updateStoreDataWithCtx(ctx, shop.Subdomain, shopID, "shop"); err != nil {
		// Log error but don't fail the main operation
		fmt.Printf("Auto-publish failed for shop %d: %v\n", shopID, err)
	}
}

// updateStoreDataWithCtx is a context-aware variant of updateStoreData
func (h *Handler) updateStoreDataWithCtx(ctx context.Context, subdomain string, shopID int64, dataType string) error {
	storeDeployerURL := os.Getenv("STORE_DEPLOYER_URL")
	if storeDeployerURL == "" {
		storeDeployerURL = "http://store-deployer:9003"
	}

	// Prepare the update request
	updateReq := map[string]interface{}{
		"shop_id":   fmt.Sprintf("%d", shopID),
		"data_type": dataType,
	}

	reqBody, err := json.Marshal(updateReq)
	if err != nil {
		return fmt.Errorf("failed to marshal update request: %v", err)
	}

	// Call the store-deployer update-data endpoint
	url := fmt.Sprintf("%s/update-data/%s", storeDeployerURL, subdomain)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create update request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	observability.InjectTraceHeaders(ctx, req)
	observability.EnsureRequestID(req)
	start := time.Now()
	resp, err := ic.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call store-deployer: %v", err)
	}
	defer resp.Body.Close()

	observability.RecordServiceRequest("store-deployer", "POST", url, resp.StatusCode, time.Since(start))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("store-deployer returned status %d", resp.StatusCode)
	}

	return nil
}
