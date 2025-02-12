package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jinzhu/copier"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
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
// @Security     OAuth2AccessCode
// @Router       /shops [post]
func (h *Handler) CreateShop(c *fiber.Ctx) error {
	// TODO: verify user exist
	userSub, _ := c.Locals("user_id").(string)
	param := db.CreateShopParams{}
	var shop models.ShopCreateParams
	c.BodyParser(&shop)
	user, err := h.Repository.GetUserBySub(c.Context(), &userSub)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusUnauthorized, "Failed to get profile", nil)
	}
	shop.Status = string(db.ProductStatusDRAFT)
	shop.CurrencyCode = "NGN"
	if !slug.IsSlug(shop.Domain) {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid domain format", nil)
	}

	validator := &models.XValidator{}
	if errs := validator.Validate(&shop); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)

		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}
	copier.Copy(&param, &shop)
	param.OwnerID = user.UserID
	objDB, err := h.Repository.CreateShop(c.Context(), param)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {

			if pgErr.Code == "23505" {
				return api.ErrorResponse(c, fiber.StatusConflict, "Shop already exist", nil)
			}
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create shop", nil)
	}
	var resp models.Shop
	copier.Copy(&resp, &objDB)
	return api.SuccessResponse(c, fiber.StatusCreated, resp, "Shop created")
}

// GetShops fetches auth user shop
// @Summary      Fetch all shops
// @Description
// @Tags         shops
// @Accept       json
// @Produce      json
// @Success      200  {object}   models.SuccessResponse{data=[]models.Shop} "Shops fetched successfully"
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
	copier.Copy(&resp, &objsDB)
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Shops retrieved")
}

// DeleteShop deletes a shop
// @Summary      Delete a shop
// @Description
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id} [delete]
func (h *Handler) DeleteShop(c *fiber.Ctx) error {
	shopID := c.Params("shop_id", "0")
	shopIDInt, _ := strconv.ParseInt(shopID, 10, 64)
	dberr := h.Repository.DeleteShop(c.Context(), shopIDInt)
	if dberr != nil {
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
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id} [get]
func (h *Handler) GetShop(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}

// UpdateShop updates a shop
// @Summary      Update a shop
// @Description
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        shop body models.ShopUpdateParams true "Shop object that needs to be updated"
// @Success      200  {object}   models.SuccessResponse{data=models.Shop} "Shop updated successfully"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id} [put]
func (h *Handler) UpdateShop(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}
