package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/stripe/stripe-go/v81/refund"
	"github.com/stripe/stripe-go/v81/webhook"
)

type StripeService struct {
	repository db.Repository
}

type StripeConfig struct {
	PublishableKey string `json:"publishable_key"`
	SecretKey      string `json:"secret_key"`
	WebhookSecret  string `json:"webhook_secret"`
	TestMode       bool   `json:"test_mode"`
}

func NewStripeService(repo db.Repository) *StripeService {
	return &StripeService{
		repository: repo,
	}
}

// GetStripeConfig retrieves Stripe configuration for a shop
func (s *StripeService) GetStripeConfig(ctx context.Context, shopID int64) (*StripeConfig, error) {
	paymentMethod, err := s.repository.GetShopPaymentMethod(ctx, db.GetShopPaymentMethodParams{
		ShopID:     shopID,
		MethodType: db.PaymentMethodTypeStripe,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get Stripe config: %w", err)
	}

	if !paymentMethod.IsEnabled {
		return nil, fmt.Errorf("Stripe is not enabled for this shop")
	}

	var config StripeConfig
	if err := json.Unmarshal(paymentMethod.Attributes, &config); err != nil {
		return nil, fmt.Errorf("failed to parse Stripe config: %w", err)
	}

	return &config, nil
}

// ConvertStripeStatusToPaymentStatus converts Stripe payment intent status to our payment status
func (s *StripeService) ConvertStripeStatusToPaymentStatus(status stripe.PaymentIntentStatus) string {
	switch status {
	case stripe.PaymentIntentStatusSucceeded:
		return "completed"
	case stripe.PaymentIntentStatusProcessing:
		return "processing"
	case stripe.PaymentIntentStatusRequiresPaymentMethod:
		return "pending"
	case stripe.PaymentIntentStatusRequiresConfirmation:
		return "pending"
	case stripe.PaymentIntentStatusRequiresAction:
		return "requires_action"
	case stripe.PaymentIntentStatusCanceled:
		return "cancelled"
	default:
		return "pending"
	}
}

// GetNextActionData extracts next action data from Stripe PaymentIntent
func (s *StripeService) GetNextActionData(intent *stripe.PaymentIntent) *models.PaymentNextAction {
	if intent.NextAction == nil {
		return nil
	}

	nextAction := &models.PaymentNextAction{
		Type: string(intent.NextAction.Type),
		Data: make(map[string]interface{}),
	}

	switch intent.NextAction.Type {
	case "redirect_to_url":
		if intent.NextAction.RedirectToURL != nil {
			nextAction.Data["redirect_url"] = intent.NextAction.RedirectToURL.URL
			nextAction.Data["return_url"] = intent.NextAction.RedirectToURL.ReturnURL
		}
	case "use_stripe_sdk":
		nextAction.Data["client_secret"] = intent.ClientSecret
	}

	return nextAction
}

// ValidateConfig validates the Stripe configuration
func (s *StripeService) ValidateConfig(config map[string]interface{}) error {
	secretKey, ok := config["secret_key"].(string)
	if !ok || secretKey == "" {
		return fmt.Errorf("secret_key is required for Stripe configuration")
	}

	publishableKey, ok := config["publishable_key"].(string)
	if !ok || publishableKey == "" {
		return fmt.Errorf("publishable_key is required for Stripe configuration")
	}

	return nil
}

// ProcessPayment processes a payment request and returns the payment response
func (s *StripeService) ProcessPayment(ctx context.Context, shopID int64, req models.PaymentRequest, amount float64, currencyCode string) (*models.PaymentResponse, error) {
	config, err := s.GetStripeConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	// Set the Stripe API key
	stripe.Key = config.SecretKey

	// Convert amount to cents for Stripe
	amountCents := int64(amount * 100)

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(currencyCode),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	// Add metadata for tracking
	params.Metadata = map[string]string{
		"shop_id":             strconv.FormatInt(shopID, 10),
		"checkout_session_id": req.CheckoutSessionID,
	}

	// Handle payment method details from request
	if paymentMethodID, ok := req.PaymentDetails["payment_method_id"].(string); ok && paymentMethodID != "" {
		params.PaymentMethod = stripe.String(paymentMethodID)
		params.ConfirmationMethod = stripe.String("manual")
		params.Confirm = stripe.Bool(true)
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	response := &models.PaymentResponse{
		PaymentID:     intent.ID,
		Status:        s.ConvertStripeStatusToPaymentStatus(intent.Status),
		Amount:        amount,
		CurrencyCode:  currencyCode,
		PaymentMethod: "stripe",
		TransactionID: intent.ID,
	}

	if intent.Status == stripe.PaymentIntentStatusSucceeded {
		now := time.Now()
		response.ProcessedAt = now.Format(time.RFC3339)
	}

	// Add next action if required
	if nextAction := s.GetNextActionData(intent); nextAction != nil {
		response.NextAction = nextAction
	}

	return response, nil
}

// CreatePaymentIntent creates a payment intent for deferred payment processing
func (s *StripeService) CreatePaymentIntent(ctx context.Context, shopID int64, req models.PaymentIntentRequest) (*models.PaymentIntentResponse, error) {
	config, err := s.GetStripeConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	// Set the Stripe API key
	stripe.Key = config.SecretKey

	// Convert amount to cents for Stripe
	amountCents := int64(req.Amount * 100)

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(req.CurrencyCode),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	// Add metadata for tracking
	params.Metadata = map[string]string{
		"shop_id":             strconv.FormatInt(shopID, 10),
		"checkout_session_id": req.CheckoutSessionID,
	}

	if req.Description != "" {
		params.Description = stripe.String(req.Description)
	}

	if req.CustomerID != "" {
		params.Customer = stripe.String(req.CustomerID)
	}

	if req.PaymentMethodID != "" {
		params.PaymentMethod = stripe.String(req.PaymentMethodID)
	}

	// Add custom metadata
	if req.Metadata != nil {
		for key, value := range req.Metadata {
			if strValue, ok := value.(string); ok {
				params.Metadata[key] = strValue
			}
		}
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	response := &models.PaymentIntentResponse{
		PaymentIntentID: intent.ID,
		ClientSecret:    intent.ClientSecret,
		Status:          s.ConvertStripeStatusToPaymentStatus(intent.Status),
		Amount:          req.Amount,
		CurrencyCode:    req.CurrencyCode,
	}

	// Add next action if required
	if nextAction := s.GetNextActionData(intent); nextAction != nil {
		response.NextAction = nextAction
	}

	return response, nil
}

// ConfirmPayment confirms a payment intent or transaction
func (s *StripeService) ConfirmPayment(ctx context.Context, shopID int64, paymentID string) (*models.PaymentResponse, error) {
	config, err := s.GetStripeConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	// Set the Stripe API key
	stripe.Key = config.SecretKey

	params := &stripe.PaymentIntentConfirmParams{}
	intent, err := paymentintent.Confirm(paymentID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment intent: %w", err)
	}

	amount := float64(intent.Amount) / 100 // Convert from cents

	response := &models.PaymentResponse{
		PaymentID:     intent.ID,
		Status:        s.ConvertStripeStatusToPaymentStatus(intent.Status),
		Amount:        amount,
		CurrencyCode:  string(intent.Currency),
		PaymentMethod: "stripe",
		TransactionID: intent.ID,
	}

	if intent.Status == stripe.PaymentIntentStatusSucceeded {
		now := time.Now()
		response.ProcessedAt = now.Format(time.RFC3339)
	}

	// Add next action if required
	if nextAction := s.GetNextActionData(intent); nextAction != nil {
		response.NextAction = nextAction
	}

	return response, nil
}

// GetPaymentStatus retrieves the current status of a payment
func (s *StripeService) GetPaymentStatus(ctx context.Context, shopID int64, paymentID string) (*models.PaymentResponse, error) {
	config, err := s.GetStripeConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	// Set the Stripe API key
	stripe.Key = config.SecretKey

	intent, err := paymentintent.Get(paymentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payment intent: %w", err)
	}

	amount := float64(intent.Amount) / 100 // Convert from cents

	response := &models.PaymentResponse{
		PaymentID:     intent.ID,
		Status:        s.ConvertStripeStatusToPaymentStatus(intent.Status),
		Amount:        amount,
		CurrencyCode:  string(intent.Currency),
		PaymentMethod: "stripe",
		TransactionID: intent.ID,
	}

	if intent.Status == stripe.PaymentIntentStatusSucceeded && intent.Created > 0 {
		processedAt := time.Unix(intent.Created, 0)
		response.ProcessedAt = processedAt.Format(time.RFC3339)
	}

	// Add next action if required
	if nextAction := s.GetNextActionData(intent); nextAction != nil {
		response.NextAction = nextAction
	}

	return response, nil
}

// RefundPayment processes a refund for a completed payment
func (s *StripeService) RefundPayment(ctx context.Context, shopID int64, paymentID string, amount float64, reason string) (*models.PaymentResponse, error) {
	config, err := s.GetStripeConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	// Set the Stripe API key
	stripe.Key = config.SecretKey

	// Convert amount to cents for Stripe
	amountCents := int64(amount * 100)

	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(paymentID),
		Amount:        stripe.Int64(amountCents),
	}

	if reason != "" {
		params.Reason = stripe.String(reason)
	}

	refundObj, err := refund.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	refundAmount := float64(refundObj.Amount) / 100 // Convert from cents
	now := time.Now()

	response := &models.PaymentResponse{
		PaymentID:     paymentID,
		Status:        "refunded",
		Amount:        refundAmount,
		CurrencyCode:  string(refundObj.Currency),
		PaymentMethod: "stripe",
		TransactionID: refundObj.ID,
		ProcessedAt:   now.Format(time.RFC3339),
		PaymentDetails: map[string]interface{}{
			"refund_id":     refundObj.ID,
			"refund_status": string(refundObj.Status),
			"reason":        reason,
		},
	}

	return response, nil
}

// HandleWebhook processes webhook events from Stripe
func (s *StripeService) HandleWebhook(ctx context.Context, payload []byte, signature string) (*models.PaymentWebhookPayload, error) {
	// Note: To properly validate webhooks, we need the webhook secret from shop configuration
	// For now, we'll parse the event without validation and add validation later
	event, err := webhook.ConstructEvent(payload, signature, "")
	if err != nil {
		// If validation fails, try to parse the event directly for development
		var eventData stripe.Event
		if parseErr := json.Unmarshal(payload, &eventData); parseErr != nil {
			return nil, fmt.Errorf("failed to parse webhook payload: %w", parseErr)
		}
		event = eventData
	}

	webhookPayload := &models.PaymentWebhookPayload{
		Provider:   "stripe",
		EventType:  string(event.Type),
		RawPayload: event,
		ReceivedAt: time.Now(),
	}

	// Extract payment information based on event type
	switch event.Type {
	case "payment_intent.succeeded", "payment_intent.payment_failed", "payment_intent.requires_action":
		if event.Data != nil && event.Data.Object != nil {
			// Convert the object to a map for easier access
			objectBytes, _ := json.Marshal(event.Data.Object)
			var paymentIntent map[string]interface{}
			json.Unmarshal(objectBytes, &paymentIntent)

			if id, exists := paymentIntent["id"].(string); exists {
				webhookPayload.PaymentID = id
				webhookPayload.TransactionID = id
			}
			if amount, exists := paymentIntent["amount"].(float64); exists {
				webhookPayload.Amount = amount / 100 // Convert from cents
			}
			if currency, exists := paymentIntent["currency"].(string); exists {
				webhookPayload.Currency = currency
			}
			if status, exists := paymentIntent["status"].(string); exists {
				webhookPayload.Status = s.ConvertStripeStatusToPaymentStatus(stripe.PaymentIntentStatus(status))
			}
			if metadata, exists := paymentIntent["metadata"].(map[string]interface{}); exists {
				webhookPayload.Metadata = metadata
			}
		}
	}

	return webhookPayload, nil
}

// Backward compatibility methods for existing code

// CreatePaymentIntentFromOrder creates a Stripe PaymentIntent from order data (for backward compatibility)
func (s *StripeService) CreatePaymentIntentFromOrder(ctx context.Context, shopID int64, req models.PaymentRequest, amount float64, currencyCode string) (*stripe.PaymentIntent, error) {
	config, err := s.GetStripeConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	// Set the Stripe API key
	stripe.Key = config.SecretKey

	// Convert amount to cents for Stripe
	amountCents := int64(amount * 100)

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(currencyCode),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	// Add metadata for tracking
	params.Metadata = map[string]string{
		"shop_id":             strconv.FormatInt(shopID, 10),
		"checkout_session_id": req.CheckoutSessionID,
	}

	// Handle payment method details from request
	if paymentDetails, ok := req.PaymentDetails["payment_method_id"].(string); ok && paymentDetails != "" {
		params.PaymentMethod = stripe.String(paymentDetails)
		params.ConfirmationMethod = stripe.String("manual")
		params.Confirm = stripe.Bool(true)
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return intent, nil
}

// ConfirmPaymentIntentByID confirms a Stripe PaymentIntent (for backward compatibility)
func (s *StripeService) ConfirmPaymentIntentByID(ctx context.Context, shopID int64, paymentIntentID string) (*stripe.PaymentIntent, error) {
	config, err := s.GetStripeConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	// Set the Stripe API key
	stripe.Key = config.SecretKey

	params := &stripe.PaymentIntentConfirmParams{}
	intent, err := paymentintent.Confirm(paymentIntentID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment intent: %w", err)
	}

	return intent, nil
}

// RetrievePaymentIntent retrieves a Stripe PaymentIntent (for backward compatibility)
func (s *StripeService) RetrievePaymentIntent(ctx context.Context, shopID int64, paymentIntentID string) (*stripe.PaymentIntent, error) {
	config, err := s.GetStripeConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	// Set the Stripe API key
	stripe.Key = config.SecretKey

	intent, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payment intent: %w", err)
	}

	return intent, nil
}
