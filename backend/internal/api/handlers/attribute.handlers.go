package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/db/errors"
)

// CreateAttribute Create a new attribute
// @Summary Create a new attribute
// @Description Create a new attribute
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_type_id path string true "Product Type ID"
// @Param attribute body models.AttributeCreateParams true "Attribute"
// @Success 200 {object} models.SuccessResponse{data=models.Attribute} "Attribute created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/product-types/{product_type_id}/attributes [post]
func (h *Handler) CreateAttribute(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	productTypeIDStr := c.Params("product_type_id", "0")
	productTypeID, _ := strconv.ParseInt(productTypeIDStr, 10, 64)

	var attribute models.AttributeCreateParams
	if err := c.BodyParser(&attribute); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validator := &models.XValidator{}
	if errs := validator.Validate(&attribute); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}

	param := db.CreateAttributeParams{
		Title:         attribute.Title,
		DataType:      attribute.DataType,
		Unit:          db.NullAttributeUnit{AttributeUnit: attribute.Unit, Valid: attribute.Unit != ""},
		Required:      attribute.Required,
		AppliesTo:     attribute.AppliesTo,
		ProductTypeID: productTypeID,
		ShopID:        shopID,
	}

	objDB, err := h.Repository.CreateAttribute(c.Context(), param)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == errors.UniqueViolation {
				return api.ErrorResponse(c, fiber.StatusConflict, "Attribute already exists", nil)
			}
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create attribute", nil)
	}

	resp := models.Attribute{
		ID:            objDB.AttributeID,
		Title:         objDB.Title,
		DataType:      objDB.DataType,
		Unit:          objDB.Unit.AttributeUnit,
		Required:      objDB.Required,
		AppliesTo:     objDB.AppliesTo,
		ProductTypeID: objDB.ProductTypeID,
	}
	return api.SuccessResponse(c, fiber.StatusCreated, resp, "Attribute created successfully")
}

// GetAttributes fetches all attributes
// @Summary Fetch all attributes
// @Description Fetch all attributes
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_type_id path string true "Product Type ID"
// @Param for query string false "For (Product or ProductVariation)"
// @Success 200 {object} models.SuccessResponse{data=[]models.Attribute} "Attributes fetched successfully"
// @Failure 500 {object} models.ErrorResponse "Failed to fetch attributes"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/product-types/{product_type_id}/attributes [get]
func (h *Handler) GetAttributes(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	productTypeIDStr := c.Params("product_type_id", "0")
	productTypeID, _ := strconv.ParseInt(productTypeIDStr, 10, 64)
	// Validate that 'for' parameter is either Product or ProductVariation
	if appliesTo := c.Query("for", ""); appliesTo != "" && appliesTo != "Product" && appliesTo != "ProductVariation" {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Parameter 'for' must be either 'Product' or 'ProductVariation'", nil)
	}
	appliesTo := c.Query("for", "")

	params := db.GetAttributesParams{
		ProductTypeID: productTypeID,
		ShopID:        shopID,
		AppliesTo:     appliesTo, // Add this to filter by type
	}

	objsDB, err := h.Repository.GetAttributes(c.Context(), params)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch attributes", nil)
	}

	resp := make([]models.Attribute, len(objsDB))
	for i, obj := range objsDB {
		var attributeOptions []models.AttributeOption
		if err := json.Unmarshal(obj.Options, &attributeOptions); err != nil {
			return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to parse attribute options", nil)
		}

		resp[i] = models.Attribute{
			ID:            obj.AttributeID,
			Title:         obj.Title,
			DataType:      obj.DataType,
			Unit:          obj.Unit.AttributeUnit,
			Required:      obj.Required,
			AppliesTo:     obj.AppliesTo,
			ProductTypeID: obj.ProductTypeID,
			Options:       attributeOptions,
		}
	}

	return api.SuccessResponse(c, fiber.StatusOK, resp, "Attributes fetched successfully")

}

// GetAttribute fetches a single attribute
// @Summary Fetch a single attribute
// @Description Fetch a single attribute
// @Tags Attributes
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param attribute_id path string true "Attribute ID"
// @Success 200 {object} models.SuccessResponse{data=models.Attribute} "Attribute fetched successfully"
// @Failure 404 {object} models.ErrorResponse "Attribute not found"
// @Failure 500 {object} models.ErrorResponse "Failed to fetch attribute"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/attributes/{attribute_id} [get]
func (h *Handler) GetAttribute(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	attributeIDStr := c.Params("attribute_id", "0")
	attributeID, _ := strconv.ParseInt(attributeIDStr, 10, 64)

	params := db.GetAttributeParams{
		AttributeID: attributeID,
		ShopID:      shopID,
	}

	objDB, err := h.Repository.GetAttribute(c.Context(), params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Attribute not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch attribute", nil)
	}

	var attributeOptions []models.AttributeOption
	if err := json.Unmarshal(objDB.Options, &attributeOptions); err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get products attribute", nil)
	}

	resp := models.Attribute{
		ID:            objDB.AttributeID,
		Title:         objDB.Title,
		DataType:      objDB.DataType,
		Unit:          objDB.Unit.AttributeUnit,
		Required:      objDB.Required,
		AppliesTo:     objDB.AppliesTo,
		ProductTypeID: objDB.ProductTypeID,
		Options:       attributeOptions,
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Attribute fetched successfully")
}

// UpdateAttribute updates an attribute
// @Summary Update an attribute
// @Description Update an attribute
// @Tags Attributes
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param attribute_id path string true "Attribute ID"
// @Param attribute body models.AttributeUpdateParams true "Attribute"
// @Success 200 {object} models.SuccessResponse{data=models.Attribute} "Attribute updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 404 {object} models.ErrorResponse "Attribute not found"
// @Failure 500 {object} models.ErrorResponse "Failed to update attribute"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/attributes/{attribute_id} [put]
func (h *Handler) UpdateAttribute(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	attributeIDStr := c.Params("attribute_id", "0")
	attributeID, _ := strconv.ParseInt(attributeIDStr, 10, 64)

	var attribute models.AttributeUpdateParams
	if err := c.BodyParser(&attribute); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validator := &models.XValidator{}
	if errs := validator.Validate(&attribute); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		fmt.Println(errMsgs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}

	param := db.UpdateAttributeParams{
		Title:       attribute.Title,
		DataType:    db.NullAttributeDataType{AttributeDataType: attribute.DataType, Valid: attribute.DataType != ""},
		Unit:        db.NullAttributeUnit{AttributeUnit: attribute.Unit, Valid: attribute.Unit != ""},
		Required:    attribute.Required,
		AppliesTo:   db.NullAttributeAppliesTo{AttributeAppliesTo: attribute.AppliesTo, Valid: attribute.AppliesTo != ""},
		AttributeID: attributeID,
		ShopID:      shopID,
	}

	objDB, err := h.Repository.UpdateAttribute(c.Context(), param)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Attribute not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update attribute", nil)
	}
	resp := models.Attribute{
		ID:            objDB.AttributeID,
		Title:         objDB.Title,
		DataType:      objDB.DataType,
		Unit:          objDB.Unit.AttributeUnit,
		Required:      objDB.Required,
		AppliesTo:     objDB.AppliesTo,
		ProductTypeID: objDB.ProductTypeID,
	}

	return api.SuccessResponse(c, fiber.StatusOK, resp, "Attribute updated successfully")
}

// DeleteAttribute deletes an attribute
// @Summary Delete an attribute
// @Description Delete an attribute
// @Tags Attributes
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param attribute_id path string true "Attribute ID"
// @Success 200 {object} models.SuccessResponse{data=models.Attribute} "Attribute deleted successfully"
// @Failure 404 {object} models.ErrorResponse "Attribute not found"
// @Failure 500 {object} models.ErrorResponse "Failed to delete attribute"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/attributes/{attribute_id} [delete]
func (h *Handler) DeleteAttribute(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	attributeIDStr := c.Params("attribute_id", "0")
	attributeID, _ := strconv.ParseInt(attributeIDStr, 10, 64)

	params := db.DeleteAttributeParams{
		AttributeID: attributeID,
		ShopID:      shopID,
	}

	objDB, err := h.Repository.DeleteAttribute(c.Context(), params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Attribute not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete attribute", nil)
	}

	resp := models.Attribute{
		ID:            objDB.AttributeID,
		Title:         objDB.Title,
		DataType:      objDB.DataType,
		Unit:          objDB.Unit.AttributeUnit,
		Required:      objDB.Required,
		AppliesTo:     objDB.AppliesTo,
		ProductTypeID: objDB.ProductTypeID,
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Attribute deleted successfully")
}

// CreateAttributeOption Create a new attribute option
// @Summary Create a new attribute option
// @Description Create a new attribute option
// @Tags Attributes
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param attribute_id path string true "Attribute ID"
// @Param option body models.AttributeOptionCreateParams true "Attribute Option"
// @Success 200 {object} models.SuccessResponse{data=models.AttributeOption} "Attribute option created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/attributes/{attribute_id}/options [post]
func (h *Handler) CreateAttributeOption(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	attributeIDStr := c.Params("attribute_id", "0")
	attributeID, _ := strconv.ParseInt(attributeIDStr, 10, 64)

	var option models.AttributeOptionCreateParams
	if err := c.BodyParser(&option); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validator := &models.XValidator{}
	if errs := validator.Validate(&option); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}

	param := db.CreateAttributeOptionParams{
		Value:       option.Value,
		ShopID:      shopID,
		AttributeID: attributeID,
	}

	objDB, err := h.Repository.CreateAttributeOption(c.Context(), param)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == errors.UniqueViolation {
				return api.ErrorResponse(c, fiber.StatusConflict, "attribute-option already exists", nil)
			}
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create attribute option", nil)
	}

	resp := models.AttributeOption{
		ID:          objDB.AttributeOptionID,
		Value:       objDB.Value,
		AttributeID: objDB.AttributeID,
	}
	return api.SuccessResponse(c, fiber.StatusCreated, resp, "Attribute option created successfully")
}

// GetAttributeOptions fetches all attribute options
// @Summary Fetch all attribute options
// @Description Fetch all attribute options
// @Tags Attributes
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param attribute_id path string true "Attribute ID"
// @Success 200 {object} models.SuccessResponse{data=[]models.AttributeOption} "Attribute options fetched successfully"
// @Failure 500 {object} models.ErrorResponse "Failed to fetch attribute options"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/attributes/{attribute_id}/options [get]
func (h *Handler) GetAttributeOptions(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	attributeIDStr := c.Params("attribute_id", "0")
	attributeID, _ := strconv.ParseInt(attributeIDStr, 10, 64)

	params := db.GetAttributeOptionsParams{
		AttributeID: attributeID,
		ShopID:      shopID,
	}

	objsDB, err := h.Repository.GetAttributeOptions(c.Context(), params)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch attribute options", nil)
	}

	resp := make([]models.AttributeOption, len(objsDB))
	for i, obj := range objsDB {
		resp[i] = models.AttributeOption{
			ID:          obj.AttributeOptionID,
			Value:       obj.Value,
			AttributeID: obj.AttributeID,
		}
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Attribute options fetched successfully")
}

// UpdateAttributeOption updates an attribute option
// @Summary Update an attribute option
// @Description Update an attribute option
// @Tags Attributes
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param attribute_option_id path string true "Attribute Option ID"
// @Param option body models.AttributeOptionUpdateParams true "Attribute Option"
// @Success 200 {object} models.SuccessResponse{data=models.AttributeOption} "Attribute option updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 404 {object} models.ErrorResponse "Attribute option not found"
// @Failure 500 {object} models.ErrorResponse "Failed to update attribute option"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/attribute-options/{attribute_option_id} [put]
func (h *Handler) UpdateAttributeOption(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	attributeOptionIDStr := c.Params("attribute_option_id", "0")
	attributeOptionID, _ := strconv.ParseInt(attributeOptionIDStr, 10, 64)

	var option models.AttributeOptionUpdateParams
	if err := c.BodyParser(&option); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validator := &models.XValidator{}
	if errs := validator.Validate(&option); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}

	param := db.UpdateAttributeOptionParams{
		Value:             option.Value,
		AttributeOptionID: attributeOptionID,
		ShopID:            shopID,
	}

	objDB, err := h.Repository.UpdateAttributeOption(c.Context(), param)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Attribute option not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update attribute option", nil)
	}

	resp := models.AttributeOption{
		ID:          objDB.AttributeOptionID,
		Value:       objDB.Value,
		AttributeID: objDB.AttributeID,
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Attribute option updated successfully")
}

// DeleteAttributeOption deletes an attribute option
// @Summary Delete an attribute option
// @Description Delete an attribute option
// @Tags Attributes
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param attribute_option_id path string true "Attribute Option ID"
// @Success 200 {object} models.SuccessResponse{data=models.AttributeOption} "Attribute option deleted successfully"
// @Failure 404 {object} models.ErrorResponse "Attribute option not found"
// @Failure 500 {object} models.ErrorResponse "Failed to delete attribute option"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/attribute-options/{attribute_option_id} [delete]
func (h *Handler) DeleteAttributeOption(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)
	attributeOptionIDStr := c.Params("attribute_option_id", "0")
	attributeOptionID, _ := strconv.ParseInt(attributeOptionIDStr, 10, 64)

	params := db.DeleteAttributeOptionParams{
		AttributeOptionID: attributeOptionID,
		ShopID:            shopID,
	}

	objDB, err := h.Repository.DeleteAttributeOption(c.Context(), params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Attribute option not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete attribute option", nil)
	}

	resp := models.AttributeOption{
		ID:          objDB.AttributeOptionID,
		Value:       objDB.Value,
		AttributeID: objDB.AttributeID,
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Attribute option deleted successfully")
}
