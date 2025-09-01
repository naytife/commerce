package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/services"
)

// GetPredefinedProductTypes returns all available predefined product type templates
// @Summary Get predefined product type templates
// @Description Get all available predefined product type templates
// @Tags ProductType
// @Produce json
// @Success 200 {object} models.SuccessResponse{data=[]models.PredefinedProductType} "Predefined product types fetched successfully"
// @Security OAuth2AccessCode
// @Router /predefined-product-types [get]
func (h *Handler) GetPredefinedProductTypes(c *fiber.Ctx) error {
	service := services.NewPredefinedProductTypeService()
	templates := service.GetPredefinedProductTypes()

	return api.SuccessResponse(c, fiber.StatusOK, templates, "Predefined product types fetched successfully")
}

// GetPredefinedProductType returns a specific predefined product type template
// @Summary Get a specific predefined product type template
// @Description Get a specific predefined product type template by ID
// @Tags ProductType
// @Produce json
// @Param template_id path string true "Template ID"
// @Success 200 {object} models.SuccessResponse{data=models.PredefinedProductType} "Predefined product type fetched successfully"
// @Failure 404 {object} models.ErrorResponse "Template not found"
// @Security OAuth2AccessCode
// @Router /predefined-product-types/{template_id} [get]
func (h *Handler) GetPredefinedProductType(c *fiber.Ctx) error {
	templateID := c.Params("template_id")

	service := services.NewPredefinedProductTypeService()
	templates := service.GetPredefinedProductTypes()

	for _, template := range templates {
		if template.ID == templateID {
			return api.SuccessResponse(c, fiber.StatusOK, template, "Predefined product type fetched successfully")
		}
	}

	return api.ErrorResponse(c, fiber.StatusNotFound, "Template not found", nil)
}

// CreateProductTypeFromTemplate creates a new product type from a predefined template
// @Summary Create product type from template
// @Description Create a new product type and its attributes from a predefined template
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param request body models.CreateProductTypeFromTemplateParams true "Template creation request"
// @Success 201 {object} models.SuccessResponse{data=models.ProductTypeWithTemplateResponse} "Product type created from template successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 404 {object} models.ErrorResponse "Template not found"
// @Failure 500 {object} models.ErrorResponse "Failed to create product type"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/product-types/from-template [post]
func (h *Handler) CreateProductTypeFromTemplate(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, _ := strconv.ParseInt(shopIDStr, 10, 64)

	var request models.CreateProductTypeFromTemplateParams
	if err := c.BodyParser(&request); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validator := &models.XValidator{}
	if errs := validator.Validate(&request); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}

	// Get the predefined template
	service := services.NewPredefinedProductTypeService()
	templates := service.GetPredefinedProductTypes()

	var selectedTemplate *models.PredefinedProductType
	for _, template := range templates {
		if template.ID == request.TemplateID {
			selectedTemplate = &template
			break
		}
	}

	if selectedTemplate == nil {
		return api.ErrorResponse(c, fiber.StatusNotFound, "Template not found", nil)
	}

	// Execute transaction to create product type and attributes
	var createdProductType models.ProductType
	var createdAttributes []models.Attribute

	err := h.Repository.WithTx(c.Context(), func(q *db.Queries) error {
		// Create the product type
		createProductTypeParams := db.CreateProductTypeParams{
			Title:        selectedTemplate.Title,
			Shippable:    selectedTemplate.Shippable,
			Digital:      selectedTemplate.Digital,
			SkuSubstring: &selectedTemplate.SkuSubstring,
			ShopID:       shopID,
		}

		productTypeDB, err := q.CreateProductType(c.Context(), createProductTypeParams)
		if err != nil {
			return err
		}

		createdProductType = models.ProductType{
			ID:           productTypeDB.ProductTypeID,
			Title:        productTypeDB.Title,
			Shippable:    productTypeDB.Shippable,
			Digital:      productTypeDB.Digital,
			SkuSubstring: productTypeDB.SkuSubstring,
		}

		// Create attributes for this product type
		for _, attrTemplate := range selectedTemplate.Attributes {
			// Map data type
			var dataType db.AttributeDataType
			switch attrTemplate.DataType {
			case "Text":
				dataType = db.AttributeDataTypeText
			case "Number":
				dataType = db.AttributeDataTypeNumber
			case "Date":
				dataType = db.AttributeDataTypeDate
			case "Option":
				dataType = db.AttributeDataTypeOption
			case "Color":
				dataType = db.AttributeDataTypeColor
			default:
				dataType = db.AttributeDataTypeText
			}

			// Map applies to
			var appliesTo db.AttributeAppliesTo
			if attrTemplate.AppliesTo == "ProductVariation" {
				appliesTo = db.AttributeAppliesToProductVariation
			} else {
				appliesTo = db.AttributeAppliesToProduct
			}

			// Map unit if provided
			var unit db.NullAttributeUnit
			if attrTemplate.Unit != nil {
				switch *attrTemplate.Unit {
				case "KG":
					unit = db.NullAttributeUnit{AttributeUnit: db.AttributeUnitKG, Valid: true}
				case "GB":
					unit = db.NullAttributeUnit{AttributeUnit: db.AttributeUnitGB, Valid: true}
				case "INCH":
					unit = db.NullAttributeUnit{AttributeUnit: db.AttributeUnitINCH, Valid: true}
				}
			}

			// Create attribute
			createAttrParams := db.CreateAttributeParams{
				Title:         attrTemplate.Title,
				DataType:      dataType,
				Unit:          unit,
				Required:      attrTemplate.Required,
				AppliesTo:     appliesTo,
				ProductTypeID: productTypeDB.ProductTypeID,
				ShopID:        shopID,
			}

			attributeDB, err := q.CreateAttribute(c.Context(), createAttrParams)
			if err != nil {
				return err
			}

			// Create attribute options if any
			var options []models.AttributeOption
			if len(attrTemplate.Options) > 0 {
				optionParams := make([]db.BatchUpsertAttributeOptionParams, len(attrTemplate.Options))
				for i, optionTemplate := range attrTemplate.Options {
					optionParams[i] = db.BatchUpsertAttributeOptionParams{
						Value:       optionTemplate.Value,
						AttributeID: attributeDB.AttributeID,
						ShopID:      shopID,
					}
				}

				batch := q.BatchUpsertAttributeOption(c.Context(), optionParams)
				var batchErr error
				var optionsDB []db.AttributeOption
				batch.Query(func(i int, optionsResult []db.AttributeOption, err error) {
					if err != nil {
						batchErr = err
					}
					if len(optionsResult) > 0 {
						optionsDB = append(optionsDB, optionsResult...)
					}
				})

				if batchErr != nil {
					return batchErr
				}

				if err := batch.Close(); err != nil {
					return err
				}

				for _, optionDB := range optionsDB {
					options = append(options, models.AttributeOption{
						ID:          optionDB.AttributeOptionID,
						Value:       optionDB.Value,
						AttributeID: optionDB.AttributeID,
					})
				}
			}

			createdAttributes = append(createdAttributes, models.Attribute{
				ID:            attributeDB.AttributeID,
				Title:         attributeDB.Title,
				DataType:      attributeDB.DataType,
				Unit:          attributeDB.Unit.AttributeUnit,
				Required:      attributeDB.Required,
				AppliesTo:     attributeDB.AppliesTo,
				ProductTypeID: attributeDB.ProductTypeID,
				Options:       options,
			})
		}

		return nil
	})

	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create product type from template", nil)
	}

	response := models.ProductTypeWithTemplateResponse{
		ProductType: createdProductType,
		Attributes:  createdAttributes,
	}

	return api.SuccessResponse(c, fiber.StatusCreated, response, "Product type created from template successfully")
}

// helper functions removed: attribute unit conversion moved/unused in this handler
