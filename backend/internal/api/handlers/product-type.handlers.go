package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/db/errors"
)

// CreateProductType Creeate a new product type
// @Summary Create a new product type
// @Description Create a new product type
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param productType body models.ProductTypeCreateParams true "Product type object that needs to be created"
// @Success 200 {object} models.SuccessResponse{data=models.ProductType} "Product type created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 409 {object} models.ErrorResponse "Product type already exists"
// @Failure 500 {object} models.ErrorResponse "Failed to create product type"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/product-types [post]
func (h *Handler) CreateProductType(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)

	var productType models.ProductTypeCreateParams
	if err := c.BodyParser(&productType); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validator := &models.XValidator{}
	if errs := validator.Validate(&productType); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}

	param := db.CreateProductTypeParams{
		Title:     productType.Title,
		Shippable: productType.Shippable,
		Digital:   productType.Digital,
		ShopID:    shopID,
	}

	objDB, err := h.Repository.CreateProductType(c.Context(), param)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == errors.UniqueViolation {
				return api.ErrorResponse(c, fiber.StatusConflict, "Product type already exists", nil)
			}
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create product type", nil)
	}

	resp := models.ProductType{
		ID:        objDB.ProductTypeID,
		Title:     objDB.Title,
		Shippable: objDB.Shippable,
		Digital:   objDB.Digital,
	}
	return api.SuccessResponse(c, fiber.StatusCreated, resp, "Product type created")
}

// GetProductTypes fetches all product types
// @Summary Fetch all product types
// @Description Fetch all product types
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Success 200 {object} models.SuccessResponse{data=models.ProductType} "Product types fetched successfully"
// @Failure 500 {object} models.ErrorResponse "Failed to fetch product types"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/product-types [get]
func (h *Handler) GetProductTypes(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)

	objsDB, err := h.Repository.GetProductTypes(c.Context(), shopID)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get product types", nil)
	}

	resp := make([]models.ProductType, len(objsDB))
	for i, objDB := range objsDB {
		resp[i] = models.ProductType{
			ID:        objDB.ProductTypeID,
			Title:     objDB.Title,
			Shippable: objDB.Shippable,
			Digital:   objDB.Digital,
		}
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Product types fetched successfully")
}

// GetProductType fetches a single product type
// @Summary Fetch a single product type
// @Description Fetch a single product type
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_type_id path string true "Product type ID"
// @Success 200 {object} models.SuccessResponse{data=models.ProductType} "Product type fetched successfully"
// @Failure 404 {object} models.ErrorResponse "Product type not found"
// @Failure 500 {object} models.ErrorResponse "Failed to fetch product type"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/product-types/{product_type_id} [get]
func (h *Handler) GetProductType(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	productTypeIDStr := c.Params("product_type_id", "0")
	productTypeID, _ := strconv.ParseInt(productTypeIDStr, 10, 64)

	param := db.GetProductTypeParams{
		ProductTypeID: productTypeID,
		ShopID:        shopID,
	}

	objDB, err := h.Repository.GetProductType(c.Context(), param)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Product type not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch product type", nil)
	}

	resp := models.ProductType{
		ID:        objDB.ProductTypeID,
		Title:     objDB.Title,
		Shippable: objDB.Shippable,
		Digital:   objDB.Digital,
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Product type fetched successfully")
}

// UpdateProductType updates a product type
// @Summary Update a product type
// @Description Update a product type
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_type_id path string true "Product type ID"
// @Param productType body models.ProductTypeUpdateParams true "Product type object that needs to be updated"
// @Success 200 {object} models.SuccessResponse{data=models.ProductType} "Product type updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 404 {object} models.ErrorResponse "Product type not found"
// @Failure 409 {object} models.ErrorResponse "Product type already exists"
// @Failure 500 {object} models.ErrorResponse "Failed to update product type"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/product-types/{product_type_id} [put]
func (h *Handler) UpdateProductType(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	productTypeIDStr := c.Params("product_type_id", "0")
	productTypeID, _ := strconv.ParseInt(productTypeIDStr, 10, 64)

	var productType models.ProductTypeUpdateParams
	if err := c.BodyParser(&productType); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validator := &models.XValidator{}
	if errs := validator.Validate(&productType); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}

	param := db.UpdateProductTypeParams{
		Title:         productType.Title,
		Shippable:     productType.Shippable,
		Digital:       productType.Digital,
		ProductTypeID: productTypeID,
		ShopID:        shopID,
	}

	objDB, err := h.Repository.UpdateProductType(c.Context(), param)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Product type not found", nil)
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == errors.UniqueViolation {
				return api.ErrorResponse(c, fiber.StatusConflict, "Product type already exists", nil)
			}
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update product type", nil)
	}

	resp := models.ProductType{
		ID:        objDB.ProductTypeID,
		Title:     objDB.Title,
		Shippable: objDB.Shippable,
		Digital:   objDB.Digital,
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Product type updated successfully")
}

// DeleteProductType deletes a product type
// @Summary Delete a product type
// @Description Delete a product type
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_type_id path string true "Product type ID"
// @Success 200 {object} models.SuccessResponse{data=models.ProductType} "Product type deleted successfully"
// @Failure 404 {object} models.ErrorResponse "Product type not found"
// @Failure 500 {object} models.ErrorResponse "Failed to delete product type"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/product-types/{product_type_id} [delete]
func (h *Handler) DeleteProductType(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	productTypeIDStr := c.Params("product_type_id", "0")
	productTypeID, _ := strconv.ParseInt(productTypeIDStr, 10, 64)

	param := db.DeleteProductTypeParams{
		ProductTypeID: productTypeID,
		ShopID:        shopID,
	}

	objDB, err := h.Repository.DeleteProductType(c.Context(), param)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Product type not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete product type", nil)
	}

	resp := models.ProductType{
		ID:        objDB.ProductTypeID,
		Title:     objDB.Title,
		Shippable: objDB.Shippable,
		Digital:   objDB.Digital,
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Product type deleted successfully")
}
