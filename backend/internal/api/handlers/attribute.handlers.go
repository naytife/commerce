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
	if attribute.AppliesTo == "ProductVariation" && len(attribute.Options) == 0 {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "At least one option is required for Variants Attribute", nil)
	}

	for _, option := range attribute.Options {
		if errs := validator.Validate(&option); len(errs) > 0 {
			errMsgs := models.FormatValidationErrors(errs)
			return &fiber.Error{
				Code:    fiber.ErrBadRequest.Code,
				Message: errMsgs,
			}
		}
	}
	var attribteDB db.Attribute
	attributeOptions := make([]db.AttributeOption, 0)
	err := h.Repository.WithTx(c.Context(), func(tx *db.Queries) error {
		param := db.CreateAttributeParams{
			Title:         attribute.Title,
			DataType:      attribute.DataType,
			Unit:          db.NullAttributeUnit{AttributeUnit: attribute.Unit, Valid: attribute.Unit != ""},
			Required:      attribute.Required,
			AppliesTo:     attribute.AppliesTo,
			ProductTypeID: productTypeID,
			ShopID:        shopID,
		}

		objDB, err := tx.CreateAttribute(c.Context(), param)
		if err != nil {
			return err
		}
		attribteDB = objDB

		if len(attribute.Options) > 0 {
			optionParams := make([]db.BatchUpsertAttributeOptionParams, len(attribute.Options))
			for i, option := range attribute.Options {
				optionParams[i] = db.BatchUpsertAttributeOptionParams{
					Value:       option.Value,
					ShopID:      shopID,
					AttributeID: objDB.AttributeID,
				}
			}
			batch := tx.BatchUpsertAttributeOption(c.Context(), optionParams)
			var batchErr error
			batch.Query(func(i int, options []db.AttributeOption, err error) {
				if err != nil {
					batchErr = err
				}
				if len(options) > 0 {
					attributeOptions = append(attributeOptions, options...)
				}
			})

			if batchErr != nil {
				return batchErr
			}

			if err := batch.Close(); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == errors.UniqueViolation {
				return api.ErrorResponse(c, fiber.StatusConflict, "Attribute already exists", nil)
			}
			return api.ErrorResponse(c, fiber.StatusConflict, fmt.Sprintf("Unique constraint violation: %s", pgErr.ConstraintName), nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create attribute", nil)
	}
	var options []models.AttributeOption
	for _, option := range attributeOptions {
		options = append(options, models.AttributeOption{
			ID:          option.AttributeOptionID,
			Value:       option.Value,
			AttributeID: option.AttributeID,
		})
	}
	resp := models.Attribute{
		ID:            attribteDB.AttributeID,
		Title:         attribteDB.Title,
		DataType:      attribteDB.DataType,
		Unit:          attribteDB.Unit.AttributeUnit,
		Required:      attribteDB.Required,
		AppliesTo:     attribteDB.AppliesTo,
		ProductTypeID: attribteDB.ProductTypeID,
		Options:       options,
	}
	return api.SuccessResponse(c, fiber.StatusCreated, resp, "Attribute created successfully")
}

// GetAttributes fetches all attributes
// @Summary Fetch all attributes
// @Description Fetch all attributes
// @Tags ProductType
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

	if attribute.AppliesTo == "ProductVariation" && len(attribute.Options) == 0 {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "At least one option is required for Variants Attribute", nil)
	}
	var attributeDB db.Attribute
	attributeOptions := make([]db.AttributeOption, 0)
	err := h.Repository.WithTx(c.Context(), func(tx *db.Queries) error {
		// Update the attribute
		param := db.UpdateAttributeParams{
			Title:       attribute.Title,
			DataType:    db.NullAttributeDataType{AttributeDataType: attribute.DataType, Valid: attribute.DataType != ""},
			Unit:        db.NullAttributeUnit{AttributeUnit: attribute.Unit, Valid: attribute.Unit != ""},
			Required:    attribute.Required,
			AppliesTo:   db.NullAttributeAppliesTo{AttributeAppliesTo: attribute.AppliesTo, Valid: attribute.AppliesTo != ""},
			AttributeID: attributeID,
			ShopID:      shopID,
		}

		objDB, err := tx.UpdateAttribute(c.Context(), param)
		if err != nil {
			return err
		}
		attributeDB = objDB

		// Handle options updates if provided
		if len(attribute.Options) > 0 {
			// Get option IDs that should be kept
			keepOptionIDs := make([]int64, len(attribute.Options))
			for i, option := range attribute.Options {
				keepOptionIDs[i] = option.ID
			}

			// Delete options not in the update list
			deleteBatch := tx.BatchDeleteAttributeOptions(c.Context(), []db.BatchDeleteAttributeOptionsParams{
				{
					ShopID:             shopID,
					AttributeID:        attributeID,
					AttributeOptionIds: keepOptionIDs,
				},
			})
			deleteBatch.Exec(func(_ int, err error) {
				if err != nil {
					return
				}
			})

			// Prepare options for batch upsert
			optionParams := make([]db.BatchUpsertAttributeOptionParams, len(attribute.Options))
			for i, option := range attribute.Options {
				optionParams[i] = db.BatchUpsertAttributeOptionParams{
					Value:       *option.Value,
					ShopID:      shopID,
					AttributeID: attributeID,
				}
			}

			// Perform batch upsert for options
			batch := tx.BatchUpsertAttributeOption(c.Context(), optionParams)
			var batchErr error
			batch.Query(func(i int, options []db.AttributeOption, err error) {
				if err != nil {
					batchErr = err
				}
				if len(options) > 0 {
					attributeOptions = append(attributeOptions, options...)
				}
			})

			if batchErr != nil {
				return batchErr
			}

			if err := batch.Close(); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == errors.UniqueViolation {
				return api.ErrorResponse(c, fiber.StatusConflict, "Attribute already exists", nil)
			}
			return api.ErrorResponse(c, fiber.StatusConflict, fmt.Sprintf("Unique constraint violation: %s", pgErr.ConstraintName), nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create attribute", nil)
	}
	var options []models.AttributeOption
	for _, option := range attributeOptions {
		options = append(options, models.AttributeOption{
			ID:          option.AttributeOptionID,
			Value:       option.Value,
			AttributeID: option.AttributeID,
		})
	}
	resp := models.Attribute{
		ID:            attributeDB.AttributeID,
		Title:         attributeDB.Title,
		DataType:      attributeDB.DataType,
		Unit:          attributeDB.Unit.AttributeUnit,
		Required:      attributeDB.Required,
		AppliesTo:     attributeDB.AppliesTo,
		ProductTypeID: attributeDB.ProductTypeID,
		Options:       options,
	}

	return api.SuccessResponse(c, fiber.StatusOK, resp, "Attribute updated successfully")
}

// DeleteAttribute deletes an attribute
// @Summary Delete an attribute
// @Description Delete an attribute
// @Tags Attributes
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

	err := h.Repository.DeleteAttribute(c.Context(), params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Attribute not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete attribute", nil)
	}

	return api.SuccessResponse(c, fiber.StatusOK, nil, "Attribute deleted successfully")
}
