package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jinzhu/copier"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/db/errors"
)

// CreateProduct Create a new product
// @Summary Create a new product
// @Description Create a new product
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_type_id path string true "Product Type ID"
// @Param product body models.ProductCreateParams true "Product"
// @Success 200 {object} models.SuccessResponse{data=models.Product} "Product created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/product-types/{product_type_id}/products [post]
func (h *Handler) CreateProduct(c *fiber.Ctx) error {
	// Extract shopID from params
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)

	productTypeIDStr := c.Params("product_type_id", "0")
	productTypeID, _ := strconv.ParseInt(productTypeIDStr, 10, 64)

	// Parse request body
	var productArg models.ProductCreateParams
	if err := c.BodyParser(&productArg); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Validate input
	validator := &models.XValidator{}
	if errs := validator.Validate(&productArg); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}

	// Prepare parameters for product creation
	createProductParams := db.CreateProductParams{
		Title:         productArg.Title,
		Description:   productArg.Description,
		ShopID:        shopID,
		ProductTypeID: productTypeID,
		Status:        db.ProductStatusDRAFT,
	}

	// Execute transaction
	var product db.Product
	var attributeValues []db.ProductAttributeValue
	err := h.Repository.WithTx(c.Context(), func(q *db.Queries) error {
		var err error
		product, err = q.CreateProduct(c.Context(), createProductParams)
		if err != nil {
			return fmt.Errorf("failed to create product: %w", err)
		}

		// Handle attribute values if provided
		if len(productArg.Attributes) == 0 {
			return nil
		}

		// Prepare attribute values for batch insert
		attributeValuesParams := make([]db.BatchUpsertProductAttributeValuesParams, len(productArg.Attributes))
		for i, attr := range productArg.Attributes {
			attributeValuesParams[i] = db.BatchUpsertProductAttributeValuesParams{
				Value:             attr.Value,
				AttributeOptionID: attr.AttributeOptionID,
				ProductID:         product.ProductID,
				AttributeID:       attr.AttributeID,
				ShopID:            shopID,
			}
		}

		// Perform batch upsert
		batch := q.BatchUpsertProductAttributeValues(c.Context(), attributeValuesParams)
		var batchErr error
		batch.Query(func(i int, result []db.ProductAttributeValue, err error) {
			if err != nil {
				batchErr = fmt.Errorf("failed to upsert product attribute values: %w", err)
			} else {
				attributeValues = append(attributeValues, result...)
			}
		})

		if batchErr != nil {
			return batchErr
		}

		if err := batch.Close(); err != nil {
			return fmt.Errorf("batch execution failed: %w", err)
		}

		return nil
	})

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == errors.UniqueViolation {
				return api.ErrorResponse(c, fiber.StatusConflict, "Product already exists", nil)
			}
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create product", nil)
	}

	// Copy result into response struct
	var resp models.Product
	copier.Copy(&resp, &product)
	copier.Copy(&resp.Attributes, &attributeValues)

	return api.SuccessResponse(c, fiber.StatusCreated, resp, "Product created successfully")
}

// GetProducts fetches all products
// @Summary Fetch all products
// @Description Fetch all products
// @Tags Product
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param after query string false "After product ID"
// @Param limit query string false "Limit"
// @Success 200 {object} models.SuccessResponse{data=[]models.Product} "Products fetched successfully"
// @Failure 500 {object} models.ErrorResponse "Failed to fetch products"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/products [get]
func (h *Handler) GetProducts(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	afterStr := c.Query("after", "0")
	after, _ := strconv.ParseInt(afterStr, 10, 64)
	limitStr := c.Query("limit", "10")
	limit, _ := strconv.ParseInt(limitStr, 10, 32)

	param := db.GetProductsParams{
		ShopID: shopID,
		After:  after,
		Limit:  int32(limit),
	}

	objsDB, err := h.Repository.GetProducts(c.Context(), param)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get products", nil)
	}

	var resp []models.Product
	copier.Copy(&resp, &objsDB)
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Products fetched successfully")
}

// GetProduct fetches a single product
// @Summary Fetch a single product
// @Description Fetch a single product
// @Tags Product
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} models.SuccessResponse{data=models.Product} "Product fetched successfully"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Failure 500 {object} models.ErrorResponse "Failed to fetch product"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/products/{product_id} [get]
func (h *Handler) GetProduct(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	productIDStr := c.Params("product_id", "0")
	productID, _ := strconv.ParseInt(productIDStr, 10, 64)

	param := db.GetProductParams{
		ShopID:    shopID,
		ProductID: productID,
	}

	objDB, err := h.Repository.GetProduct(c.Context(), param)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Product not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch product", nil)
	}

	attributeValuesParams := db.GetProductAttributeValuesParams{
		ProductID: productID,
		ShopID:    shopID,
	}

	attributeValues, err := h.Repository.GetProductAttributeValues(c.Context(), attributeValuesParams)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch product attribute values", nil)
	}

	var resp models.Product
	copier.Copy(&resp, &objDB)
	copier.Copy(&resp.Attributes, &attributeValues)

	return api.SuccessResponse(c, fiber.StatusOK, resp, "Product fetched successfully")
}

// UpdateProduct updates a product
// @Summary Update a product
// @Description Update a product
// @Tags Product
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_id path string true "Product ID"
// @Param product body models.ProductUpdateParams true "Product"
// @Success 200 {object} models.SuccessResponse{data=models.Product} "Product updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Failure 500 {object} models.ErrorResponse "Failed to update product"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/products/{product_id} [put]
func (h *Handler) UpdateProduct(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	productIDStr := c.Params("product_id", "0")
	productID, _ := strconv.ParseInt(productIDStr, 10, 64)

	var product models.ProductUpdateParams
	if err := c.BodyParser(&product); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validator := &models.XValidator{}
	if errs := validator.Validate(&product); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}

	param := db.UpdateProductParams{
		Title:       product.Title,
		Description: product.Description,
		ProductID:   productID,
		ShopID:      shopID,
	}

	objDB, err := h.Repository.UpdateProduct(c.Context(), param)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Product not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update product", nil)
	}

	var resp models.Product
	copier.Copy(&resp, &objDB)
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Product updated successfully")
}

// DeleteProduct deletes a product
// @Summary Delete a product
// @Description Delete a product
// @Tags Product
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} models.SuccessResponse{data=models.Product} "Product deleted successfully"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Failure 500 {object} models.ErrorResponse "Failed to delete product"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/products/{product_id} [delete]
func (h *Handler) DeleteProduct(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	productIDStr := c.Params("product_id", "0")
	productID, _ := strconv.ParseInt(productIDStr, 10, 64)

	param := db.DeleteProductParams{
		ProductID: productID,
		ShopID:    shopID,
	}

	objDB, err := h.Repository.DeleteProduct(c.Context(), param)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Product not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete product", nil)
	}

	var resp models.Product
	copier.Copy(&resp, &objDB)
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Product deleted successfully")
}
