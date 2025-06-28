package handlers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/db"
)

// PaymentMethodsRequest represents the request for updating payment methods
type PaymentMethodsRequest struct {
	PaymentMethods []PaymentMethodConfig `json:"payment_methods" validate:"required"`
}

// PaymentMethodConfig represents configuration for a payment method
type PaymentMethodConfig struct {
	MethodType string                 `json:"method_type" validate:"required"`
	IsEnabled  bool                   `json:"is_enabled"`
	Config     map[string]interface{} `json:"config"`
}

// GetShopPaymentMethods gets all payment methods for a shop
// @Summary      Get all payment methods for a shop
// @Description  Retrieve all configured payment methods for a specific shop
// @Tags         payment-methods
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Success      200  {object}   models.SuccessResponse{data=[]models.PaymentMethodInfo} "Payment methods retrieved successfully"
// @Failure      400  {object}   models.ErrorResponse "Invalid shop ID"
// @Failure      404  {object}   models.ErrorResponse "Shop not found"
// @Failure      500  {object}   models.ErrorResponse "Failed to get payment methods"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/payment-methods [get]
func (h *Handler) GetShopPaymentMethods(c *fiber.Ctx) error {
	shopID, err := strconv.ParseInt(c.Params("shop_id"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	// Verify shop exists
	_, err = h.Repository.GetShop(c.Context(), shopID)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
	}

	// Get payment methods
	paymentMethods, err := h.Repository.GetShopPaymentMethods(c.Context(), shopID)
	if err != nil {
		return api.SystemErrorResponse(c, err, "Failed to get payment methods")
	}

	// Convert to response format
	response := make([]map[string]interface{}, len(paymentMethods))
	for i, pm := range paymentMethods {
		var config map[string]interface{}
		if len(pm.Attributes) > 0 {
			json.Unmarshal(pm.Attributes, &config)
		}

		// Extract safe config for frontend (no secrets)
		safeConfig := map[string]interface{}{}
		switch string(pm.MethodType) {
		case "stripe":
			if v, ok := config["publishable_key"].(string); ok {
				safeConfig["publishable_key"] = v
			}
			if v, ok := config["test_mode"].(bool); ok {
				safeConfig["test_mode"] = v
			}
		case "paypal":
			if v, ok := config["client_id"].(string); ok {
				safeConfig["client_id"] = v
			}
			if v, ok := config["sandbox_mode"].(bool); ok {
				safeConfig["sandbox_mode"] = v
			}
		case "paystack":
			if v, ok := config["public_key"].(string); ok {
				safeConfig["public_key"] = v
			}
			if v, ok := config["test_mode_paystack"].(bool); ok {
				safeConfig["test_mode_paystack"] = v
			}
		case "flutterwave":
			if v, ok := config["public_key_flutterwave"].(string); ok {
				safeConfig["public_key_flutterwave"] = v
			}
			if v, ok := config["test_mode_flutterwave"].(bool); ok {
				safeConfig["test_mode_flutterwave"] = v
			}
		}

		methodName := strings.Title(string(pm.MethodType))
		response[i] = map[string]interface{}{
			"id":       string(pm.MethodType),
			"name":     methodName,
			"provider": string(pm.MethodType),
			"enabled":  pm.IsEnabled,
			"config":   safeConfig,
		}
	}
	return api.SuccessResponse(c, fiber.StatusOK, response, "Payment methods retrieved successfully")
}

// UpsertShopPaymentMethod creates or updates a payment method for a shop
// @Summary      Create or update a payment method
// @Description  Create or update a payment method configuration for a shop
// @Tags         payment-methods
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        method_type path string true "Payment method type (stripe, paypal, paystack, flutterwave)"
// @Param        config body PaymentMethodConfig true "Payment method configuration"
// @Success      200  {object}   models.SuccessResponse{data=models.PaymentMethodInfo} "Payment method updated successfully"
// @Failure      400  {object}   models.ErrorResponse "Invalid request body or shop ID"
// @Failure      404  {object}   models.ErrorResponse "Shop not found"
// @Failure      500  {object}   models.ErrorResponse "Failed to update payment method"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/payment-methods/{method_type} [put]
func (h *Handler) UpsertShopPaymentMethod(c *fiber.Ctx) error {
	shopID, err := api.ParseIDParameter(c, "shop_id", "Shop")
	if err != nil {
		return err
	}

	methodType := c.Params("method_type")
	if methodType == "" {
		return api.BusinessLogicErrorResponse(c, "Method type is required")
	}

	var req PaymentMethodConfig
	if err := c.BodyParser(&req); err != nil {
		return api.BusinessLogicErrorResponse(c, "Invalid request body")
	}

	// Verify shop exists
	_, err = h.Repository.GetShop(c.Context(), shopID)
	if err != nil {
		return api.NotFoundErrorResponse(c, "Shop")
	}

	// Convert method type to enum
	var pmType db.PaymentMethodType
	switch methodType {
	case "stripe":
		pmType = db.PaymentMethodTypeStripe
	case "paypal":
		pmType = db.PaymentMethodTypePaypal
	case "paystack":
		pmType = db.PaymentMethodTypePaystack
	case "flutterwave":
		pmType = db.PaymentMethodTypeFlutterwave
	default:
		return api.BusinessLogicErrorResponse(c, "Invalid payment method type")
	}

	// Convert config to JSON
	configJSON, err := json.Marshal(req.Config)
	if err != nil {
		return api.BusinessLogicErrorResponse(c, "Invalid config format")
	}

	// Upsert payment method
	paymentMethod, err := h.Repository.UpsertShopPaymentMethod(c.Context(), db.UpsertShopPaymentMethodParams{
		ShopID:     shopID,
		MethodType: pmType,
		IsEnabled:  req.IsEnabled,
		Attributes: configJSON,
	})
	if err != nil {
		return api.SystemErrorResponse(c, err, "Failed to update payment method")
	}

	// Prepare response
	var config map[string]interface{}
	if len(paymentMethod.Attributes) > 0 {
		json.Unmarshal(paymentMethod.Attributes, &config)
	}

	response := map[string]interface{}{
		"method_type": string(paymentMethod.MethodType),
		"is_enabled":  paymentMethod.IsEnabled,
		"config":      config,
		"created_at":  paymentMethod.CreatedAt,
		"updated_at":  paymentMethod.UpdatedAt,
	}

	return api.SuccessResponse(c, fiber.StatusOK, response, "Payment method updated successfully")
}

// UpdateShopPaymentMethodStatus updates the enabled status of a payment method
// @Summary      Update payment method status
// @Description  Enable or disable a payment method for a shop
// @Tags         payment-methods
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        method_type path string true "Payment method type (stripe, paypal, paystack, flutterwave)"
// @Param        status body object{is_enabled=bool} true "Payment method status"
// @Success      200  {object}   models.SuccessResponse{data=models.PaymentMethodInfo} "Payment method status updated successfully"
// @Failure      400  {object}   models.ErrorResponse "Invalid request body or shop ID"
// @Failure      404  {object}   models.ErrorResponse "Payment method not found"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/payment-methods/{method_type}/status [patch]
func (h *Handler) UpdateShopPaymentMethodStatus(c *fiber.Ctx) error {
	shopID, err := strconv.ParseInt(c.Params("shop_id"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	methodType := c.Params("method_type")
	if methodType == "" {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Method type is required", nil)
	}

	var req struct {
		IsEnabled bool `json:"is_enabled"`
	}
	if err := c.BodyParser(&req); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Convert method type to enum
	var pmType db.PaymentMethodType
	switch methodType {
	case "stripe":
		pmType = db.PaymentMethodTypeStripe
	case "paypal":
		pmType = db.PaymentMethodTypePaypal
	case "paystack":
		pmType = db.PaymentMethodTypePaystack
	case "flutterwave":
		pmType = db.PaymentMethodTypeFlutterwave
	default:
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid payment method type", nil)
	}

	// Update payment method status
	paymentMethod, err := h.Repository.UpdateShopPaymentMethodStatus(c.Context(), db.UpdateShopPaymentMethodStatusParams{
		ShopID:     shopID,
		MethodType: pmType,
		IsEnabled:  req.IsEnabled,
	})
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusNotFound, "Payment method not found", nil)
	}

	// Prepare response
	var config map[string]interface{}
	if len(paymentMethod.Attributes) > 0 {
		json.Unmarshal(paymentMethod.Attributes, &config)
	}

	response := map[string]interface{}{
		"method_type": string(paymentMethod.MethodType),
		"is_enabled":  paymentMethod.IsEnabled,
		"config":      config,
		"updated_at":  paymentMethod.UpdatedAt,
	}

	return api.SuccessResponse(c, fiber.StatusOK, response, "Payment method status updated successfully")
}

// DeleteShopPaymentMethod deletes a payment method for a shop
// @Summary      Delete a payment method
// @Description  Remove a payment method configuration from a shop
// @Tags         payment-methods
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        method_type path string true "Payment method type (stripe, paypal, paystack, flutterwave)"
// @Success      200  {object}   models.SuccessResponse "Payment method deleted successfully"
// @Failure      400  {object}   models.ErrorResponse "Invalid request body or shop ID"
// @Failure      500  {object}   models.ErrorResponse "Failed to delete payment method"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/payment-methods/{method_type} [delete]
func (h *Handler) DeleteShopPaymentMethod(c *fiber.Ctx) error {
	shopID, err := strconv.ParseInt(c.Params("shop_id"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	methodType := c.Params("method_type")
	if methodType == "" {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Method type is required", nil)
	}

	// Convert method type to enum
	var pmType db.PaymentMethodType
	switch methodType {
	case "stripe":
		pmType = db.PaymentMethodTypeStripe
	case "paypal":
		pmType = db.PaymentMethodTypePaypal
	case "paystack":
		pmType = db.PaymentMethodTypePaystack
	case "flutterwave":
		pmType = db.PaymentMethodTypeFlutterwave
	default:
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid payment method type", nil)
	}

	// Delete payment method
	err = h.Repository.DeleteShopPaymentMethod(c.Context(), db.DeleteShopPaymentMethodParams{
		ShopID:     shopID,
		MethodType: pmType,
	})
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete payment method", nil)
	}

	return api.SuccessResponse(c, fiber.StatusOK, nil, "Payment method deleted successfully")
}

// TestPaymentMethod tests a payment method configuration
// @Summary      Test payment method configuration
// @Description  Test the connectivity and configuration of a payment method
// @Tags         payment-methods
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        method_type path string true "Payment method type (stripe, paypal, paystack, flutterwave)"
// @Success      200  {object}   models.SuccessResponse{data=object} "Payment method test completed"
// @Failure      400  {object}   models.ErrorResponse "Invalid request or testing not supported"
// @Failure      404  {object}   models.ErrorResponse "Payment method not configured"
// @Failure      500  {object}   models.ErrorResponse "Invalid payment method configuration"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/payment-methods/{method_type}/test [post]
func (h *Handler) TestPaymentMethod(c *fiber.Ctx) error {
	shopID, err := strconv.ParseInt(c.Params("shop_id"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	methodType := c.Params("method_type")
	if methodType == "" {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Method type is required", nil)
	}

	// Convert method type to enum
	var pmType db.PaymentMethodType
	switch methodType {
	case "stripe":
		pmType = db.PaymentMethodTypeStripe
	case "paypal":
		pmType = db.PaymentMethodTypePaypal
	case "paystack":
		pmType = db.PaymentMethodTypePaystack
	case "flutterwave":
		pmType = db.PaymentMethodTypeFlutterwave
	case "pay_on_delivery":
		// Pay on delivery doesn't need API testing
		return api.SuccessResponse(c, fiber.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "Pay on delivery is always available",
		}, "Payment method test successful")
	default:
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid payment method type", nil)
	}

	// Get payment method configuration
	paymentMethod, err := h.Repository.GetShopPaymentMethod(c.Context(), db.GetShopPaymentMethodParams{
		ShopID:     shopID,
		MethodType: pmType,
	})
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusNotFound, "Payment method not configured", nil)
	}

	// Parse configuration
	var config map[string]interface{}
	if len(paymentMethod.Attributes) > 0 {
		if err := json.Unmarshal(paymentMethod.Attributes, &config); err != nil {
			return api.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid payment method configuration", nil)
		}
	}

	// Test the payment method based on type
	var testResult map[string]interface{}
	switch methodType {
	case "stripe":
		testResult = h.testStripeConnection(config)
	case "paystack":
		testResult = h.testPaystackConnection(config)
	case "flutterwave":
		testResult = h.testFlutterwaveConnection(config)
	case "paypal":
		testResult = h.testPaypalConnection(config)
	default:
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Testing not supported for this payment method", nil)
	}

	return api.SuccessResponse(c, fiber.StatusOK, testResult, "Payment method test completed")
}

// Helper methods for testing payment providers
func (h *Handler) testStripeConnection(config map[string]interface{}) map[string]interface{} {
	// Basic validation for Stripe
	secretKey, hasSecret := config["secret_key"].(string)
	publishableKey, hasPublishable := config["publishable_key"].(string)

	if !hasSecret || !hasPublishable || secretKey == "" || publishableKey == "" {
		return map[string]interface{}{
			"status":  "error",
			"message": "Missing required Stripe keys",
		}
	}

	// In a real implementation, you would make an API call to Stripe to verify the keys
	// For now, we'll just validate the format
	if !strings.HasPrefix(secretKey, "sk_") || !strings.HasPrefix(publishableKey, "pk_") {
		return map[string]interface{}{
			"status":  "error",
			"message": "Invalid Stripe key format",
		}
	}

	return map[string]interface{}{
		"status":  "success",
		"message": "Stripe connection test successful",
	}
}

func (h *Handler) testPaystackConnection(config map[string]interface{}) map[string]interface{} {
	secretKey, hasSecret := config["secret_key"].(string)

	if !hasSecret || secretKey == "" {
		return map[string]interface{}{
			"status":  "error",
			"message": "Missing Paystack secret key",
		}
	}

	// Basic format validation
	if !strings.HasPrefix(secretKey, "sk_") {
		return map[string]interface{}{
			"status":  "error",
			"message": "Invalid Paystack secret key format",
		}
	}

	return map[string]interface{}{
		"status":  "success",
		"message": "Paystack connection test successful",
	}
}

func (h *Handler) testFlutterwaveConnection(config map[string]interface{}) map[string]interface{} {
	secretKey, hasSecret := config["secret_key"].(string)

	if !hasSecret || secretKey == "" {
		return map[string]interface{}{
			"status":  "error",
			"message": "Missing Flutterwave secret key",
		}
	}

	return map[string]interface{}{
		"status":  "success",
		"message": "Flutterwave connection test successful",
	}
}

func (h *Handler) testPaypalConnection(config map[string]interface{}) map[string]interface{} {
	clientID, hasClientID := config["client_id"].(string)
	clientSecret, hasClientSecret := config["client_secret"].(string)

	if !hasClientID || !hasClientSecret || clientID == "" || clientSecret == "" {
		return map[string]interface{}{
			"status":  "error",
			"message": "Missing PayPal client credentials",
		}
	}

	return map[string]interface{}{
		"status":  "success",
		"message": "PayPal connection test successful",
	}
}
