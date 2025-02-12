package handlers

import "github.com/gofiber/fiber/v2"

// CreateAttribute Create a new attribute
// @Summary Create a new attribute
// @Description Create a new attribute
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_type_id path string true "Product Type ID"
// @Success 200 {object} models.SuccessResponse{data=models.Attribute} "Attribute created successfully"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/product-types/{product_type_id}/attributes [post]
func (h *Handler) CreateAttribute(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}

// CreateAttributeOption Create a new attribute option
// @Summary Create a new attribute option
// @Description Create a new attribute option
// @Tags Attributes
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param attribute_id path string true "Attribute ID"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/attributes/{attribute_id}/options [post]
func (h *Handler) CreateAttributeOption(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}
