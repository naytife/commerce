package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/observability"
	"go.uber.org/zap"
)

// AddProductImage adds a new image to a product
// @Summary Add a new image to a product
// @Description Add a new image to a product by providing the URL and alt text
// @Tags Product Images
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_id path string true "Product ID"
// @Param image body models.ProductImageCreateParams true "Image data"
// @Success 200 {object} models.SuccessResponse{data=models.ProductImageResponse} "Image added successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/products/{product_id}/images [post]
func (h *Handler) AddProductImage(c *fiber.Ctx) error {
	// Parse shop ID and product ID from params
	shopIDStr := c.Params("shop_id", "0")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	productIDStr := c.Params("product_id", "0")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid product ID", nil)
	}

	// Parse request body
	var imageParams models.ProductImageCreateParams
	if err := c.BodyParser(&imageParams); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Validate input
	validator := &models.XValidator{}
	if errs := validator.Validate(&imageParams); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}

	// Verify product exists
	productExists := false
	err = h.Repository.WithTx(c.Context(), func(q *db.Queries) error {
		_, err := q.GetProductById(c.Context(), db.GetProductByIdParams{
			ProductID: productID,
			ShopID:    shopID,
		})
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil // Product not found, we'll handle this after the transaction
			}
			return err
		}

		productExists = true
		return nil
	})

	if err != nil {
		zap.L().Error("AddProductImage: failed to verify product", zap.Error(err), zap.Int64("shop_id", shopID), zap.Int64("product_id", productID))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to verify product", nil)
	}

	if !productExists {
		zap.L().Warn("AddProductImage: product not found", zap.Int64("shop_id", shopID), zap.Int64("product_id", productID))
		return api.ErrorResponse(c, fiber.StatusNotFound, "Product not found", nil)
	}

	// Create product image
	var createdImage db.ProductImage
	err = h.Repository.WithTx(c.Context(), func(q *db.Queries) error {
		image, err := q.CreateProductImage(c.Context(), db.CreateProductImageParams{
			Url:       imageParams.URL,
			Alt:       imageParams.Alt,
			ProductID: productID,
			ShopID:    shopID,
		})
		if err != nil {
			return err
		}

		createdImage = image
		return nil
	})

	if err != nil {
		zap.L().Error("AddProductImage: failed to add product image", zap.Error(err), zap.Int64("shop_id", shopID), zap.Int64("product_id", productID))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to add product image", nil)
	}

	// Fetch shop for subdomain
	shop, err := h.Repository.GetShop(c.Context(), shopID)
	if err != nil {
		zap.L().Warn("Auto-publish: failed to get shop for auto-publish", zap.Int64("shop_id", shopID), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Image added, but failed to fetch shop for auto-publish", nil)
	}

	// Auto-publish asynchronously via StoreDeployerClient
	go func(shopID int64, subdomain string) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		ctx, finish := observability.StartSpan(ctx, "autoPublishProductImages", "store-deployer", "POST", "update-data")
		defer finish(0, nil)

		if err := h.StoreDeployerClient.UpdateData(ctx, subdomain, shopID, "products"); err != nil {
			zap.L().Warn("auto-publish shop data failed",
				zap.Int64("shop_id", shopID),
				zap.Error(err))
		}
	}(shopID, shop.Subdomain)

	// Return success response
	response := models.ProductImageResponse{
		ID:  createdImage.ProductImageID,
		URL: createdImage.Url,
		Alt: createdImage.Alt,
	}
	return api.SuccessResponse(c, fiber.StatusOK, response, "Image added successfully")
}

// GetProductImages gets all images for a product
// @Summary Get all images for a product
// @Description Get all images for a product
// @Tags Product Images
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} models.SuccessResponse{data=models.ProductImagesResponse} "Images retrieved successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/products/{product_id}/images [get]
func (h *Handler) GetProductImages(c *fiber.Ctx) error {
	// Parse shop ID and product ID from params
	shopIDStr := c.Params("shop_id", "0")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	productIDStr := c.Params("product_id", "0")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid product ID", nil)
	}

	// Verify product exists
	productExists := false
	err = h.Repository.WithTx(c.Context(), func(q *db.Queries) error {
		_, err := q.GetProductById(c.Context(), db.GetProductByIdParams{
			ProductID: productID,
			ShopID:    shopID,
		})
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil // Product not found, we'll handle this after the transaction
			}
			return err
		}

		productExists = true
		return nil
	})

	if err != nil {
		zap.L().Error("GetProductImages: failed to verify product", zap.Error(err), zap.Int64("shop_id", shopID), zap.Int64("product_id", productID))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to verify product", nil)
	}

	if !productExists {
		zap.L().Warn("GetProductImages: product not found", zap.Int64("shop_id", shopID), zap.Int64("product_id", productID))
		return api.ErrorResponse(c, fiber.StatusNotFound, "Product not found", nil)
	}

	// Get product images
	var productImages []db.ProductImage
	err = h.Repository.WithTx(c.Context(), func(q *db.Queries) error {
		images, err := q.GetProductImages(c.Context(), db.GetProductImagesParams{
			ProductID: productID,
			ShopID:    shopID,
		})
		if err != nil {
			return err
		}

		productImages = images
		return nil
	})

	if err != nil {
		zap.L().Error("GetProductImages: failed to get product images", zap.Error(err), zap.Int64("shop_id", shopID), zap.Int64("product_id", productID))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get product images", nil)
	}

	// Map database images to response images
	responseImages := make([]models.ProductImageResponse, len(productImages))
	for i, image := range productImages {
		responseImages[i] = models.ProductImageResponse{
			ID:  image.ProductImageID,
			URL: image.Url,
			Alt: image.Alt,
		}
	}

	// Return success response
	response := models.ProductImagesResponse{
		Images: responseImages,
	}
	return api.SuccessResponse(c, fiber.StatusOK, response, "Images retrieved successfully")
}

// DeleteProductImage deletes a product image
// @Summary Delete a product image
// @Description Delete a product image by ID
// @Tags Product Images
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_id path string true "Product ID"
// @Param image_id path string true "Image ID"
// @Success 200 {object} models.SuccessResponse{data=nil} "Image deleted successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/products/{product_id}/images/{image_id} [delete]
func (h *Handler) DeleteProductImage(c *fiber.Ctx) error {
	// Parse shop ID from params
	shopIDStr := c.Params("shop_id", "0")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	// We don't need to verify product ID here since we're deleting by image ID
	// and the database constraint ensures image belongs to a valid product

	// Parse image ID from params
	imageIDStr := c.Params("image_id", "0")
	imageID, err := strconv.ParseInt(imageIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid image ID", nil)
	}

	// Delete product image
	err = h.Repository.WithTx(c.Context(), func(q *db.Queries) error {
		return q.DeleteProductImage(c.Context(), db.DeleteProductImageParams{
			ProductImageID: imageID,
			ShopID:         shopID,
		})
	})

	if err != nil {
		zap.L().Error("DeleteProductImage: failed to delete product image", zap.Error(err), zap.Int64("shop_id", shopID), zap.Int64("image_id", imageID))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete product image", nil)
	}

	// Fetch shop for subdomain
	shop, err := h.Repository.GetShop(c.Context(), shopID)
	if err != nil {
		zap.L().Warn("Auto-publish: failed to get shop for auto-publish", zap.Int64("shop_id", shopID), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Image deleted, but failed to fetch shop for auto-publish", nil)
	}

	// Auto-publish asynchronously via StoreDeployerClient
	go func(shopID int64, subdomain string) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		ctx, finish := observability.StartSpan(ctx, "autoPublishProductImages", "store-deployer", "POST", "update-data")
		defer finish(0, nil)

		if err := h.StoreDeployerClient.UpdateData(ctx, subdomain, shopID, "products"); err != nil {
			zap.L().Warn("auto-publish shop data failed",
				zap.Int64("shop_id", shopID),
				zap.Error(err))
		}
	}(shopID, shop.Subdomain)

	// Return success response
	return api.SuccessResponse(c, fiber.StatusOK, nil, "Image deleted successfully")
}
