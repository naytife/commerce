package handlers

import "github.com/gofiber/fiber/v2"

// CreateProductType Creeate a new product type
// @Summary Create a new product type
// @Description Create a new product type
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param productType body models.ProductTypeCreate true "Product type object that needs to be created"
// @Success 200 {object} models.SuccessResponse{data=models.ProductType} "Product type created successfully"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/product-types [post]
func (h *Handler) CreateProductType(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}

// GetProductTypes fetches all product types
// @Summary Fetch all product types
// @Description Fetch all product types
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Success 200 {object} models.SuccessResponse{data=models.ProductType} "Product types fetched successfully"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/product-types [get]
func (h *Handler) GetProductTypes(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}

// GetProductType fetches a single product type
// @Summary Fetch a single product type
// @Description Fetch a single product type
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_type_id path string true "Product type ID"
// @Success 200 {object} models.SuccessResponse{data=models.ProductType} "Product type fetched successfully"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/product-types/{product_type_id} [get]
func (h *Handler) GetProductType(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}

// UpdateProductType updates a product type
// @Summary Update a product type
// @Description Update a product type
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_type_id path string true "Product type ID"
// @Param productType body models.ProductTypeUpdate true "Product type object that needs to be updated"
// @Success 200 {object} models.SuccessResponse{data=models.ProductType} "Product type updated successfully"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/product-types/{product_type_id} [put]
func (h *Handler) UpdateProductType(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}

// DeleteProductType deletes a product type
// @Summary Delete a product type
// @Description Delete a product type
// @Tags ProductType
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Param product_type_id path string true "Product type ID"
// @Success 200 {object} models.SuccessResponse{data=models.ProductType} "Product type deleted successfully"
// @Security OAuth2AccessCode
// @Router /shops/{shop_id}/product-types/{product_type_id} [delete]
func (h *Handler) DeleteProductType(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}
