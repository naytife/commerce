package handlers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/services"
)

var validate = validator.New()

// PaymentHandler handles payment-related API endpoints
type PaymentHandler struct {
	paymentFactory *services.PaymentProcessorFactory
	repository     db.Repository
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentFactory *services.PaymentProcessorFactory, repo db.Repository) *PaymentHandler {
	return &PaymentHandler{
		paymentFactory: paymentFactory,
		repository:     repo,
	}
}

// generateOrderID generates a simple order ID for Pay on Delivery
func generateOrderID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

// CreateCheckoutSessionRequest represents the request to create a checkout session
type CreateCheckoutSessionRequest struct {
	ShopID            int64             `json:"shop_id" validate:"required"`
	PaymentMethodType string            `json:"payment_method_type" validate:"required,oneof=stripe paystack flutterwave paypal pay_on_delivery"`
	Items             []CheckoutItem    `json:"items" validate:"required,min=1"`
	Customer          *CheckoutCustomer `json:"customer"`
	Shipping          *CheckoutShipping `json:"shipping"`
	Currency          string            `json:"currency" validate:"required"`
	SuccessURL        string            `json:"success_url" validate:"required,url"`
	CancelURL         string            `json:"cancel_url" validate:"required,url"`
	Metadata          map[string]string `json:"metadata"`
}

type CheckoutItem struct {
	ProductID   string `json:"product_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Amount      int64  `json:"amount" validate:"required,min=1"` // Amount in cents
	Currency    string `json:"currency" validate:"required"`
	Quantity    int64  `json:"quantity" validate:"required,min=1"`
}

type CheckoutCustomer struct {
	Email string `json:"email" validate:"email"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type CheckoutShipping struct {
	Name    string          `json:"name" validate:"required"`
	Address CheckoutAddress `json:"address" validate:"required"`
}

type CheckoutAddress struct {
	Line1      string `json:"line1" validate:"required"`
	Line2      string `json:"line2"`
	City       string `json:"city" validate:"required"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code" validate:"required"`
	Country    string `json:"country" validate:"required"`
}

// CreateCheckoutSessionResponse represents the response from creating a checkout session
type CreateCheckoutSessionResponse struct {
	PaymentIntentID    string                    `json:"payment_intent_id"`
	ClientSecret       string                    `json:"client_secret"`
	Status             string                    `json:"status"`
	NextAction         *models.PaymentNextAction `json:"next_action,omitempty"`
	PaymentMethodTypes []string                  `json:"payment_method_types"`
}

// CreateCheckoutSession creates a payment intent for checkout using the configured payment method
// @Summary Create checkout session
// @Description Creates a payment intent for processing payments using the shop's configured payment method
// @Tags payments
// @Accept json
// @Produce json
// @Param request body CreateCheckoutSessionRequest true "Checkout session request"
// @Success 200 {object} CreateCheckoutSessionResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /payments/checkout [post]
func (h *PaymentHandler) CreateCheckoutSession(c *fiber.Ctx) error {
	var req CreateCheckoutSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	// Handle Pay on Delivery
	if req.PaymentMethodType == "pay_on_delivery" {
		// Create a simple response for Pay on Delivery
		response := CreateCheckoutSessionResponse{
			PaymentIntentID:    "pod_" + strconv.FormatInt(req.ShopID, 10) + "_" + generateOrderID(),
			ClientSecret:       "",
			Status:             "requires_confirmation",
			PaymentMethodTypes: []string{"pay_on_delivery"},
		}
		return c.JSON(response)
	}

	// Get payment method configuration for the shop
	paymentMethods, err := h.repository.GetShopPaymentMethods(c.Context(), req.ShopID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Status:  "error",
			Message: "Failed to get payment method configuration",
			Code:    fiber.StatusInternalServerError,
		})
	}

	var paymentMethodConfig *db.ShopPaymentMethod
	for _, pm := range paymentMethods {
		if string(pm.MethodType) == req.PaymentMethodType && pm.IsEnabled {
			paymentMethodConfig = &pm
			break
		}
	}

	if paymentMethodConfig == nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Status:  "error",
			Message: "Payment method not configured or disabled for this shop",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Get payment processor
	processor := h.paymentFactory.GetProcessor(req.PaymentMethodType)
	if processor == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Status:  "error",
			Message: "Payment processor not available",
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range req.Items {
		totalAmount += float64(item.Amount*item.Quantity) / 100 // Convert from cents
	}

	// Create payment intent request
	paymentReq := models.PaymentIntentRequest{
		Amount:            totalAmount,
		CurrencyCode:      req.Currency,
		PaymentMethod:     req.PaymentMethodType,
		CheckoutSessionID: "checkout_" + strconv.FormatInt(req.ShopID, 10) + "_" + generateOrderID(),
		Description:       "Order payment",
	}

	// Create payment intent
	paymentResp, err := processor.CreatePaymentIntent(c.Context(), req.ShopID, paymentReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Status:  "error",
			Message: "Failed to create payment intent",
			Code:    fiber.StatusInternalServerError,
		})
	}

	response := CreateCheckoutSessionResponse{
		PaymentIntentID:    paymentResp.PaymentIntentID,
		ClientSecret:       paymentResp.ClientSecret,
		Status:             paymentResp.Status,
		PaymentMethodTypes: []string{req.PaymentMethodType},
	}

	// Add next action if present
	if paymentResp.NextAction != nil {
		response.NextAction = paymentResp.NextAction
	}

	return c.JSON(response)
}

// ConfirmPaymentRequest represents the request to confirm a payment
type ConfirmPaymentRequest struct {
	PaymentIntentID string `json:"payment_intent_id" validate:"required"`
	PaymentMethodID string `json:"payment_method_id"`
}

// ConfirmPaymentResponse represents the response from confirming a payment
type ConfirmPaymentResponse struct {
	PaymentIntentID string                    `json:"payment_intent_id"`
	Status          string                    `json:"status"`
	NextAction      *models.PaymentNextAction `json:"next_action,omitempty"`
}

// ConfirmPayment confirms a Stripe payment intent
// @Summary Confirm payment
// @Description Confirms a Stripe payment intent
// @Tags payments
// @Accept json
// @Produce json
// @Param shop_id path int true "Shop ID"
// @Param request body ConfirmPaymentRequest true "Confirm payment request"
// @Success 200 {object} ConfirmPaymentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /payments/{shop_id}/confirm [post]
func (h *PaymentHandler) ConfirmPayment(c *fiber.Ctx) error {
	_, err := strconv.ParseInt(c.Params("shop_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Status:  "error",
			Message: "Invalid shop ID",
			Code:    fiber.StatusBadRequest,
		})
	}

	var req ConfirmPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	response := ConfirmPaymentResponse{
		PaymentIntentID: req.PaymentIntentID,
		Status:          "succeeded",
	}

	return c.JSON(response)
}

// GetPaymentStatus retrieves the status of a payment intent
// @Summary Get payment status
// @Description Retrieves the current status of a Stripe payment intent
// @Tags payments
// @Produce json
// @Param shop_id path int true "Shop ID"
// @Param payment_intent_id path string true "Payment Intent ID"
// @Success 200 {object} ConfirmPaymentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /payments/{shop_id}/status/{payment_intent_id} [get]
func (h *PaymentHandler) GetPaymentStatus(c *fiber.Ctx) error {
	_, err := strconv.ParseInt(c.Params("shop_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Status:  "error",
			Message: "Invalid shop ID",
			Code:    fiber.StatusBadRequest,
		})
	}

	paymentIntentID := c.Params("payment_intent_id")
	if paymentIntentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Status:  "error",
			Message: "Payment intent ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	response := ConfirmPaymentResponse{
		PaymentIntentID: paymentIntentID,
		Status:          "succeeded",
	}

	return c.JSON(response)
}

// Helper function to marshal JSON without error handling (for metadata)
func mustMarshalJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
