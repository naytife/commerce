package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"

	ic "github.com/petrejonn/naytife/internal/httpclient"

	"github.com/petrejonn/naytife/internal/observability"
)

type PayPalService struct {
	repository db.Repository
}

type PayPalConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Environment  string `json:"environment"` // sandbox or live
	TestMode     bool   `json:"test_mode"`
}

type PayPalAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type PayPalOrderRequest struct {
	Intent        string                 `json:"intent"`
	PurchaseUnits []PayPalPurchaseUnit   `json:"purchase_units"`
	PaymentSource *PayPalPaymentSource   `json:"payment_source,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

type PayPalPurchaseUnit struct {
	Amount      PayPalAmount `json:"amount"`
	Description string       `json:"description,omitempty"`
	CustomID    string       `json:"custom_id,omitempty"`
}

type PayPalAmount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type PayPalPaymentSource struct {
	PayPal *PayPalPaymentSourcePayPal `json:"paypal,omitempty"`
}

type PayPalPaymentSourcePayPal struct {
	ExperienceContext *PayPalExperienceContext `json:"experience_context,omitempty"`
}

type PayPalExperienceContext struct {
	ReturnURL string `json:"return_url,omitempty"`
	CancelURL string `json:"cancel_url,omitempty"`
}

type PayPalOrderResponse struct {
	ID     string        `json:"id"`
	Status string        `json:"status"`
	Links  []PayPalLink  `json:"links"`
	Amount *PayPalAmount `json:"amount,omitempty"`
}

type PayPalLink struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

func NewPayPalService(repo db.Repository) *PayPalService {
	return &PayPalService{
		repository: repo,
	}
}

// GetPayPalConfig retrieves PayPal configuration for a shop
func (p *PayPalService) GetPayPalConfig(ctx context.Context, shopID int64) (*PayPalConfig, error) {
	paymentMethod, err := p.repository.GetShopPaymentMethod(ctx, db.GetShopPaymentMethodParams{
		ShopID:     shopID,
		MethodType: db.PaymentMethodTypePaypal,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get PayPal config: %w", err)
	}

	if !paymentMethod.IsEnabled {
		return nil, fmt.Errorf("PayPal is not enabled for this shop")
	}

	var config PayPalConfig
	if err := json.Unmarshal(paymentMethod.Attributes, &config); err != nil {
		return nil, fmt.Errorf("failed to parse PayPal config: %w", err)
	}

	return &config, nil
}

// GetPayPalBaseURL returns the appropriate PayPal API base URL
func (p *PayPalService) GetPayPalBaseURL(config *PayPalConfig) string {
	if config.Environment == "live" && !config.TestMode {
		return "https://api-m.paypal.com"
	}
	return "https://api-m.sandbox.paypal.com"
}

// GetAccessToken retrieves a PayPal access token
func (p *PayPalService) GetAccessToken(ctx context.Context, config *PayPalConfig) (*PayPalAccessTokenResponse, error) {
	baseURL := p.GetPayPalBaseURL(config)

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/v1/oauth2/token",
		bytes.NewBuffer([]byte("grant_type=client_credentials")))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.ClientID, config.ClientSecret)

	observability.InjectTraceHeaders(ctx, req)
	observability.EnsureRequestID(req)
	_, finish := observability.StartSpan(ctx, "GetAccessToken", "paypal", "POST", req.URL.String())
	defer func() { finish(0, nil) }()
	start := time.Now()
	resp, err := ic.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("paypal api error: %s", string(body))
	}

	var tokenResp PayPalAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	observability.RecordServiceRequest("paypal", req.Method, req.URL.String(), resp.StatusCode, time.Since(start))

	return &tokenResp, nil
}

// ValidateConfig validates the PayPal configuration
func (p *PayPalService) ValidateConfig(config map[string]interface{}) error {
	clientID, ok := config["client_id"].(string)
	if !ok || clientID == "" {
		return fmt.Errorf("client_id is required for PayPal configuration")
	}

	clientSecret, ok := config["client_secret"].(string)
	if !ok || clientSecret == "" {
		return fmt.Errorf("client_secret is required for PayPal configuration")
	}

	environment, ok := config["environment"].(string)
	if !ok || (environment != "sandbox" && environment != "live") {
		return fmt.Errorf("environment must be either 'sandbox' or 'live'")
	}

	return nil
}

// ProcessPayment processes a payment request and returns the payment response
func (p *PayPalService) ProcessPayment(ctx context.Context, shopID int64, req models.PaymentRequest, amount float64, currencyCode string) (*models.PaymentResponse, error) {
	config, err := p.GetPayPalConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	// Get access token
	tokenResp, err := p.GetAccessToken(ctx, config)
	if err != nil {
		return nil, err
	}

	// Create PayPal order
	orderReq := PayPalOrderRequest{
		Intent: "CAPTURE",
		PurchaseUnits: []PayPalPurchaseUnit{
			{
				Amount: PayPalAmount{
					CurrencyCode: currencyCode,
					Value:        fmt.Sprintf("%.2f", amount),
				},
				Description: "Payment for order",
				CustomID:    req.CheckoutSessionID,
			},
		},
		Metadata: map[string]interface{}{
			"shop_id":             strconv.FormatInt(shopID, 10),
			"checkout_session_id": req.CheckoutSessionID,
		},
	}

	// Add payment source for immediate capture if payment method is provided
	if returnURL, ok := req.PaymentDetails["return_url"].(string); ok {
		cancelURL, _ := req.PaymentDetails["cancel_url"].(string)
		orderReq.PaymentSource = &PayPalPaymentSource{
			PayPal: &PayPalPaymentSourcePayPal{
				ExperienceContext: &PayPalExperienceContext{
					ReturnURL: returnURL,
					CancelURL: cancelURL,
				},
			},
		}
	}

	orderResp, err := p.createPayPalOrder(ctx, config, tokenResp.AccessToken, &orderReq)
	if err != nil {
		return nil, err
	}

	status := p.ConvertPayPalStatusToPaymentStatus(orderResp.Status)

	response := &models.PaymentResponse{
		PaymentID:     orderResp.ID,
		Status:        status,
		Amount:        amount,
		CurrencyCode:  currencyCode,
		PaymentMethod: "paypal",
		TransactionID: orderResp.ID,
	}

	// Add next action for approval if needed
	if orderResp.Status == "CREATED" || orderResp.Status == "PAYER_ACTION_REQUIRED" {
		for _, link := range orderResp.Links {
			if link.Rel == "approve" {
				response.NextAction = &models.PaymentNextAction{
					Type: "redirect_to_url",
					Data: map[string]interface{}{
						"redirect_url": link.Href,
					},
				}
				break
			}
		}
	}

	if orderResp.Status == "COMPLETED" {
		now := time.Now()
		response.ProcessedAt = now.Format(time.RFC3339)
	}

	return response, nil
}

// CreatePaymentIntent creates a payment intent for deferred payment processing
func (p *PayPalService) CreatePaymentIntent(ctx context.Context, shopID int64, req models.PaymentIntentRequest) (*models.PaymentIntentResponse, error) {
	config, err := p.GetPayPalConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	// Get access token
	tokenResp, err := p.GetAccessToken(ctx, config)
	if err != nil {
		return nil, err
	}

	// Create PayPal order (this serves as a payment intent)
	orderReq := PayPalOrderRequest{
		Intent: "CAPTURE",
		PurchaseUnits: []PayPalPurchaseUnit{
			{
				Amount: PayPalAmount{
					CurrencyCode: req.CurrencyCode,
					Value:        fmt.Sprintf("%.2f", req.Amount),
				},
				Description: req.Description,
				CustomID:    req.CheckoutSessionID,
			},
		},
		Metadata: map[string]interface{}{
			"shop_id":             strconv.FormatInt(shopID, 10),
			"checkout_session_id": req.CheckoutSessionID,
		},
	}

	orderResp, err := p.createPayPalOrder(ctx, config, tokenResp.AccessToken, &orderReq)
	if err != nil {
		return nil, err
	}

	response := &models.PaymentIntentResponse{
		PaymentIntentID: orderResp.ID,
		Status:          p.ConvertPayPalStatusToPaymentStatus(orderResp.Status),
		Amount:          req.Amount,
		CurrencyCode:    req.CurrencyCode,
	}

	// Add next action for approval
	for _, link := range orderResp.Links {
		if link.Rel == "approve" {
			response.NextAction = &models.PaymentNextAction{
				Type: "redirect_to_url",
				Data: map[string]interface{}{
					"redirect_url": link.Href,
				},
			}
			break
		}
	}

	return response, nil
}

// ConfirmPayment confirms a payment intent or transaction
func (p *PayPalService) ConfirmPayment(ctx context.Context, shopID int64, paymentID string) (*models.PaymentResponse, error) {
	config, err := p.GetPayPalConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	// Get access token
	tokenResp, err := p.GetAccessToken(ctx, config)
	if err != nil {
		return nil, err
	}

	// Capture the PayPal order
	captureResp, err := p.capturePayPalOrder(ctx, config, tokenResp.AccessToken, paymentID)
	if err != nil {
		return nil, err
	}

	var amount float64
	if captureResp.Amount != nil {
		amount, _ = strconv.ParseFloat(captureResp.Amount.Value, 64)
	}

	response := &models.PaymentResponse{
		PaymentID:     captureResp.ID,
		Status:        p.ConvertPayPalStatusToPaymentStatus(captureResp.Status),
		Amount:        amount,
		CurrencyCode:  captureResp.Amount.CurrencyCode,
		PaymentMethod: "paypal",
		TransactionID: captureResp.ID,
	}

	if captureResp.Status == "COMPLETED" {
		now := time.Now()
		response.ProcessedAt = now.Format(time.RFC3339)
	}

	return response, nil
}

// GetPaymentStatus retrieves the current status of a payment
func (p *PayPalService) GetPaymentStatus(ctx context.Context, shopID int64, paymentID string) (*models.PaymentResponse, error) {
	config, err := p.GetPayPalConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	// Get access token
	tokenResp, err := p.GetAccessToken(ctx, config)
	if err != nil {
		return nil, err
	}

	// Get PayPal order details
	orderResp, err := p.getPayPalOrder(ctx, config, tokenResp.AccessToken, paymentID)
	if err != nil {
		return nil, err
	}

	var amount float64
	if orderResp.Amount != nil {
		amount, _ = strconv.ParseFloat(orderResp.Amount.Value, 64)
	}

	response := &models.PaymentResponse{
		PaymentID:     orderResp.ID,
		Status:        p.ConvertPayPalStatusToPaymentStatus(orderResp.Status),
		Amount:        amount,
		CurrencyCode:  orderResp.Amount.CurrencyCode,
		PaymentMethod: "paypal",
		TransactionID: orderResp.ID,
	}

	if orderResp.Status == "COMPLETED" {
		now := time.Now()
		response.ProcessedAt = now.Format(time.RFC3339)
	}

	return response, nil
}

// RefundPayment processes a refund for a completed payment
func (p *PayPalService) RefundPayment(ctx context.Context, shopID int64, paymentID string, amount float64, reason string) (*models.PaymentResponse, error) {
	// PayPal refunds require the capture ID, not the order ID
	// This is a simplified implementation
	return &models.PaymentResponse{
		PaymentID:     paymentID,
		Status:        "refund_pending",
		Amount:        amount,
		PaymentMethod: "paypal",
		PaymentDetails: map[string]interface{}{
			"message": "PayPal refunds require additional implementation with capture IDs",
			"reason":  reason,
		},
	}, nil
}

// HandleWebhook processes webhook events from PayPal
func (p *PayPalService) HandleWebhook(ctx context.Context, payload []byte, signature string) (*models.PaymentWebhookPayload, error) {
	var webhookEvent map[string]interface{}
	if err := json.Unmarshal(payload, &webhookEvent); err != nil {
		return nil, fmt.Errorf("failed to parse PayPal webhook payload: %w", err)
	}

	webhookPayload := &models.PaymentWebhookPayload{
		Provider:   "paypal",
		RawPayload: webhookEvent,
		ReceivedAt: time.Now(),
	}

	if eventType, ok := webhookEvent["event_type"].(string); ok {
		webhookPayload.EventType = eventType
	}

	// Extract resource information
	if resource, ok := webhookEvent["resource"].(map[string]interface{}); ok {
		if id, exists := resource["id"].(string); exists {
			webhookPayload.PaymentID = id
			webhookPayload.TransactionID = id
		}
		if status, exists := resource["status"].(string); exists {
			webhookPayload.Status = p.ConvertPayPalStatusToPaymentStatus(status)
		}
		if amount, exists := resource["amount"].(map[string]interface{}); exists {
			if value, ok := amount["value"].(string); ok {
				if amountFloat, err := strconv.ParseFloat(value, 64); err == nil {
					webhookPayload.Amount = amountFloat
				}
			}
			if currency, ok := amount["currency_code"].(string); ok {
				webhookPayload.Currency = currency
			}
		}
	}

	return webhookPayload, nil
}

// Helper methods

func (p *PayPalService) createPayPalOrder(ctx context.Context, config *PayPalConfig, accessToken string, orderReq *PayPalOrderRequest) (*PayPalOrderResponse, error) {
	baseURL := p.GetPayPalBaseURL(config)

	reqBody, err := json.Marshal(orderReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/v2/checkout/orders", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	observability.InjectTraceHeaders(ctx, req)
	var finish func(int, error)
	_, finish = observability.StartSpan(ctx, "createPayPalOrder", "paypal", "POST", req.URL.String())
	defer func() { finish(0, nil) }()
	resp, err := ic.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create PayPal order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("paypal api error: %s", string(body))
	}

	var orderResp PayPalOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return nil, fmt.Errorf("failed to decode order response: %w", err)
	}

	return &orderResp, nil
}

func (p *PayPalService) capturePayPalOrder(ctx context.Context, config *PayPalConfig, accessToken, orderID string) (*PayPalOrderResponse, error) {
	baseURL := p.GetPayPalBaseURL(config)

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/v2/checkout/orders/"+orderID+"/capture", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create capture request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	observability.InjectTraceHeaders(ctx, req)
	var finish func(int, error)
	_, finish = observability.StartSpan(ctx, "capturePayPalOrder", "paypal", "POST", req.URL.String())
	defer func() { finish(0, nil) }()
	resp, err := ic.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to capture PayPal order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("paypal api error: %s", string(body))
	}

	var captureResp PayPalOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&captureResp); err != nil {
		return nil, fmt.Errorf("failed to decode capture response: %w", err)
	}

	return &captureResp, nil
}

func (p *PayPalService) getPayPalOrder(ctx context.Context, config *PayPalConfig, accessToken, orderID string) (*PayPalOrderResponse, error) {
	baseURL := p.GetPayPalBaseURL(config)

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/v2/checkout/orders/"+orderID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get order request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	observability.InjectTraceHeaders(ctx, req)
	var finish func(int, error)
	_, finish = observability.StartSpan(ctx, "getPayPalOrder", "paypal", "GET", req.URL.String())
	defer func() { finish(0, nil) }()
	resp, err := ic.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get PayPal order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("paypal api error: %s", string(body))
	}

	var orderResp PayPalOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return nil, fmt.Errorf("failed to decode order response: %w", err)
	}

	return &orderResp, nil
}

// ConvertPayPalStatusToPaymentStatus converts PayPal order status to our payment status
func (p *PayPalService) ConvertPayPalStatusToPaymentStatus(status string) string {
	switch status {
	case "COMPLETED":
		return "completed"
	case "APPROVED":
		return "processing"
	case "CREATED":
		return "pending"
	case "SAVED":
		return "pending"
	case "VOIDED":
		return "cancelled"
	case "PAYER_ACTION_REQUIRED":
		return "requires_action"
	default:
		return "pending"
	}
}
