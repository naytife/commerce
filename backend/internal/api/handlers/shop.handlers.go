package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/db/errors"
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
	// TODO: verify user exist
	userSub, _ := c.Locals("user_id").(string)
	var shop models.ShopCreateParams
	c.BodyParser(&shop)
	user, err := h.Repository.GetUserBySub(c.Context(), &userSub)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusUnauthorized, "Failed to get profile", nil)
	}
	shop.Status = string(db.ProductStatusDRAFT)
	shop.CurrencyCode = "NGN"
	if !slug.IsSlug(shop.Subdomain) {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid subdomain format", nil)
	}

	validator := &models.XValidator{}
	if errs := validator.Validate(&shop); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)

		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}
	param := db.CreateShopParams{
		OwnerID:      user.UserID,
		Title:        shop.Title,
		Subdomain:    shop.Subdomain,
		CurrencyCode: shop.CurrencyCode,
		Status:       shop.Status,
	}
	objDB, err := h.Repository.CreateShop(c.Context(), param)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {

			if pgErr.Code == errors.UniqueViolation {
				return api.ErrorResponse(c, fiber.StatusConflict, "Shop already exist", nil)
			}
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create shop", nil)
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
	}
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
	user, err := h.Repository.GetUserBySub(c.Context(), &userSub)
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
		}

		// Fetch shop images
		shopImages, imgErr := h.Repository.GetShopImages(c.Context(), obj.ShopID)
		if imgErr == nil {
			// Only add images if successfully retrieved
			shop.Images = &models.ShopImagesResponse{
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
	err := h.Repository.DeleteShop(c.Context(), shopIDInt)
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
	}

	// Fetch shop images
	shopImages, err := h.Repository.GetShopImages(c.Context(), shopID)
	if err == nil {
		// Only add images if successfully retrieved
		resp.Images = &models.ShopImagesResponse{
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
// @Router       /shops/subdomain/{subdomain} [get]
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
	}

	// Fetch shop images
	shopImages, err := h.Repository.GetShopImages(c.Context(), objDB.ShopID)
	if err == nil {
		// Only add images if successfully retrieved
		resp.Images = &models.ShopImagesResponse{
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
	}
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
// @Success      200  {object}   models.SuccessResponse{data=models.ShopImagesResponse} "Shop images updated successfully"
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
	response := models.ShopImagesResponse{
		ID:                updatedImages.ShopImageID,
		FaviconUrl:        updatedImages.FaviconUrl,
		LogoUrl:           updatedImages.LogoUrl,
		LogoUrlDark:       updatedImages.LogoUrlDark,
		BannerUrl:         updatedImages.BannerUrl,
		BannerUrlDark:     updatedImages.BannerUrlDark,
		CoverImageUrl:     updatedImages.CoverImageUrl,
		CoverImageUrlDark: updatedImages.CoverImageUrlDark,
	}

	return api.SuccessResponse(c, fiber.StatusOK, response, "Shop images updated successfully")
}
