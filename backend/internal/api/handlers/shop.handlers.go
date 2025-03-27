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
// @Accept       json
// @Produce      json
// @Success      200  {object}   models.SuccessResponse{data=[]models.Shop} "Shops fetched successfully"
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
	var resp []models.Shop
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
		resp = append(resp, shop)
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Shops retrieved")
}

// DeleteShop deletes a shop
// @Summary      Delete a shop
// @Description
// @Tags         shops
// @Accept       json
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
// @Accept       json
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
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Shop retrieved")
}

// GetShopBySubDomain fetches a shop by subdomain
// @Summary      Fetch a shop by subdomain
// @Description
// @Tags         shops
// @Accept       json
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
		CustomDomain:        objDB.Domain,
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
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Shop retrieved")
}

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
