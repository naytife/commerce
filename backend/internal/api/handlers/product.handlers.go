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
// @Router /shops/{shop_id}/product-types/{product_type_id}/product [post]
func (h *Handler) CreateProduct(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}
