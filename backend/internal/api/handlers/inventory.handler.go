package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
)

// GetLowStockVariants fetches product variants with low stock
// @Summary      Get low stock variants
// @Description  Get product variants that are running low on stock
// @Tags         inventory
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        threshold query int false "Stock threshold" default(10)
// @Success      200  {object}   models.SuccessResponse{data=[]models.LowStockVariantResponse} "Low stock variants fetched successfully"
// @Failure      400  {object}   models.ErrorResponse "Bad request"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/inventory/low-stock [get]
func (h *Handler) GetLowStockVariants(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	thresholdStr := c.Query("threshold", "10")
	threshold, _ := strconv.ParseInt(thresholdStr, 10, 32)

	variants, err := h.Repository.GetLowStockVariants(c.Context(), db.GetLowStockVariantsParams{
		ShopID:            shopID,
		AvailableQuantity: threshold,
	})
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch low stock variants", nil)
	}

	// Convert to response format
	variantResponses := make([]models.LowStockVariantResponse, len(variants))
	for i, variant := range variants {
		variantResponses[i] = models.LowStockVariantResponse{
			VariantID:    variant.ProductVariationID,
			ProductID:    variant.ProductID,
			ProductTitle: variant.ProductTitle,
			VariantTitle: &variant.Description, // Using description as variant title
			SKU:          &variant.Sku,
			CurrentStock: int32(variant.AvailableQuantity),
			ReorderLevel: int32(threshold),
			Price:        models.NumericToFloat64(variant.Price),
		}
	}

	return api.SuccessResponse(c, fiber.StatusOK, variantResponses, "Low stock variants fetched successfully")
}

// UpdateVariantStock updates the stock quantity for a product variant
// @Summary      Update variant stock
// @Description  Update the stock quantity for a specific product variant
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        variant_id path string true "Variant ID"
// @Param        stock body models.UpdateStockParams true "Stock update parameters"
// @Success      200  {object}   models.SuccessResponse{data=models.VariantStockResponse} "Stock updated successfully"
// @Failure      400  {object}   models.ErrorResponse "Bad request"
// @Failure      404  {object}   models.ErrorResponse "Variant not found"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/inventory/variants/{variant_id}/stock [put]
func (h *Handler) UpdateVariantStock(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	variantIDStr := c.Params("variant_id")
	variantID, err := strconv.ParseInt(variantIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid variant ID", nil)
	}

	var param models.UpdateStockParams
	err = c.BodyParser(&param)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse request body", nil)
	}

	validator := &models.XValidator{}
	if errs := validator.Validate(&param); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return api.ErrorResponse(c, fiber.StatusBadRequest, errMsgs, nil)
	}

	// Get current stock before update to track the movement
	currentVariant, err := h.Repository.GetProductVariation(c.Context(), db.GetProductVariationParams{
		ProductVariationID: variantID,
		ShopID:             shopID,
	})
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusNotFound, "Product variant not found", nil)
	}

	quantityBefore := int32(currentVariant.AvailableQuantity)
	quantityAfter := param.Quantity
	quantityChange := quantityAfter - quantityBefore

	updatedVariant, err := h.Repository.UpdateVariantStock(c.Context(), db.UpdateVariantStockParams{
		ProductVariationID: variantID,
		ShopID:             shopID,
		AvailableQuantity:  int64(param.Quantity),
	})
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update stock", nil)
	}

	// Create stock movement record
	var notes *string
	if param.Reason != nil {
		notes = param.Reason
	}

	_, err = h.Repository.CreateStockMovement(c.Context(), db.CreateStockMovementParams{
		ProductVariationID: variantID,
		ShopID:             shopID,
		MovementType:       param.MovementType,
		QuantityChange:     quantityChange,
		QuantityBefore:     quantityBefore,
		QuantityAfter:      quantityAfter,
		ReferenceID:        nil, // No reference for manual adjustments
		Notes:              notes,
	})
	if err != nil {
		// Log error but don't fail the request since stock was already updated
		// In production, you might want to use a proper logger
	}

	response := models.VariantStockResponse{
		VariantID: updatedVariant.ProductVariationID,
		Stock:     int32(updatedVariant.AvailableQuantity),
		UpdatedAt: models.TimestamptzToTime(updatedVariant.UpdatedAt),
	}

	return api.SuccessResponse(c, fiber.StatusOK, response, "Stock updated successfully")
}

// AddVariantStock adds stock to a product variant
// @Summary      Add stock to variant
// @Description  Add stock quantity to a specific product variant
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        variant_id path string true "Variant ID"
// @Param        stock body models.AddStockParams true "Stock addition parameters"
// @Success      200  {object}   models.SuccessResponse{data=models.VariantStockResponse} "Stock added successfully"
// @Failure      400  {object}   models.ErrorResponse "Bad request"
// @Failure      404  {object}   models.ErrorResponse "Variant not found"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/inventory/variants/{variant_id}/add-stock [post]
func (h *Handler) AddVariantStock(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	variantIDStr := c.Params("variant_id")
	variantID, err := strconv.ParseInt(variantIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid variant ID", nil)
	}

	var param models.AddStockParams
	err = c.BodyParser(&param)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse request body", nil)
	}

	validator := &models.XValidator{}
	if errs := validator.Validate(&param); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return api.ErrorResponse(c, fiber.StatusBadRequest, errMsgs, nil)
	}

	updatedVariant, err := h.Repository.AddVariantStock(c.Context(), db.AddVariantStockParams{
		ProductVariationID: variantID,
		ShopID:             shopID,
		AvailableQuantity:  int64(param.Quantity),
	})
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to add stock", nil)
	}

	// Create stock movement record
	// TODO: Implement stock movement tracking

	response := models.VariantStockResponse{
		VariantID: updatedVariant.ProductVariationID,
		Stock:     int32(updatedVariant.AvailableQuantity),
		UpdatedAt: models.TimestamptzToTime(updatedVariant.UpdatedAt),
	}

	return api.SuccessResponse(c, fiber.StatusOK, response, "Stock added successfully")
}

// DeductVariantStock deducts stock from a product variant
// @Summary      Deduct stock from variant
// @Description  Deduct stock quantity from a specific product variant
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        variant_id path string true "Variant ID"
// @Param        stock body models.DeductStockParams true "Stock deduction parameters"
// @Success      200  {object}   models.SuccessResponse{data=models.VariantStockResponse} "Stock deducted successfully"
// @Failure      400  {object}   models.ErrorResponse "Bad request"
// @Failure      404  {object}   models.ErrorResponse "Variant not found"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/inventory/variants/{variant_id}/deduct-stock [post]
func (h *Handler) DeductVariantStock(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	variantIDStr := c.Params("variant_id")
	variantID, err := strconv.ParseInt(variantIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid variant ID", nil)
	}

	var param models.DeductStockParams
	err = c.BodyParser(&param)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse request body", nil)
	}

	validator := &models.XValidator{}
	if errs := validator.Validate(&param); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return api.ErrorResponse(c, fiber.StatusBadRequest, errMsgs, nil)
	}

	updatedVariant, err := h.Repository.DeductVariantStock(c.Context(), db.DeductVariantStockParams{
		ProductVariationID: variantID,
		ShopID:             shopID,
		AvailableQuantity:  int64(param.Quantity),
	})
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to deduct stock", nil)
	}

	// Create stock movement record
	var referenceID *int64
	if param.ReferenceID != nil {
		id, err := strconv.ParseInt(*param.ReferenceID, 10, 64)
		if err == nil {
			referenceID = &id
		}
	}

	_, err = h.Repository.CreateStockMovement(c.Context(), db.CreateStockMovementParams{
		ProductVariationID: variantID,
		ShopID:             shopID,
		MovementType:       "sale",
		QuantityChange:     -param.Quantity,                                                 // Negative for deduction
		QuantityBefore:     int32(updatedVariant.AvailableQuantity + int64(param.Quantity)), // Previous stock
		QuantityAfter:      int32(updatedVariant.AvailableQuantity),                         // Current stock
		ReferenceID:        referenceID,
		Notes:              &param.Reason,
	})
	if err != nil {
		// Log error but don't fail the request
		// In production, you might want to use a proper logger
	}

	response := models.VariantStockResponse{
		VariantID: updatedVariant.ProductVariationID,
		Stock:     int32(updatedVariant.AvailableQuantity),
		UpdatedAt: models.TimestamptzToTime(updatedVariant.UpdatedAt),
	}

	return api.SuccessResponse(c, fiber.StatusOK, response, "Stock deducted successfully")
}

// GetInventoryReport generates an inventory report
// @Summary      Get inventory report
// @Description  Generate comprehensive inventory report for a shop
// @Tags         inventory
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Success      200  {object}   models.SuccessResponse{data=models.InventoryReportResponse} "Inventory report generated successfully"
// @Failure      400  {object}   models.ErrorResponse "Bad request"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/inventory/report [get]
func (h *Handler) GetInventoryReport(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	report, err := h.Repository.GetInventoryReport(c.Context(), db.GetInventoryReportParams{
		ShopID:            shopID,
		AvailableQuantity: 10, // Default threshold for low stock
	})
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate inventory report", nil)
	}

	// Calculate totals from the report data
	totalProducts := len(report)
	totalVariants := len(report)
	var totalStockValue float64
	lowStockCount := 0
	outOfStockCount := 0

	// Create inventory items array
	items := make([]models.InventoryItemResponse, len(report))

	for i, item := range report {
		totalStockValue += models.NumericToFloat64(item.Price) * float64(item.AvailableQuantity)
		if item.StockStatus == "LOW_STOCK" {
			lowStockCount++
		}
		if item.StockStatus == "OUT_OF_STOCK" {
			outOfStockCount++
		}

		// Create inventory item response
		items[i] = models.InventoryItemResponse{
			VariantID:         item.ProductVariationID,
			ProductID:         item.ProductID,
			ProductTitle:      item.ProductTitle,
			VariantTitle:      item.VariantDescription,
			SKU:               item.Sku,
			CurrentStock:      item.AvailableQuantity,
			ReservedStock:     int64(item.ReservedQuantity),
			AvailableStock:    item.AvailableStock,
			LowStockThreshold: 10,  // TODO: Get from product variation settings
			Location:          nil, // TODO: Add location field
			LastUpdated:       models.TimestamptzToTime(item.UpdatedAt).Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response := models.InventoryReportResponse{
		TotalProducts:   totalProducts,
		TotalVariants:   totalVariants,
		TotalStockValue: totalStockValue,
		LowStockCount:   lowStockCount,
		OutOfStockCount: outOfStockCount,
		GeneratedAt:     time.Now(),
		Items:           items,
	}

	return api.SuccessResponse(c, fiber.StatusOK, response, "Inventory report generated successfully")
}

// GetStockMovements fetches stock movement history
// @Summary      Get stock movement history
// @Description  Get paginated stock movement history for a shop or specific variant
// @Tags         inventory
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        variant_id query string false "Variant ID to filter by"
// @Param        limit query int false "Limit" default(50)
// @Param        offset query int false "Offset" default(0)
// @Success      200  {object}   models.SuccessResponse{data=models.StockMovementsResponse} "Stock movements fetched successfully"
// @Failure      400  {object}   models.ErrorResponse "Bad request"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/inventory/movements [get]
func (h *Handler) GetStockMovements(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	variantIDStr := c.Query("variant_id")
	// Note: Currently the query doesn't support filtering by variant_id
	// This could be implemented by modifying the SQL query
	if variantIDStr != "" {
		_, err := strconv.ParseInt(variantIDStr, 10, 64)
		if err != nil {
			return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid variant ID", nil)
		}
		// TODO: Add variant filtering to the query
	}

	limitStr := c.Query("limit", "50")
	offsetStr := c.Query("offset", "0")

	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	offset, _ := strconv.ParseInt(offsetStr, 10, 32)

	movements, err := h.Repository.GetStockMovements(c.Context(), db.GetStockMovementsParams{
		ShopID: shopID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch stock movements", nil)
	}

	// Convert to response format
	movementResponses := make([]models.StockMovementResponse, len(movements))
	for i, movement := range movements {
		var referenceIDStr *string
		if movement.ReferenceID != nil {
			idStr := strconv.FormatInt(*movement.ReferenceID, 10)
			referenceIDStr = &idStr
		}

		// Convert int32 to *int32 for optional fields
		previousStock := &movement.QuantityBefore
		newStock := &movement.QuantityAfter

		movementResponses[i] = models.StockMovementResponse{
			MovementID:    movement.MovementID,
			VariantID:     movement.ProductVariationID,
			MovementType:  movement.MovementType,
			Quantity:      movement.QuantityChange,
			PreviousStock: previousStock,
			NewStock:      newStock,
			Reason:        movement.Notes,
			ReferenceType: nil, // Not available in current schema
			ReferenceID:   referenceIDStr,
			CreatedAt:     models.TimestamptzToTime(movement.CreatedAt),
			ProductTitle:  movement.ProductTitle,
			VariantTitle:  &movement.VariantTitle,
			SKU:           movement.Sku,
		}
	}

	response := models.StockMovementsResponse{
		Movements: movementResponses,
		Total:     len(movementResponses),
		Page:      int(offset/limit) + 1,
		Limit:     int(limit),
	}

	return api.SuccessResponse(c, fiber.StatusOK, response, "Stock movements fetched successfully")
}
