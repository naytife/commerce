package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jinzhu/copier"
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

func (h *Handler) CreateShop(c *fiber.Ctx) error {
	// TODO: verify user exist
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve user ID")
	}

	param := db.CreateShopParams{}
	var shop models.Shop
	c.BodyParser(&shop)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: "Invalid user ID format",
		}
	}
	shop.OwnerID = userID
	shop.Status = string(DRAFT)
	shop.CurrencyCode = "NGN"
	if !slug.IsSlug(shop.Domain) {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: "Invalid domain format",
		}
	}

	validator := &models.XValidator{}
	if errs := validator.Validate(&shop); len(errs) > 0 && errs[0].Error {
		errMsgs := models.FormatValidationErrors(errs)

		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}
	copier.Copy(&param, &shop)
	objDB, err := h.Repository.CreateShop(c.Context(), param)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {

			if pgErr.Code == "23505" {
				return &fiber.Error{
					Code:    fiber.ErrConflict.Code,
					Message: "Shop already exists",
				}
			}
		}
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: "Failed to create shop",
		}
	}
	return c.JSON(models.ResponseHTTP{
		Success: true,
		Data:    objDB,
		Message: "Shop created",
	})
}

func (h *Handler) GetShops(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve user ID")
	}
	objsDB, err := h.Repository.GetShopsByOwner(c.Context(), uuid.MustParse(userIDStr))
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: "Failed to get shops",
		}
	}
	return c.JSON(models.ResponseHTTP{
		Success: true,
		Data:    objsDB,
		Message: "Shops retrieved",
	})
}

func (h *Handler) DeleteShop(c *fiber.Ctx) error {
	shopID := c.Params("id")
	shopIDInt, err := strconv.ParseInt(shopID, 10, 64)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: "Invalid shop ID format",
		}
	}
	dberr := h.Repository.DeleteShop(c.Context(), shopIDInt)
	if dberr != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: "Failed to delete shop",
		}
	}
	return c.JSON(models.ResponseHTTP{
		Success: true,
		Data:    nil,
		Message: "Shop deleted",
	})
}
