package handlers

import (
	"encoding/json"
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

// CreateProduct Create a new product
// @Summary Create a new product
// @Description Create a new product
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_type_id path string true "Product Type ID"
// @Param product body models.ProductCreateParams true "Product"
// @Success 200 {object} models.SuccessResponse{data=nil} "Product created successfully"
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
	err := h.Repository.WithTx(c.Context(), func(q *db.Queries) error {
		product, err := q.CreateProduct(c.Context(), createProductParams)
		if err != nil {
			return fmt.Errorf("failed to create product: %w", err)
		}

		// Handle attribute values if provided
		if len(productArg.Attributes) > 0 {
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

			// Perform batch upsert for attributes
			batch := q.BatchUpsertProductAttributeValues(c.Context(), attributeValuesParams)
			var batchErr error
			batch.Exec(func(i int, err error) {
				if err != nil {
					batchErr = fmt.Errorf("failed to upsert product attribute values: %w", err)
				}
			})

			if batchErr != nil {
				return batchErr
			}

			if err := batch.Close(); err != nil {
				return fmt.Errorf("failed to close attribute batch: %w", err)
			}
		}

		// Handle variants if provided
		if len(productArg.Variants) > 0 {
			fmt.Println("LEN: ", len(productArg.Variants))
			variantParams := make([]db.UpsertProductVariantsParams, len(productArg.Variants))
			for i, variant := range productArg.Variants {
				// Generate a slug for the variant
				slug := slug.MakeLang(variant.Description, "en")

				variantParams[i] = db.UpsertProductVariantsParams{
					Slug:              slug,
					Description:       variant.Description,
					Price:             variant.Price,
					AvailableQuantity: variant.AvailableQuantity,
					SeoDescription:    variant.SeoDescription,
					SeoKeywords:       variant.SeoKeywords,
					SeoTitle:          variant.SeoTitle,
					ProductID:         product.ProductID,
					ShopID:            shopID,
				}
			}

			// Perform batch upsert for variants
			variantBatch := q.UpsertProductVariants(c.Context(), variantParams)
			var batchErr error
			variantBatch.Exec(func(i int, err error) {
				if err != nil {
					batchErr = fmt.Errorf("failed to upsert product variant %d: %w", i, err)
				}
			})

			if batchErr != nil {
				return batchErr
			}

			if err := variantBatch.Close(); err != nil {
				return fmt.Errorf("failed to upsert product variants: %w", err)
			}
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

	return api.SuccessResponse(c, fiber.StatusCreated, nil, "Product created successfully")
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
	// Parse query params safely
	shopID, err := strconv.ParseInt(c.Params("shop_id", "0"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop_id", nil)
	}

	after, err := strconv.ParseInt(c.Query("after", "0"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid after parameter", nil)
	}

	limit, err := strconv.ParseInt(c.Query("limit", "10"), 10, 32)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid limit parameter", nil)
	}

	param := db.GetProductsParams{
		ShopID: shopID,
		After:  after,
		Limit:  int32(limit),
	}

	objsDB, err := h.Repository.GetProducts(c.Context(), param)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get products", nil)
	}
	products := make([]models.Product, len(objsDB))
	for i, prod := range objsDB {
		var attributes []models.ProductAttribute
		var variants []models.ProductVariant

		if err := json.Unmarshal(prod.Attributes, &attributes); err != nil {
			return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get products attribute", nil)
		}

		if err := json.Unmarshal(prod.Variants, &variants); err != nil {
			return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get products variants", nil)
		}
		products[i] = models.Product{
			ID:          prod.ProductID,
			Title:       prod.Title,
			Description: prod.Description,
			CreatedAt:   prod.CreatedAt,
			UpdatedAt:   prod.UpdatedAt,
			Attributes:  attributes,
			Variants:    variants,
		}
	}

	// No need for extra queries per product. SQL already includes attributes.
	return api.SuccessResponse(c, fiber.StatusOK, products, "Products fetched successfully")
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
	// Parse shop_id
	shopID, err := strconv.ParseInt(c.Params("shop_id", "0"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop_id", nil)
	}

	// Parse product_id
	productID, err := strconv.ParseInt(c.Params("product_id", "0"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid product_id", nil)
	}

	// Fetch product and attributes in a single query
	param := db.GetProductParams{
		ProductID: productID,
		ShopID:    shopID,
	}

	objDB, err := h.Repository.GetProduct(c.Context(), param)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Product not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch product", nil)
	}
	var attributes []models.ProductAttribute
	var variants []models.ProductVariant

	if err := json.Unmarshal(objDB.Attributes, &attributes); err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get products attribute", nil)
	}

	if err := json.Unmarshal(objDB.Variants, &variants); err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get products variants", nil)
	}

	// Return the single query result directly
	return api.SuccessResponse(c, fiber.StatusOK, models.Product{
		ID:          objDB.ProductID,
		Title:       objDB.Title,
		Description: objDB.Description,
		Status:      objDB.Status,
		// CategoryID:  *objDB.CategoryID,
		Attributes: attributes,
		Variants:   variants,
		UpdatedAt:  objDB.UpdatedAt,
		CreatedAt:  objDB.CreatedAt,
	}, "Product fetched successfully")
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
// @Success 200 {object} models.SuccessResponse{data=nil} "Product updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Failure 500 {object} models.ErrorResponse "Failed to update product"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/products/{product_id} [put]
func (h *Handler) UpdateProduct(c *fiber.Ctx) error {
	// Parse path parameters
	shopID, err := strconv.ParseInt(c.Params("shop_id", "0"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop_id", nil)
	}
	productID, err := strconv.ParseInt(c.Params("product_id", "0"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid product_id", nil)
	}

	// Parse request body
	var product models.ProductUpdateParams
	if err := c.BodyParser(&product); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Validate input
	validator := &models.XValidator{}
	if errs := validator.Validate(&product); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}

	// Execute transaction
	err = h.Repository.WithTx(c.Context(), func(q *db.Queries) error {
		// Update base product
		err := q.UpdateProduct(c.Context(), db.UpdateProductParams{
			Title:       product.Title,
			Description: product.Description,
			ProductID:   productID,
			ShopID:      shopID,
		})
		if err != nil {
			if err == pgx.ErrNoRows {
				return fiber.NewError(fiber.StatusNotFound, "Product not found")
			}
			return fmt.Errorf("failed to update product: %w", err)
		}

		// Handle attribute updates if provided
		if len(product.Attributes) > 0 {
			// Get attribute IDs that should be kept
			keepAttributeIDs := make([]int32, len(product.Attributes))
			for i, attr := range product.Attributes {
				keepAttributeIDs[i] = int32(attr.AttributeID)
			}

			// Delete attributes not in the update list
			deleteBatch := q.BatchDeleteProductAttributeValues(c.Context(), []db.BatchDeleteProductAttributeValuesParams{{
				ProductID: productID,
				ShopID:    shopID,
				Column3:   keepAttributeIDs,
			}})
			deleteBatch.Exec(func(_ int, err error) {
				if err != nil {
					return
				}
			})

			// Prepare attribute values for batch upsert
			attributeValuesParams := make([]db.BatchUpsertProductAttributeValuesParams, len(product.Attributes))
			for i, attr := range product.Attributes {
				attributeValuesParams[i] = db.BatchUpsertProductAttributeValuesParams{
					Value:             attr.Value,
					AttributeOptionID: attr.AttributeOptionID,
					ProductID:         productID,
					AttributeID:       attr.AttributeID,
					ShopID:            shopID,
				}
			}

			// Perform batch upsert
			batch := q.BatchUpsertProductAttributeValues(c.Context(), attributeValuesParams)
			batch.Exec(func(i int, err error) {
				if err != nil {
					return
				}
			})

			if err := batch.Close(); err != nil {
				return fmt.Errorf("failed to update attributes: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "Product not found" {
			statusCode = fiber.StatusNotFound
		}
		return api.ErrorResponse(c, statusCode, err.Error(), nil)
	}

	return api.SuccessResponse(c, fiber.StatusOK, nil, "Product updated successfully")
}

// DeleteProduct deletes a product
// @Summary Delete a product
// @Description Delete a product
// @Tags Product
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} models.SuccessResponse{data=nil} "Product deleted successfully"
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

	err := h.Repository.DeleteProduct(c.Context(), param)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Product not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete product", nil)
	}

	return api.SuccessResponse(c, fiber.StatusOK, nil, "Product deleted successfully")
}

// GetProductsByType Get products by product type
// @Summary Get products by product type
// @Description Get products by product type
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_type_id path string true "Product Type ID"
// @Param after query string false "After cursor"
// @Param limit query string false "Limit"
// @Success 200 {object} models.SuccessResponse{data=[]models.Product} "Products fetched successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/product-types/{product_type_id}/products [get]
func (h *Handler) GetProductsByType(c *fiber.Ctx) error {
	// Parse path parameters
	shopID, err := strconv.ParseInt(c.Params("shop_id", "0"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop_id", nil)
	}

	productTypeID, err := strconv.ParseInt(c.Params("product_type_id", "0"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid product_type_id", nil)
	}

	// Parse query parameters
	after, err := strconv.ParseInt(c.Query("after", "0"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid after parameter", nil)
	}

	limit, err := strconv.ParseInt(c.Query("limit", "10"), 10, 32)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid limit parameter", nil)
	}

	param := db.GetProductsByTypeParams{
		ShopID:        shopID,
		ProductTypeID: productTypeID,
		After:         after,
		Limit:         int32(limit),
	}

	products, err := h.Repository.GetProductsByType(c.Context(), param)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch products", nil)
	}

	return api.SuccessResponse(c, fiber.StatusOK, products, "Products fetched successfully")
}
