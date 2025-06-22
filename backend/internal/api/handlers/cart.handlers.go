package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
)

// GetCart retrieves cart items for a session/user
// @Summary      Get cart items
// @Description  Get cart items for the current session
// @Tags         cart
// @Produce      json
// @Param        session_id query string false "Session ID for anonymous users"
// @Success      200  {object}  models.SuccessResponse{data=models.Cart}  "Cart retrieved successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid request"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Router       /cart [get]
func (h *Handler) GetCart(c *fiber.Ctx) error {
	// For now, return empty cart since we're using client-side storage
	// This can be extended later to support server-side cart storage
	cart := models.Cart{
		Items: []models.CartItem{},
		Total: 0,
	}

	return api.SuccessResponse(c, fiber.StatusOK, cart, "Cart retrieved successfully")
}

// AddToCart adds an item to the cart
// @Summary      Add item to cart
// @Description  Add a product variant to the cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        item body models.AddToCartRequest true "Cart item to add"
// @Success      200  {object}  models.SuccessResponse{data=models.CartItem}  "Item added to cart successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid request body"
// @Failure      404  {object}  models.ErrorResponse "Product not found"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Router       /cart/items [post]
func (h *Handler) AddToCart(c *fiber.Ctx) error {
	var req models.AddToCartRequest
	if err := c.BodyParser(&req); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Validate request
	validator := &models.XValidator{}
	if errs := validator.Validate(&req); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}

	// Get product variation details to validate and get current price
	// For now, we'll create a mock response since this is client-side cart
	cartItem := models.CartItem{
		ID:                 req.ProductVariationID,
		ProductVariationID: req.ProductVariationID,
		Quantity:           req.Quantity,
		Title:              req.Title,
		Price:              req.Price,
		Image:              req.Image,
		Slug:               req.Slug,
	}

	return api.SuccessResponse(c, fiber.StatusOK, cartItem, "Item added to cart successfully")
}

// UpdateCartItem updates quantity of an item in cart
// @Summary      Update cart item quantity
// @Description  Update the quantity of an item in the cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        item_id path string true "Cart Item ID"
// @Param        update body models.UpdateCartItemRequest true "Updated quantity"
// @Success      200  {object}  models.SuccessResponse{data=models.CartItem}  "Cart item updated successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid request"
// @Failure      404  {object}  models.ErrorResponse "Cart item not found"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Router       /cart/items/{item_id} [put]
func (h *Handler) UpdateCartItem(c *fiber.Ctx) error {
	itemID := c.Params("item_id")
	if itemID == "" {
		return api.BusinessLogicErrorResponse(c, "Item ID is required")
	}

	var req models.UpdateCartItemRequest
	if err := c.BodyParser(&req); err != nil {
		return api.BusinessLogicErrorResponse(c, "Invalid request body")
	}

	// Validate request
	if err := api.ValidateRequest(c, &req); err != nil {
		return err
	}

	// Since we're using client-side cart, return success response
	cartItem := models.CartItem{
		ID:       itemID,
		Quantity: req.Quantity,
	}

	return api.SuccessResponse(c, fiber.StatusOK, cartItem, "Cart item updated successfully")
}

// RemoveFromCart removes an item from the cart
// @Summary      Remove item from cart
// @Description  Remove an item from the cart
// @Tags         cart
// @Produce      json
// @Param        item_id path string true "Cart Item ID"
// @Success      200  {object}  models.SuccessResponse  "Item removed from cart successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid item ID"
// @Failure      404  {object}  models.ErrorResponse "Cart item not found"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Router       /cart/items/{item_id} [delete]
func (h *Handler) RemoveFromCart(c *fiber.Ctx) error {
	itemID := c.Params("item_id")
	if itemID == "" {
		return api.BusinessLogicErrorResponse(c, "Item ID is required")
	}

	// Since we're using client-side cart, return success response
	return api.SuccessResponse(c, fiber.StatusOK, nil, "Item removed from cart successfully")
}

// ClearCart removes all items from the cart
// @Summary      Clear cart
// @Description  Remove all items from the cart
// @Tags         cart
// @Produce      json
// @Param        session_id query string false "Session ID for anonymous users"
// @Success      200  {object}  models.SuccessResponse  "Cart cleared successfully"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Router       /cart [delete]
func (h *Handler) ClearCart(c *fiber.Ctx) error {
	// Since we're using client-side cart, return success response
	return api.SuccessResponse(c, fiber.StatusOK, nil, "Cart cleared successfully")
}
