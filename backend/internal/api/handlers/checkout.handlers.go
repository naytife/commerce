package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"go.uber.org/zap"
)

// InitiateCheckout initiates the checkout process
// @Summary      Initiate checkout
// @Description  Begin the checkout process for cart items
// @Tags         checkout
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        checkout body models.CheckoutRequest true "Checkout details"
// @Success      200  {object}  models.SuccessResponse{data=models.CheckoutResponse}  "Checkout initiated successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid request"
// @Failure      404  {object}  models.ErrorResponse "Shop not found"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Router       /shops/{shop_id}/checkout [post]
func (h *Handler) InitiateCheckout(c *fiber.Ctx) error {
	shopID, err := strconv.ParseInt(c.Params("shop_id"), 10, 64)
	if err != nil {
		zap.L().Warn("InitiateCheckout: invalid shop id", zap.String("param", c.Params("shop_id")))
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	var req models.CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		zap.L().Warn("InitiateCheckout: failed to parse body", zap.Error(err), zap.Int64("shop_id", shopID))
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

	// Verify shop exists and get payment methods
	shop, err := h.Repository.GetShop(c.Context(), shopID)
	if err != nil {
		zap.L().Error("InitiateCheckout: shop not found", zap.Error(err), zap.Int64("shop_id", shopID))
		return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
	}

	// Get available payment methods for the shop
	paymentMethods, err := h.Repository.GetShopPaymentMethods(c.Context(), shopID)
	if err != nil {
		zap.L().Error("InitiateCheckout: failed to get payment methods", zap.Error(err), zap.Int64("shop_id", shopID))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get payment methods", nil)
	}

	enabledMethods := []models.PaymentMethodInfo{}
	for _, pm := range paymentMethods {
		if pm.IsEnabled {
			method := models.PaymentMethodInfo{
				Type:        string(pm.MethodType),
				Enabled:     true,
				DisplayName: getPaymentMethodDisplayName(pm.MethodType),
			}

			// Parse attributes for Stripe-specific config
			if pm.MethodType == db.PaymentMethodTypeStripe && len(pm.Attributes) > 0 {
				var attributes map[string]interface{}
				if err := json.Unmarshal(pm.Attributes, &attributes); err == nil {
					config := make(map[string]interface{})
					if publishableKey, ok := attributes["publishable_key"].(string); ok {
						config["publishable_key"] = publishableKey
					}
					if testMode, ok := attributes["test_mode"].(bool); ok {
						config["test_mode"] = testMode
					}
					method.Config = config
				}
			}

			enabledMethods = append(enabledMethods, method)
		}
	}

	response := models.CheckoutResponse{
		ShopID:          shopID,
		ShopName:        shop.Title,
		CurrencyCode:    shop.CurrencyCode,
		PaymentMethods:  enabledMethods,
		CustomerInfo:    req.CustomerInfo,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		ExpiresAt:       time.Now().Add(30 * time.Minute), // 30 minutes checkout session
	}

	return api.SuccessResponse(c, fiber.StatusOK, response, "Checkout initiated successfully")
}

// ProcessPayment processes payment for an order
// @Summary      Process payment
// @Description  Process payment for an order using the selected payment method
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        payment body models.PaymentRequest true "Payment details"
// @Success      200  {object}  models.SuccessResponse{data=models.PaymentResponse}  "Payment processed successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid request"
// @Failure      402  {object}  models.ErrorResponse "Payment failed"
// @Failure      404  {object}  models.ErrorResponse "Shop not found"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Router       /shops/{shop_id}/payment [post]
func (h *Handler) ProcessPayment(c *fiber.Ctx) error {
	shopID, err := strconv.ParseInt(c.Params("shop_id"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	var req models.PaymentRequest
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

	// Verify shop exists
	shop, err := h.Repository.GetShop(c.Context(), shopID)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
	}

	// Extract amount from payment details
	amount, ok := req.PaymentDetails["amount"].(float64)
	if !ok {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Missing amount in payment details", nil)
	}

	// Use PaymentProcessorFactory for processing payments
	if h.PaymentProcessorFactory != nil {
		return h.processPaymentWithFactory(c, shopID, shop, req, amount)
	}

	// If PaymentProcessorFactory is not available, return error
	return api.ErrorResponse(c, fiber.StatusInternalServerError, "Payment processing is not configured", nil)
}

// CreatePaymentIntent creates a payment intent for checkout
// @Summary      Create payment intent
// @Description  Create a payment intent for the checkout session
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        payment body models.PaymentIntentRequest true "Payment intent details"
// @Success      200  {object}  models.SuccessResponse{data=models.PaymentIntentResponse}  "Payment intent created successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid request"
// @Failure      404  {object}  models.ErrorResponse "Shop not found"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Router       /shops/{shop_id}/payment/intent [post]
func (h *Handler) CreatePaymentIntent(c *fiber.Ctx) error {
	shopID, err := strconv.ParseInt(c.Params("shop_id"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	var req models.PaymentIntentRequest
	if err := c.BodyParser(&req); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Verify shop exists
	_, err = h.Repository.GetShop(c.Context(), shopID)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
	}

	// Use PaymentProcessorFactory for creating payment intents
	if h.PaymentProcessorFactory != nil {
		processor := h.PaymentProcessorFactory.GetProcessor(req.PaymentMethod)
		if processor == nil {
			return api.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Unsupported payment method: %s", req.PaymentMethod), nil)
		}

		// Create payment intent with the selected processor
		response, err := processor.CreatePaymentIntent(c.Context(), shopID, req)
		if err != nil {
			return api.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to create payment intent: %v", err), nil)
		}

		return api.SuccessResponse(c, fiber.StatusOK, *response, "Payment intent created successfully")
	}

	// If PaymentProcessorFactory is not available, return error
	return api.ErrorResponse(c, fiber.StatusInternalServerError, "Payment processing is not configured", nil)
}

// processPaymentWithFactory handles payment processing using the PaymentProcessorFactory
func (h *Handler) processPaymentWithFactory(c *fiber.Ctx, shopID int64, shop db.Shop, req models.PaymentRequest, amount float64) error {
	// Get the appropriate payment processor
	processor := h.PaymentProcessorFactory.GetProcessor(req.PaymentMethod)
	if processor == nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Unsupported payment method: %s", req.PaymentMethod), nil)
	}

	// Process the payment using the selected processor
	response, err := processor.ProcessPayment(c.Context(), shopID, req, amount, shop.CurrencyCode)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusPaymentRequired, fmt.Sprintf("Payment failed: %v", err), nil)
	}

	// Create order if payment successful
	// TODO: Implement actual order creation from checkout session
	response.OrderID = 1

	return api.SuccessResponse(c, fiber.StatusOK, *response, "Payment processed successfully")
}

// processStripePayment handles Stripe payment processing (legacy)
// func (h *Handler) processStripePayment(c *fiber.Ctx, shopID int64, shop db.Shop, req models.PaymentRequest) error {
// 	// Check if Stripe service is available
// 	if h.StripeService == nil {
// 		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Stripe service not available", nil)
// 	}

// 	// Extract payment intent ID from payment details
// 	paymentIntentID, ok := req.PaymentDetails["payment_intent_id"].(string)
// 	if !ok {
// 		return api.ErrorResponse(c, fiber.StatusBadRequest, "Missing payment intent ID", nil)
// 	}

// 	// Confirm the payment with Stripe
// 	response, err := h.StripeService.ConfirmPayment(c.Context(), shopID, paymentIntentID)
// 	if err != nil {
// 		return api.ErrorResponse(c, fiber.StatusPaymentRequired, fmt.Sprintf("Payment failed: %v", err), nil)
// 	}

// 	// Create order if payment successful
// 	response.OrderID = 1 // TODO: Implement actual order creation from checkout session

// 	return api.SuccessResponse(c, fiber.StatusOK, *response, "Payment processed successfully")
// }

// getPaymentMethodDisplayName returns a user-friendly display name for payment methods
func getPaymentMethodDisplayName(methodType db.PaymentMethodType) string {
	switch methodType {
	case db.PaymentMethodTypeStripe:
		return "Credit/Debit Card"
	case db.PaymentMethodTypePaypal:
		return "PayPal"
	case db.PaymentMethodTypePaystack:
		return "Paystack"
	case db.PaymentMethodTypeFlutterwave:
		return "Flutterwave"
	default:
		return string(methodType)
	}
}
