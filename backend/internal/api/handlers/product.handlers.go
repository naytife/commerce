package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/db/errors"
	"github.com/petrejonn/naytife/internal/observability"
	"go.uber.org/zap"
	// ic and observability not required here; updateStoreDataWithCtx in shop.handlers.go provides tracing
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
// @Success 201 {object} models.SuccessResponse{data=models.Product} "Product created successfully"
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

	if len(productArg.Variants) == 0 {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "At least one product variant is required", nil)
	}

	for i, variant := range productArg.Variants {
		if errs := validator.Validate(&variant); len(errs) > 0 {
			errMsgs := models.FormatValidationErrors(errs)
			return api.ErrorResponse(
				c,
				fiber.StatusBadRequest,
				fmt.Sprintf("Invalid variant at position %d: %s", i+1, errMsgs),
				nil,
			)
		}

		if !variant.Price.Valid || variant.Price.Int == nil || variant.Price.Int.Sign() <= 0 {
			return api.ErrorResponse(
				c,
				fiber.StatusBadRequest,
				fmt.Sprintf("Invalid price for variant %d: price must be greater than 0", i+1),
				nil,
			)
		}
	}
	slug := slug.MakeLang(fmt.Sprint(productArg.Title), "en")
	createProductParams := db.CreateProductParams{
		Title:         productArg.Title,
		Description:   productArg.Description,
		ShopID:        shopID,
		ProductTypeID: productTypeID,
		Status:        db.ProductStatusDRAFT,
		Slug:          slug,
	}

	// Get product type to access the sku_substring
	productTypeParam := db.GetProductTypeParams{
		ProductTypeID: productTypeID,
		ShopID:        shopID,
	}

	productType, err := h.Repository.GetProductType(c.Context(), productTypeParam)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Product type not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch product type", nil)
	}

	// Determine the SKU substring to use
	skuSubstring := "SKU"
	if productType.SkuSubstring != nil {
		skuSubstring = *productType.SkuSubstring
	} else {
		// Use first 3 characters of product type title in uppercase
		if len(productType.Title) >= 3 {
			skuSubstring = strings.ToUpper(productType.Title[:3])
		} else {
			skuSubstring = strings.ToUpper(productType.Title)
		}
	}

	// Execute transaction and capture created product
	var createdProduct db.Product
	err = h.Repository.WithTx(c.Context(), func(q *db.Queries) error {
		product, err := q.CreateProduct(c.Context(), createProductParams)
		if err != nil {
			return err
		}

		// Store the created product for later use
		createdProduct = product

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
					batchErr = err
				}
			})

			if batchErr != nil {
				return batchErr
			}

			if err := batch.Close(); err != nil {
				return err
			}
		}

		// Create variants one by one with proper SKUs
		for i, variant := range productArg.Variants {
			if i > 0 && variant.IsDefault {
				variant.IsDefault = false // Ensure only the first marked default is actually default
			} else if i == 0 && !variant.IsDefault {
				variant.IsDefault = true // Make first variant default if none specified
			}

			// Create the variant first with a temporary placeholder SKU
			// This ensures we get a variation_id
			createVariantParam := db.CreateProductVariationParams{
				Description:       variant.Description,
				Price:             variant.Price,
				AvailableQuantity: variant.AvailableQuantity,
				SeoDescription:    variant.SeoDescription,
				SeoKeywords:       variant.SeoKeywords,
				SeoTitle:          variant.SeoTitle,
				ProductID:         product.ProductID,
				ShopID:            shopID,
				IsDefault:         variant.IsDefault,
				Sku:               fmt.Sprintf("TEMP-%d-%d", product.ProductID, i), // Temporary SKU
			}

			createdVariant, err := q.CreateProductVariation(c.Context(), createVariantParam)
			if err != nil {
				return err
			}

			// Now update with the proper SKU format
			sku := fmt.Sprintf("%s-%d", skuSubstring, createdVariant.ProductVariationID)
			updateParam := db.UpdateProductVariationSkuParams{
				ProductVariationID: createdVariant.ProductVariationID,
				Sku:                sku,
				ShopID:             shopID,
			}

			_, err = q.UpdateProductVariationSku(c.Context(), updateParam)
			if err != nil {
				return err
			}

			// Handle variant attributes if any
			if len(variant.Attributes) > 0 {
				variantAttrParams := make([]db.BatchUpsertProductVariationAttributeValuesParams, len(variant.Attributes))
				for i, attr := range variant.Attributes {
					variantAttrParams[i] = db.BatchUpsertProductVariationAttributeValuesParams{
						Value:              attr.Value,
						AttributeOptionID:  attr.AttributeOptionID,
						ProductVariationID: createdVariant.ProductVariationID,
						AttributeID:        attr.AttributeID,
						ShopID:             shopID,
					}
				}

				batch := q.BatchUpsertProductVariationAttributeValues(c.Context(), variantAttrParams)
				var batchErr error
				batch.Exec(func(i int, err error) {
					if err != nil {
						batchErr = err
					}
				})

				if batchErr != nil {
					return fmt.Errorf("failed to create variant attributes: %w", batchErr)
				}

				if err := batch.Close(); err != nil {
					return fmt.Errorf("failed to commit variant attribute creation: %w", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == errors.UniqueViolation {
				if pgErr.ConstraintName == "products_title_shop_id_key" {
					return api.ErrorResponse(c, fiber.StatusConflict, "A product with this title already exists in your shop", nil)
				}
				return api.ErrorResponse(c, fiber.StatusConflict, fmt.Sprintf("Unique constraint violation: %s", pgErr.ConstraintName), nil)
			}
			if pgErr.Code == errors.ForeignKeyViolation {
				if pgErr.ConstraintName == "fk_attribute" {
					return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid attribute ID provided", nil)
				}
				return api.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Invalid reference: %s", pgErr.ConstraintName), nil)
			}
			return api.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Database error: %s", pgErr.Message), nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create product", nil)
	}

	// Return the created product with its ID
	productResponse := models.Product{
		ID:          createdProduct.ProductID,
		Title:       createdProduct.Title,
		Description: createdProduct.Description,
		Status:      createdProduct.Status,
		CreatedAt:   createdProduct.CreatedAt,
		UpdatedAt:   createdProduct.UpdatedAt,
	}

	shop, err := h.Repository.GetShop(c.Context(), shopID)
	if err != nil {
		h.Logger.Warn("Auto-publish: failed to get shop for auto-publish", zap.Int64("shop_id", shopID), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Product created, but failed to fetch shop for auto-publish", nil)
	}
	// Auto-publish asynchronously via StoreDeployerClient
	go func(shopID int64, subdomain string) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		ctx, finish := observability.StartSpan(ctx, "autoPublishProductChanges", "store-deployer", "POST", "update-data")
		defer finish(0, nil)

		if err := h.StoreDeployerClient.UpdateData(ctx, subdomain, shopID, "products"); err != nil {
			h.Logger.Warn("auto-publish failed",
				zap.Int64("shop_id", shopID),
				zap.Error(err))
		}
	}(shopID, shop.Subdomain)

	return api.SuccessResponse(c, fiber.StatusCreated, productResponse, "Product created successfully")
}

// GetProducts fetches all products
// @Summary Fetch all products
// @Description Fetch all products
// @Tags Product
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
	shopID, err := api.ParseIDParameter(c, "shop_id", "Shop")
	if err != nil {
		return err
	}

	after, err := strconv.ParseInt(c.Query("after", "0"), 10, 64)
	if err != nil {
		return api.BusinessLogicErrorResponse(c, "Invalid after parameter")
	}

	limit, err := strconv.ParseInt(c.Query("limit", "10"), 10, 32)
	if err != nil {
		return api.BusinessLogicErrorResponse(c, "Invalid limit parameter")
	}

	param := db.GetProductsParams{
		ShopID: shopID,
		After:  after,
		Limit:  int32(limit),
	}

	objsDB, err := h.Repository.GetProducts(c.Context(), param)
	if err != nil {
		return api.SystemErrorResponse(c, err, "Failed to get products")
	}
	products := make([]models.Product, len(objsDB))
	for i, prod := range objsDB {
		var attributes []models.ProductAttribute
		var variants []models.ProductVariant
		var images []models.ProductImageResponse

		if err := json.Unmarshal(prod.Attributes, &attributes); err != nil {
			return api.SystemErrorResponse(c, err, "Failed to get products attribute")
		}

		if err := json.Unmarshal(prod.Variants, &variants); err != nil {
			return api.SystemErrorResponse(c, err, "Failed to get products variants")
		}

		if err := json.Unmarshal(prod.Images, &images); err != nil {
			return api.SystemErrorResponse(c, err, "Failed to get product images")
		}

		products[i] = models.Product{
			ID:          prod.ProductID,
			Title:       prod.Title,
			Description: prod.Description,
			CreatedAt:   prod.CreatedAt,
			UpdatedAt:   prod.UpdatedAt,
			Attributes:  attributes,
			Variants:    variants,
			Images:      images,
			Status:      prod.Status,
		}

		// Only set DefaultVariant if the variants array isn't empty
		if len(variants) > 0 {
			products[i].DefaultVariant = variants[0]
			for j := range variants {
				if variants[j].IsDefault {
					products[i].DefaultVariant = variants[j]
					break
				}
			}
		}
	}

	// No need for extra queries per product. SQL already includes attributes.
	return api.SuccessResponse(c, fiber.StatusOK, products, "Products fetched successfully")
}

// GetProduct fetches a single product
// @Summary Fetch a single product
// @Description Fetch a single product
// @Tags Product
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
	shopID, err := api.ParseIDParameter(c, "shop_id", "Shop")
	if err != nil {
		return err
	}

	// Parse product_id
	productID, err := api.ParseIDParameter(c, "product_id", "Product")
	if err != nil {
		return err
	}

	// Fetch product and attributes in a single query
	param := db.GetProductParams{
		ProductID: productID,
		ShopID:    shopID,
	}

	objDB, err := h.Repository.GetProduct(c.Context(), param)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.NotFoundErrorResponse(c, "Product")
		}
		return api.SystemErrorResponse(c, err, "Failed to fetch product")
	}
	var attributes []models.ProductAttribute
	var variants []models.ProductVariant
	var images []models.ProductImageResponse

	if err := json.Unmarshal(objDB.Attributes, &attributes); err != nil {
		return api.SystemErrorResponse(c, err, "Failed to get products attribute")
	}

	if err := json.Unmarshal(objDB.Variants, &variants); err != nil {
		return api.SystemErrorResponse(c, err, "Failed to get products variants")
	}

	if err := json.Unmarshal(objDB.Images, &images); err != nil {
		return api.SystemErrorResponse(c, err, "Failed to get product images")
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
		Images:     images,
		UpdatedAt:  objDB.UpdatedAt,
		CreatedAt:  objDB.CreatedAt,
	}, "Product fetched successfully")
}

// UpdateProduct updates a product
// @Summary Update a product
// @Description Update a product
// @Tags Product
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
	shopID, err := api.ParseIDParameter(c, "shop_id", "Shop")
	if err != nil {
		return err
	}

	productID, err := api.ParseIDParameter(c, "product_id", "Product")
	if err != nil {
		return err
	}

	// Parse request body
	var product models.ProductUpdateParams
	if err := c.BodyParser(&product); err != nil {
		return api.BusinessLogicErrorResponse(c, "Invalid request body")
	}

	// Validate input
	if err := api.ValidateRequest(c, &product); err != nil {
		return err
	}

	// Validate variants if provided
	if len(product.Variants) > 0 {
		for i, variant := range product.Variants {
			validator := &models.XValidator{}
			if errs := validator.Validate(&variant); len(errs) > 0 {
				errMsgs := models.FormatValidationErrors(errs)
				return api.BusinessLogicErrorResponse(c, fmt.Sprintf("Invalid variant at position %d: %s", i+1, errMsgs))
			}

			if !variant.Price.Valid || variant.Price.Int == nil || variant.Price.Int.Sign() <= 0 {
				return api.BusinessLogicErrorResponse(c, fmt.Sprintf("Invalid price for variant %d: price must be greater than 0", i+1))
			}
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

		// Handle variants update if provided
		if len(product.Variants) > 0 {
			// Track variant IDs that should be kept
			variantIDs := make([]int64, 0, len(product.Variants))

			// Update or insert variants
			for i, variant := range product.Variants {
				if variant.ID > 0 {
					// Update existing variant
					variantIDs = append(variantIDs, variant.ID)

					// Create parameters for updating existing variant directly
					updateVariantParam := db.UpdateProductVariationParams{
						ProductVariationID: variant.ID,
						Description:        &variant.Description,
						Price:              variant.Price,
						AvailableQuantity:  &variant.AvailableQuantity,
						SeoDescription:     variant.SeoDescription,
						SeoKeywords:        variant.SeoKeywords,
						SeoTitle:           variant.SeoTitle,
						IsDefault:          &variant.IsDefault,
						ShopID:             shopID,
					}

					_, err := q.UpdateProductVariation(c.Context(), updateVariantParam)
					if err != nil {
						return fmt.Errorf("failed to update variant %d: %w", variant.ID, err)
					}
				} else {
					// It's a new variant, generate a SKU and create
					// Get product type to access the sku_substring
					productTypeParam := db.GetProductTypeByProductParams{
						ProductID: productID,
						ShopID:    shopID,
					}

					productType, err := q.GetProductTypeByProduct(c.Context(), productTypeParam)
					if err != nil {
						return fmt.Errorf("failed to get product type: %w", err)
					}

					// Determine the SKU substring to use
					skuSubstring := "SKU"
					if productType.SkuSubstring != nil {
						skuSubstring = *productType.SkuSubstring
					} else {
						// Use first 3 characters of product type title in uppercase
						if len(productType.Title) >= 3 {
							skuSubstring = strings.ToUpper(productType.Title[:3])
						} else {
							skuSubstring = strings.ToUpper(productType.Title)
						}
					}

					// Create with temporary SKU
					createVariantParam := db.CreateProductVariationParams{
						Description:       variant.Description,
						Price:             variant.Price,
						AvailableQuantity: variant.AvailableQuantity,
						SeoDescription:    variant.SeoDescription,
						SeoKeywords:       variant.SeoKeywords,
						SeoTitle:          variant.SeoTitle,
						ProductID:         productID,
						ShopID:            shopID,
						IsDefault:         variant.IsDefault,
						Sku:               fmt.Sprintf("TEMP-%d-%d", productID, i), // Temporary SKU
					}

					createdVariant, err := q.CreateProductVariation(c.Context(), createVariantParam)
					if err != nil {
						return fmt.Errorf("failed to create variant: %w", err)
					}

					variantIDs = append(variantIDs, createdVariant.ProductVariationID)

					// Update with proper SKU
					sku := fmt.Sprintf("%s-%d", skuSubstring, createdVariant.ProductVariationID)
					updateParam := db.UpdateProductVariationSkuParams{
						ProductVariationID: createdVariant.ProductVariationID,
						Sku:                sku,
						ShopID:             shopID,
					}

					_, err = q.UpdateProductVariationSku(c.Context(), updateParam)
					if err != nil {
						return fmt.Errorf("failed to update variant SKU: %w", err)
					}
				}

				// Handle variant attributes if provided
				if len(variant.Attributes) > 0 {
					variantID := variant.ID
					if variantID == 0 {
						// If it was a newly created variant, get the ID from the last added ID
						variantID = variantIDs[len(variantIDs)-1]
					}

					// Prepare attribute values for batch upsert
					variantAttrParams := make([]db.BatchUpsertProductVariationAttributeValuesParams, len(variant.Attributes))
					for j, attr := range variant.Attributes {
						variantAttrParams[j] = db.BatchUpsertProductVariationAttributeValuesParams{
							Value:              attr.Value,
							AttributeOptionID:  attr.AttributeOptionID,
							ProductVariationID: variantID,
							AttributeID:        attr.AttributeID,
							ShopID:             shopID,
						}
					}

					// Perform batch upsert
					batch := q.BatchUpsertProductVariationAttributeValues(c.Context(), variantAttrParams)
					var batchErr error
					batch.Exec(func(i int, err error) {
						if err != nil {
							batchErr = err
						}
					})

					if batchErr != nil {
						return fmt.Errorf("failed to update variant attributes: %w", batchErr)
					}

					if err := batch.Close(); err != nil {
						return fmt.Errorf("failed to commit variant attribute update: %w", err)
					}
				}
			}

			// Delete any variants not in the update list
			if len(variantIDs) > 0 {
				deleteParams := db.DeleteProductVariantsParams{
					ShopID:              shopID,
					ProductID:           productID,
					ProductVariationIds: variantIDs,
				}
				err := q.DeleteProductVariants(c.Context(), deleteParams)
				if err != nil {
					return fmt.Errorf("failed to delete removed variants: %w", err)
				}
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

	shopIDCopy := shopID
	productIDCopy := productID
	go func() {
		ctx := context.Background()
		h.autoPublishProductChangesWithCtx(ctx, shopIDCopy, "product_update", fmt.Sprintf("product:%d", productIDCopy), "Product updated")
	}()

	return api.SuccessResponse(c, fiber.StatusOK, nil, "Product updated successfully")
}

// DeleteProduct deletes a product
// @Summary Delete a product
// @Description Delete a product
// @Tags Product
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

	shopIDCopy := shopID
	productIDCopy := productID
	go func() {
		ctx := context.Background()
		h.autoPublishProductChangesWithCtx(ctx, shopIDCopy, "product_delete", fmt.Sprintf("product:%d", productIDCopy), "Product deleted")
	}()

	return api.SuccessResponse(c, fiber.StatusOK, nil, "Product deleted successfully")
}

// GetProductsByType Get products by product type
// @Summary Get products by product type
// @Description Get products by product type
// @Tags ProductType
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

	productsDB, err := h.Repository.GetProductsByType(c.Context(), param)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch products", nil)
	}

	products := make([]models.Product, len(productsDB))
	for i, prod := range productsDB {
		var variants []models.ProductVariant
		var images []models.ProductImageResponse

		// Convert attributes to expected type
		var attributesList []models.ProductAttribute
		if attrs, ok := prod.Attributes.([]interface{}); ok {
			for _, attr := range attrs {
				if attrMap, ok := attr.(map[string]interface{}); ok {
					attributesList = append(attributesList, models.ProductAttribute{
						AttributeID:    int64(attrMap["attribute_id"].(float64)),
						AttributeTitle: attrMap["title"].(string),
					})
				}
			}
		}

		if err := json.Unmarshal(prod.Variants, &variants); err != nil {
			return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get products variants", nil)
		}

		if err := json.Unmarshal(prod.Images, &images); err != nil {
			return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get product images", nil)
		}

		products[i] = models.Product{
			ID:          prod.ProductID,
			Title:       prod.Title,
			Description: prod.Description,
			CreatedAt:   prod.CreatedAt,
			UpdatedAt:   prod.UpdatedAt,
			Attributes:  attributesList,
			Variants:    variants,
			Images:      images,
			Status:      prod.Status,
		}

		// Only set DefaultVariant if the variants array isn't empty
		if len(variants) > 0 {
			products[i].DefaultVariant = variants[0]
			for j := range variants {
				if variants[j].IsDefault {
					products[i].DefaultVariant = variants[j]
					break
				}
			}
		}
	}

	return api.SuccessResponse(c, fiber.StatusOK, products, "Products fetched successfully")
}

func (h *Handler) autoPublishProductChangesWithCtx(parentCtx context.Context, shopID int64, changeType, entity, description string) {
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()
	// Get shop details for subdomain
	_, err := h.Repository.GetShop(ctx, shopID)
	if err != nil {
		fmt.Printf("Auto-publish: failed to get shop %d: %v\n", shopID, err)
		return // Silently fail for auto-publish
	}

	// if err := h.updateStoreDataWithCtx(ctx, shop.Subdomain, shopID, "products"); err != nil {
	// 	// Log error but don't block the main operation
	// 	fmt.Printf("Auto-publish failed for shop %d: %v\n", shopID, err)
	// }
}
