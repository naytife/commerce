package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/petrejonn/naytife/internal/observability"
	"go.uber.org/zap"
)

type PaystackService struct {
	repository  db.Repository
	RetryClient *retryablehttp.Client
}

type PaystackConfig struct {
	PublicKey string `json:"public_key"`
	SecretKey string `json:"secret_key"`
	TestMode  bool   `json:"test_mode"`
}

// Paystack API structures
type PaystackInitializeRequest struct {
	Email    string                 `json:"email"`
	Amount   string                 `json:"amount"` // Amount in kobo (Nigerian currency subunit)
	Currency string                 `json:"currency,omitempty"`
	Callback string                 `json:"callback_url,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type PaystackInitializeResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

type PaystackVerifyResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID              int64  `json:"id"`
		Domain          string `json:"domain"`
		Status          string `json:"status"`
		Reference       string `json:"reference"`
		Amount          int64  `json:"amount"`
		GatewayResponse string `json:"gateway_response"`
		PaidAt          string `json:"paid_at"`
		CreatedAt       string `json:"created_at"`
		Channel         string `json:"channel"`
		Currency        string `json:"currency"`
		IPAddress       string `json:"ip_address"`
		Customer        struct {
			Email string `json:"email"`
		} `json:"customer"`
	} `json:"data"`
}

type PaystackRefundRequest struct {
	Transaction string `json:"transaction"`
	Amount      int64  `json:"amount,omitempty"` // Amount in kobo
}

type PaystackRefundResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Transaction struct {
			ID        int64  `json:"id"`
			Reference string `json:"reference"`
			Amount    int64  `json:"amount"`
			Status    string `json:"status"`
		} `json:"transaction"`
		Integration struct {
			ID int64 `json:"id"`
		} `json:"integration"`
		Deducted struct {
			Amount int64 `json:"amount"`
		} `json:"deducted"`
		Channel    string `json:"channel"`
		FullRefund bool   `json:"full_refund"`
	} `json:"data"`
}

type PaystackErrorResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func NewPaystackService(repo db.Repository) *PaystackService {
	return &PaystackService{
		repository: repo,
	}
}

// GetPaystackConfig retrieves Paystack configuration for a shop
func (p *PaystackService) GetPaystackConfig(ctx context.Context, shopID int64) (*PaystackConfig, error) {
	paymentMethod, err := p.repository.GetShopPaymentMethod(ctx, db.GetShopPaymentMethodParams{
		ShopID:     shopID,
		MethodType: db.PaymentMethodTypePaystack,
	})
	if err != nil {
		zap.L().Error("GetPaystackConfig: failed to get paystack config", zap.Int64("shop_id", shopID), zap.Error(err))
		return nil, fmt.Errorf("failed to get paystack config: %w", err)
	}

	if !paymentMethod.IsEnabled {
		zap.L().Warn("GetPaystackConfig: paystack not enabled for shop", zap.Int64("shop_id", shopID))
		return nil, fmt.Errorf("paystack is not enabled for this shop")
	}

	var config PaystackConfig
	if err := json.Unmarshal(paymentMethod.Attributes, &config); err != nil {
		zap.L().Error("GetPaystackConfig: failed to parse Paystack config", zap.Int64("shop_id", shopID), zap.Error(err))
		return nil, fmt.Errorf("failed to parse Paystack config: %w", err)
	}

	return &config, nil
}

// getBaseURL returns the appropriate Paystack API base URL
func (p *PaystackService) getBaseURL(testMode bool) string {
	// Paystack uses the same URL for both test and live, differentiated by API keys
	return "https://api.paystack.co"
}

// makePaystackRequest makes an HTTP request to Paystack API
func (p *PaystackService) makePaystackRequest(ctx context.Context, method, url string, headers map[string]string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")

	observability.InjectTraceHeaders(ctx, req)
	observability.EnsureRequestID(req)
	_, finish := observability.StartSpan(ctx, "makePaystackRequest", "paystack", method, url)
	defer finish(0, nil)
	start := time.Now()
	var resp *http.Response
	if p.RetryClient != nil {
		resp, err = p.RetryClient.StandardClient().Do(req)
	} else {
		resp, err = http.DefaultClient.Do(req)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	observability.RecordServiceRequest("paystack", req.Method, req.URL.String(), resp.StatusCode, time.Since(start))

	return resp, nil
}

// ValidateConfig validates the Paystack configuration
func (p *PaystackService) ValidateConfig(config map[string]interface{}) error {
	requiredFields := []string{"public_key", "secret_key"}

	for _, field := range requiredFields {
		if value, exists := config[field]; !exists || value == "" {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validate key formats
	if publicKey, ok := config["public_key"].(string); ok {
		if !strings.HasPrefix(publicKey, "pk_") {
			return fmt.Errorf("invalid public key format")
		}
	}

	if secretKey, ok := config["secret_key"].(string); ok {
		if !strings.HasPrefix(secretKey, "sk_") {
			return fmt.Errorf("invalid secret key format")
		}
	}

	return nil
}

// ConvertPaystackStatusToPaymentStatus converts Paystack transaction status to our payment status
func (p *PaystackService) ConvertPaystackStatusToPaymentStatus(status string) string {
	switch strings.ToLower(status) {
	case "success":
		return "completed"
	case "pending":
		return "pending"
	case "failed":
		return "failed"
	case "abandoned":
		return "cancelled"
	case "reversed":
		return "refunded"
	default:
		return "pending"
	}
}

// ProcessPayment processes a payment request and returns the payment response
func (p *PaystackService) ProcessPayment(ctx context.Context, shopID int64, req models.PaymentRequest, amount float64, currencyCode string) (*models.PaymentResponse, error) {
	config, err := p.GetPaystackConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	baseURL := p.getBaseURL(config.TestMode)

	// Convert amount to kobo (smallest currency unit for NGN)
	// For other currencies, we might need different conversion logic
	amountKobo := int64(amount * 100)

	// Create transaction initialization request
	initReq := PaystackInitializeRequest{
		Email:    "customer@example.com", // TODO: Get customer email from checkout session
		Amount:   strconv.FormatInt(amountKobo, 10),
		Currency: currencyCode,
		Metadata: map[string]interface{}{
			"shop_id":             strconv.FormatInt(shopID, 10),
			"checkout_session_id": req.CheckoutSessionID,
		},
	}

	// Add custom metadata from request payment details
	if req.PaymentDetails != nil {
		for key, value := range req.PaymentDetails {
			initReq.Metadata[key] = value
		}
	}

	initReqJSON, err := json.Marshal(initReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal initialization request: %w", err)
	}

	// Make request to initialize transaction
	headers := map[string]string{
		"Authorization": "Bearer " + config.SecretKey,
	}

	resp, err := p.makePaystackRequest(ctx, "POST", baseURL+"/transaction/initialize", headers, bytes.NewBuffer(initReqJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize transaction: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp PaystackErrorResponse
		if jsonErr := json.Unmarshal(respBody, &errorResp); jsonErr == nil {
			return nil, fmt.Errorf("paystack api error: %s", errorResp.Message)
		}
		return nil, fmt.Errorf("paystack api error: status %d", resp.StatusCode)
	}

	var initResp PaystackInitializeResponse
	if err := json.Unmarshal(respBody, &initResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !initResp.Status {
		return nil, fmt.Errorf("failed to initialize transaction: %s", initResp.Message)
	}

	response := &models.PaymentResponse{
		PaymentID:     initResp.Data.Reference,
		Status:        "pending",
		Amount:        amount,
		CurrencyCode:  currencyCode,
		PaymentMethod: "paystack",
		TransactionID: initResp.Data.Reference,
		NextAction: &models.PaymentNextAction{
			Type: "redirect",
			Data: map[string]interface{}{
				"redirect_url": initResp.Data.AuthorizationURL,
			},
		},
	}

	return response, nil
}

// CreatePaymentIntent creates a payment intent for deferred payment processing
func (p *PaystackService) CreatePaymentIntent(ctx context.Context, shopID int64, req models.PaymentIntentRequest) (*models.PaymentIntentResponse, error) {
	config, err := p.GetPaystackConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	baseURL := p.getBaseURL(config.TestMode)

	// Convert amount to kobo
	amountKobo := int64(req.Amount * 100)

	// Create transaction initialization request
	initReq := PaystackInitializeRequest{
		Email:    "customer@example.com", // TODO: Get customer email from checkout session
		Amount:   strconv.FormatInt(amountKobo, 10),
		Currency: req.CurrencyCode,
		Metadata: map[string]interface{}{
			"shop_id":             strconv.FormatInt(shopID, 10),
			"checkout_session_id": req.CheckoutSessionID,
		},
	}

	// Add custom metadata from request
	if req.Metadata != nil {
		for key, value := range req.Metadata {
			initReq.Metadata[key] = value
		}
	}

	initReqJSON, err := json.Marshal(initReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal initialization request: %w", err)
	}

	// Make request to initialize transaction
	headers := map[string]string{
		"Authorization": "Bearer " + config.SecretKey,
	}

	resp, err := p.makePaystackRequest(ctx, "POST", baseURL+"/transaction/initialize", headers, bytes.NewBuffer(initReqJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize transaction: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp PaystackErrorResponse
		if jsonErr := json.Unmarshal(respBody, &errorResp); jsonErr == nil {
			return nil, fmt.Errorf("paystack api error: %s", errorResp.Message)
		}
		return nil, fmt.Errorf("paystack api error: status %d", resp.StatusCode)
	}

	var initResp PaystackInitializeResponse
	if err := json.Unmarshal(respBody, &initResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !initResp.Status {
		return nil, fmt.Errorf("failed to initialize transaction: %s", initResp.Message)
	}

	response := &models.PaymentIntentResponse{
		PaymentIntentID: initResp.Data.Reference,
		ClientSecret:    initResp.Data.AccessCode,
		Status:          "pending",
		Amount:          req.Amount,
		CurrencyCode:    req.CurrencyCode,
		NextAction: &models.PaymentNextAction{
			Type: "redirect",
			Data: map[string]interface{}{
				"redirect_url": initResp.Data.AuthorizationURL,
			},
		},
	}

	return response, nil
}

// ConfirmPayment confirms a payment intent or transaction
func (p *PaystackService) ConfirmPayment(ctx context.Context, shopID int64, paymentID string) (*models.PaymentResponse, error) {
	// In Paystack, payment confirmation is done by verifying the transaction
	return p.GetPaymentStatus(ctx, shopID, paymentID)
}

// GetPaymentStatus retrieves the current status of a payment
func (p *PaystackService) GetPaymentStatus(ctx context.Context, shopID int64, paymentID string) (*models.PaymentResponse, error) {
	config, err := p.GetPaystackConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	baseURL := p.getBaseURL(config.TestMode)

	// Make request to verify transaction
	headers := map[string]string{
		"Authorization": "Bearer " + config.SecretKey,
	}

	resp, err := p.makePaystackRequest(ctx, "GET", baseURL+"/transaction/verify/"+paymentID, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to verify transaction: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp PaystackErrorResponse
		if jsonErr := json.Unmarshal(respBody, &errorResp); jsonErr == nil {
			return nil, fmt.Errorf("paystack api error: %s", errorResp.Message)
		}
		return nil, fmt.Errorf("paystack api error: status %d", resp.StatusCode)
	}

	var verifyResp PaystackVerifyResponse
	if err := json.Unmarshal(respBody, &verifyResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !verifyResp.Status {
		return nil, fmt.Errorf("failed to verify transaction: %s", verifyResp.Message)
	}

	// Convert amount from kobo to main currency unit
	amount := float64(verifyResp.Data.Amount) / 100

	response := &models.PaymentResponse{
		PaymentID:     verifyResp.Data.Reference,
		Status:        p.ConvertPaystackStatusToPaymentStatus(verifyResp.Data.Status),
		Amount:        amount,
		CurrencyCode:  verifyResp.Data.Currency,
		PaymentMethod: "paystack",
		TransactionID: verifyResp.Data.Reference,
	}

	// Set processed time if payment was successful and we have a paid_at timestamp
	if verifyResp.Data.Status == "success" && verifyResp.Data.PaidAt != "" {
		if processedAt, err := time.Parse(time.RFC3339, verifyResp.Data.PaidAt); err == nil {
			response.ProcessedAt = processedAt.Format(time.RFC3339)
		}
	}

	return response, nil
}

// RefundPayment processes a refund for a completed payment
func (p *PaystackService) RefundPayment(ctx context.Context, shopID int64, paymentID string, amount float64, reason string) (*models.PaymentResponse, error) {
	config, err := p.GetPaystackConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	baseURL := p.getBaseURL(config.TestMode)

	// Convert amount to kobo (if amount is provided, otherwise do full refund)
	var amountKobo int64
	if amount > 0 {
		amountKobo = int64(amount * 100)
	}

	// Create refund request
	refundReq := PaystackRefundRequest{
		Transaction: paymentID,
	}

	if amountKobo > 0 {
		refundReq.Amount = amountKobo
	}

	refundReqJSON, err := json.Marshal(refundReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal refund request: %w", err)
	}

	// Make request to create refund
	headers := map[string]string{
		"Authorization": "Bearer " + config.SecretKey,
	}

	resp, err := p.makePaystackRequest(ctx, "POST", baseURL+"/refund", headers, bytes.NewBuffer(refundReqJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResp PaystackErrorResponse
		if jsonErr := json.Unmarshal(respBody, &errorResp); jsonErr == nil {
			return nil, fmt.Errorf("paystack api error: %s", errorResp.Message)
		}
		return nil, fmt.Errorf("paystack api error: status %d", resp.StatusCode)
	}

	var refundResp PaystackRefundResponse
	if err := json.Unmarshal(respBody, &refundResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !refundResp.Status {
		return nil, fmt.Errorf("failed to create refund: %s", refundResp.Message)
	}

	// Convert amount from kobo to main currency unit
	refundAmount := float64(refundResp.Data.Deducted.Amount) / 100
	now := time.Now()

	response := &models.PaymentResponse{
		PaymentID:     paymentID,
		Status:        "refunded",
		Amount:        refundAmount,
		CurrencyCode:  "NGN", // Paystack primarily deals with NGN
		PaymentMethod: "paystack",
		TransactionID: refundResp.Data.Transaction.Reference,
		ProcessedAt:   now.Format(time.RFC3339),
		PaymentDetails: map[string]interface{}{
			"refund_reference": refundResp.Data.Transaction.Reference,
			"refund_status":    refundResp.Data.Transaction.Status,
			"full_refund":      refundResp.Data.FullRefund,
			"reason":           reason,
		},
	}

	return response, nil
}

// HandleWebhook processes webhook events from Paystack
func (p *PaystackService) HandleWebhook(ctx context.Context, payload []byte, signature string) (*models.PaymentWebhookPayload, error) {
	// TODO: Implement webhook signature validation
	// Paystack uses HMAC SHA512 for webhook verification

	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	webhookPayload := &models.PaymentWebhookPayload{
		Provider:   "paystack",
		EventType:  "",
		RawPayload: event,
		ReceivedAt: time.Now(),
	}

	// Extract event type
	if eventType, exists := event["event"].(string); exists {
		webhookPayload.EventType = eventType
	}

	// Extract payment information based on event type
	if data, exists := event["data"].(map[string]interface{}); exists {
		switch webhookPayload.EventType {
		case "charge.success", "charge.failed", "charge.pending":
			if reference, exists := data["reference"].(string); exists {
				webhookPayload.PaymentID = reference
				webhookPayload.TransactionID = reference
			}
			if amount, exists := data["amount"].(float64); exists {
				webhookPayload.Amount = amount / 100 // Convert from kobo
			}
			if currency, exists := data["currency"].(string); exists {
				webhookPayload.Currency = currency
			}
			if status, exists := data["status"].(string); exists {
				webhookPayload.Status = p.ConvertPaystackStatusToPaymentStatus(status)
			}
			if metadata, exists := data["metadata"].(map[string]interface{}); exists {
				webhookPayload.Metadata = metadata
			}
		}
	}

	return webhookPayload, nil
}
