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

	ic "github.com/petrejonn/naytife/internal/httpclient"

	"github.com/petrejonn/naytife/internal/observability"
	"go.uber.org/zap"
)

type FlutterwaveService struct {
	repository db.Repository
}

type FlutterwaveConfig struct {
	PublicKey     string `json:"public_key"`
	SecretKey     string `json:"secret_key"`
	EncryptionKey string `json:"encryption_key"`
	TestMode      bool   `json:"test_mode"`
}

// Flutterwave API structures
type FlutterwavePaymentLinkRequest struct {
	TxRef       string                 `json:"tx_ref"`
	Amount      string                 `json:"amount"`
	Currency    string                 `json:"currency"`
	RedirectURL string                 `json:"redirect_url,omitempty"`
	Customer    FlutterwaveCustomer    `json:"customer"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
}

type FlutterwaveCustomer struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phonenumber,omitempty"`
	Name        string `json:"name,omitempty"`
}

type FlutterwavePaymentLinkResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Link string `json:"link"`
	} `json:"data"`
}

type FlutterwavePaymentRequest struct {
	TxRef          string                    `json:"tx_ref"`
	Amount         string                    `json:"amount"`
	Currency       string                    `json:"currency"`
	RedirectURL    string                    `json:"redirect_url,omitempty"`
	Customer       FlutterwaveCustomer       `json:"customer"`
	Customizations FlutterwaveCustomizations `json:"customizations,omitempty"`
	PaymentOptions string                    `json:"payment_options,omitempty"`
	Meta           map[string]interface{}    `json:"meta,omitempty"`
}

type FlutterwaveCustomizations struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Logo        string `json:"logo,omitempty"`
}

type FlutterwavePaymentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Link string `json:"link"`
	} `json:"data"`
}

type FlutterwaveVerifyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID                int64   `json:"id"`
		TxRef             string  `json:"tx_ref"`
		FlwRef            string  `json:"flw_ref"`
		DeviceFingerprint string  `json:"device_fingerprint"`
		Amount            float64 `json:"amount"`
		Currency          string  `json:"currency"`
		ChargedAmount     float64 `json:"charged_amount"`
		AppFee            float64 `json:"app_fee"`
		MerchantFee       float64 `json:"merchant_fee"`
		ProcessorResponse string  `json:"processor_response"`
		AuthModel         string  `json:"auth_model"`
		IP                string  `json:"ip"`
		Narration         string  `json:"narration"`
		Status            string  `json:"status"`
		PaymentType       string  `json:"payment_type"`
		CreatedAt         string  `json:"created_at"`
		AccountID         int64   `json:"account_id"`
		Customer          struct {
			ID          int64  `json:"id"`
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
			Email       string `json:"email"`
			CreatedAt   string `json:"created_at"`
		} `json:"customer"`
		Card struct {
			First6Digits string `json:"first_6digits"`
			Last4Digits  string `json:"last_4digits"`
			Issuer       string `json:"issuer"`
			Country      string `json:"country"`
			Type         string `json:"type"`
			Token        string `json:"token"`
			Expiry       string `json:"expiry"`
		} `json:"card"`
		Meta map[string]interface{} `json:"meta"`
	} `json:"data"`
}

type FlutterwaveRefundRequest struct {
	Amount string `json:"amount,omitempty"`
}

type FlutterwaveRefundResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID             int64   `json:"id"`
		AccountID      int64   `json:"account_id"`
		TxID           int64   `json:"tx_id"`
		FlwRef         string  `json:"flw_ref"`
		WalletID       int64   `json:"wallet_id"`
		AmountRefunded float64 `json:"amount_refunded"`
		Status         string  `json:"status"`
		Destination    string  `json:"destination"`
		Meta           struct {
			Source string `json:"source"`
		} `json:"meta"`
		CreatedAt string `json:"created_at"`
	} `json:"data"`
}

type FlutterwaveErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewFlutterwaveService(repo db.Repository) *FlutterwaveService {
	return &FlutterwaveService{
		repository: repo,
	}
}

// GetFlutterwaveConfig retrieves Flutterwave configuration for a shop
func (f *FlutterwaveService) GetFlutterwaveConfig(ctx context.Context, shopID int64) (*FlutterwaveConfig, error) {
	paymentMethod, err := f.repository.GetShopPaymentMethod(ctx, db.GetShopPaymentMethodParams{
		ShopID:     shopID,
		MethodType: db.PaymentMethodTypeFlutterwave,
	})
	if err != nil {
		zap.L().Error("GetFlutterwaveConfig: failed to get Flutterwave config", zap.Int64("shop_id", shopID), zap.Error(err))
		return nil, fmt.Errorf("failed to get Flutterwave config: %w", err)
	}

	if !paymentMethod.IsEnabled {
		zap.L().Warn("GetFlutterwaveConfig: flutterwave not enabled for shop", zap.Int64("shop_id", shopID))
		return nil, fmt.Errorf("flutterwave is not enabled for this shop")
	}

	var config FlutterwaveConfig
	if err := json.Unmarshal(paymentMethod.Attributes, &config); err != nil {
		zap.L().Error("GetFlutterwaveConfig: failed to parse Flutterwave config", zap.Int64("shop_id", shopID), zap.Error(err))
		return nil, fmt.Errorf("failed to parse Flutterwave config: %w", err)
	}

	return &config, nil
}

// getBaseURL returns the appropriate Flutterwave API base URL
func (f *FlutterwaveService) getBaseURL(testMode bool) string {
	// Flutterwave uses the same URL for both test and live, differentiated by API keys
	return "https://api.flutterwave.com/v3"
}

// makeFlutterwaveRequest makes an HTTP request to Flutterwave API
func (f *FlutterwaveService) makeFlutterwaveRequest(ctx context.Context, method, url string, headers map[string]string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		zap.L().Error("makeFlutterwaveRequest: failed to create request", zap.String("method", method), zap.String("url", url), zap.Error(err))
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
	_, finish := observability.StartSpan(ctx, "makeFlutterwaveRequest", "flutterwave", method, url)
	defer finish(0, nil)
	start := time.Now()
	resp, err := ic.DefaultClient.Do(req)
	if err != nil {
		zap.L().Error("makeFlutterwaveRequest: request failed", zap.String("method", method), zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	observability.RecordServiceRequest("flutterwave", req.Method, req.URL.String(), resp.StatusCode, time.Since(start))

	return resp, nil
}

// ValidateConfig validates the Flutterwave configuration
func (f *FlutterwaveService) ValidateConfig(config map[string]interface{}) error {
	requiredFields := []string{"public_key", "secret_key", "encryption_key"}

	for _, field := range requiredFields {
		if value, exists := config[field]; !exists || value == "" {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validate key formats
	if publicKey, ok := config["public_key"].(string); ok {
		if !strings.HasPrefix(publicKey, "FLWPUBK_") && !strings.HasPrefix(publicKey, "FLWPUBK-") {
			return fmt.Errorf("invalid public key format")
		}
	}

	if secretKey, ok := config["secret_key"].(string); ok {
		if !strings.HasPrefix(secretKey, "FLWSECK_") && !strings.HasPrefix(secretKey, "FLWSECK-") {
			return fmt.Errorf("invalid secret key format")
		}
	}

	return nil
}

// ConvertFlutterwaveStatusToPaymentStatus converts Flutterwave transaction status to our payment status
func (f *FlutterwaveService) ConvertFlutterwaveStatusToPaymentStatus(status string) string {
	switch strings.ToLower(status) {
	case "successful":
		return "completed"
	case "pending":
		return "pending"
	case "failed":
		return "failed"
	case "cancelled":
		return "cancelled"
	default:
		return "pending"
	}
}

// generateTxRef generates a unique transaction reference
func (f *FlutterwaveService) generateTxRef(shopID int64, checkoutSessionID string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("flw_%d_%s_%d", shopID, checkoutSessionID, timestamp)
}

// ProcessPayment processes a payment request and returns the payment response
func (f *FlutterwaveService) ProcessPayment(ctx context.Context, shopID int64, req models.PaymentRequest, amount float64, currencyCode string) (*models.PaymentResponse, error) {
	config, err := f.GetFlutterwaveConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	baseURL := f.getBaseURL(config.TestMode)
	txRef := f.generateTxRef(shopID, req.CheckoutSessionID)

	// Create payment request
	paymentReq := FlutterwavePaymentRequest{
		TxRef:    txRef,
		Amount:   fmt.Sprintf("%.2f", amount),
		Currency: currencyCode,
		Customer: FlutterwaveCustomer{
			Email: "customer@example.com", // TODO: Get customer email from checkout session
		},
		PaymentOptions: "card,banktransfer,ussd,mobilemoney",
		Meta: map[string]interface{}{
			"shop_id":             strconv.FormatInt(shopID, 10),
			"checkout_session_id": req.CheckoutSessionID,
		},
	}

	// Add custom metadata from request payment details
	if req.PaymentDetails != nil {
		for key, value := range req.PaymentDetails {
			paymentReq.Meta[key] = value
		}
	}

	paymentReqJSON, err := json.Marshal(paymentReq)
	if err != nil {
		zap.L().Error("ProcessPayment: failed to marshal payment request", zap.Int64("shop_id", shopID), zap.String("tx_ref", txRef), zap.Error(err))
		return nil, fmt.Errorf("failed to marshal payment request: %w", err)
	}

	// Make request to create payment
	headers := map[string]string{
		"Authorization": "Bearer " + config.SecretKey,
	}

	resp, err := f.makeFlutterwaveRequest(ctx, "POST", baseURL+"/payments", headers, bytes.NewBuffer(paymentReqJSON))
	if err != nil {
		zap.L().Error("ProcessPayment: failed to create payment", zap.Int64("shop_id", shopID), zap.String("tx_ref", txRef), zap.Error(err))
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("ProcessPayment: failed to read response body", zap.Int64("shop_id", shopID), zap.String("tx_ref", txRef), zap.Error(err))
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp FlutterwaveErrorResponse
		if jsonErr := json.Unmarshal(respBody, &errorResp); jsonErr == nil {
			zap.L().Error("ProcessPayment: flutterwave api error", zap.Int64("shop_id", shopID), zap.String("tx_ref", txRef), zap.String("message", errorResp.Message))
			return nil, fmt.Errorf("flutterwave api error: %s", errorResp.Message)
		}
		zap.L().Error("ProcessPayment: flutterwave api non-200 response", zap.Int64("shop_id", shopID), zap.String("tx_ref", txRef), zap.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("flutterwave api error: status %d", resp.StatusCode)
	}

	var paymentResp FlutterwavePaymentResponse
	if err := json.Unmarshal(respBody, &paymentResp); err != nil {
		zap.L().Error("ProcessPayment: failed to unmarshal response", zap.Int64("shop_id", shopID), zap.String("tx_ref", txRef), zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if paymentResp.Status != "success" {
		zap.L().Error("ProcessPayment: payment creation failed", zap.Int64("shop_id", shopID), zap.String("tx_ref", txRef), zap.String("message", paymentResp.Message))
		return nil, fmt.Errorf("failed to create payment: %s", paymentResp.Message)
	}

	response := &models.PaymentResponse{
		PaymentID:     txRef,
		Status:        "pending",
		Amount:        amount,
		CurrencyCode:  currencyCode,
		PaymentMethod: "flutterwave",
		TransactionID: txRef,
		NextAction: &models.PaymentNextAction{
			Type: "redirect",
			Data: map[string]interface{}{
				"redirect_url": paymentResp.Data.Link,
			},
		},
	}

	return response, nil
}

// CreatePaymentIntent creates a payment intent for deferred payment processing
func (f *FlutterwaveService) CreatePaymentIntent(ctx context.Context, shopID int64, req models.PaymentIntentRequest) (*models.PaymentIntentResponse, error) {
	config, err := f.GetFlutterwaveConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	baseURL := f.getBaseURL(config.TestMode)
	txRef := f.generateTxRef(shopID, req.CheckoutSessionID)

	// Create payment request
	paymentReq := FlutterwavePaymentRequest{
		TxRef:    txRef,
		Amount:   fmt.Sprintf("%.2f", req.Amount),
		Currency: req.CurrencyCode,
		Customer: FlutterwaveCustomer{
			Email: "customer@example.com", // TODO: Get customer email from checkout session
		},
		PaymentOptions: "card,banktransfer,ussd,mobilemoney",
		Meta: map[string]interface{}{
			"shop_id":             strconv.FormatInt(shopID, 10),
			"checkout_session_id": req.CheckoutSessionID,
		},
	}

	// Add custom metadata from request
	if req.Metadata != nil {
		for key, value := range req.Metadata {
			paymentReq.Meta[key] = value
		}
	}

	paymentReqJSON, err := json.Marshal(paymentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payment request: %w", err)
	}

	// Make request to create payment
	headers := map[string]string{
		"Authorization": "Bearer " + config.SecretKey,
	}

	resp, err := f.makeFlutterwaveRequest(ctx, "POST", baseURL+"/payments", headers, bytes.NewBuffer(paymentReqJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp FlutterwaveErrorResponse
		if jsonErr := json.Unmarshal(respBody, &errorResp); jsonErr == nil {
			return nil, fmt.Errorf("flutterwave api error: %s", errorResp.Message)
		}
		return nil, fmt.Errorf("flutterwave api error: status %d", resp.StatusCode)
	}

	var paymentResp FlutterwavePaymentResponse
	if err := json.Unmarshal(respBody, &paymentResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if paymentResp.Status != "success" {
		return nil, fmt.Errorf("failed to create payment: %s", paymentResp.Message)
	}

	response := &models.PaymentIntentResponse{
		PaymentIntentID: txRef,
		ClientSecret:    txRef, // Flutterwave doesn't have client secrets like Stripe
		Status:          "pending",
		Amount:          req.Amount,
		CurrencyCode:    req.CurrencyCode,
		NextAction: &models.PaymentNextAction{
			Type: "redirect",
			Data: map[string]interface{}{
				"redirect_url": paymentResp.Data.Link,
			},
		},
	}

	return response, nil
}

// ConfirmPayment confirms a payment intent or transaction
func (f *FlutterwaveService) ConfirmPayment(ctx context.Context, shopID int64, paymentID string) (*models.PaymentResponse, error) {
	// In Flutterwave, payment confirmation is done by verifying the transaction
	return f.GetPaymentStatus(ctx, shopID, paymentID)
}

// GetPaymentStatus retrieves the current status of a payment
func (f *FlutterwaveService) GetPaymentStatus(ctx context.Context, shopID int64, paymentID string) (*models.PaymentResponse, error) {
	config, err := f.GetFlutterwaveConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	baseURL := f.getBaseURL(config.TestMode)

	// Make request to verify transaction
	headers := map[string]string{
		"Authorization": "Bearer " + config.SecretKey,
	}

	resp, err := f.makeFlutterwaveRequest(ctx, "GET", baseURL+"/transactions/verify_by_reference?tx_ref="+paymentID, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to verify transaction: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp FlutterwaveErrorResponse
		if jsonErr := json.Unmarshal(respBody, &errorResp); jsonErr == nil {
			return nil, fmt.Errorf("flutterwave api error: %s", errorResp.Message)
		}
		return nil, fmt.Errorf("flutterwave api error: status %d", resp.StatusCode)
	}

	var verifyResp FlutterwaveVerifyResponse
	if err := json.Unmarshal(respBody, &verifyResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if verifyResp.Status != "success" {
		return nil, fmt.Errorf("failed to verify transaction: %s", verifyResp.Message)
	}

	response := &models.PaymentResponse{
		PaymentID:     verifyResp.Data.TxRef,
		Status:        f.ConvertFlutterwaveStatusToPaymentStatus(verifyResp.Data.Status),
		Amount:        verifyResp.Data.Amount,
		CurrencyCode:  verifyResp.Data.Currency,
		PaymentMethod: "flutterwave",
		TransactionID: verifyResp.Data.FlwRef,
	}

	// Set processed time if payment was successful and we have a created_at timestamp
	if verifyResp.Data.Status == "successful" && verifyResp.Data.CreatedAt != "" {
		if processedAt, err := time.Parse(time.RFC3339, verifyResp.Data.CreatedAt); err == nil {
			response.ProcessedAt = processedAt.Format(time.RFC3339)
		}
	}

	return response, nil
}

// RefundPayment processes a refund for a completed payment
func (f *FlutterwaveService) RefundPayment(ctx context.Context, shopID int64, paymentID string, amount float64, reason string) (*models.PaymentResponse, error) {
	config, err := f.GetFlutterwaveConfig(ctx, shopID)
	if err != nil {
		return nil, err
	}

	// First, get the transaction details to find the transaction ID
	paymentResp, err := f.GetPaymentStatus(ctx, shopID, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment details: %w", err)
	}

	if paymentResp.Status != "completed" {
		return nil, fmt.Errorf("payment must be completed to process refund")
	}

	baseURL := f.getBaseURL(config.TestMode)

	// Create refund request
	refundReq := FlutterwaveRefundRequest{}
	if amount > 0 {
		refundReq.Amount = fmt.Sprintf("%.2f", amount)
	}

	refundReqJSON, err := json.Marshal(refundReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal refund request: %w", err)
	}

	// Make request to create refund using the Flutterwave reference
	headers := map[string]string{
		"Authorization": "Bearer " + config.SecretKey,
	}

	var refundURL string
	if paymentResp.TransactionID != "" {
		refundURL = fmt.Sprintf("%s/transactions/%s/refund", baseURL, paymentResp.TransactionID)
	} else {
		return nil, fmt.Errorf("transaction ID not found for refund")
	}

	resp, err := f.makeFlutterwaveRequest(ctx, "POST", refundURL, headers, bytes.NewBuffer(refundReqJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResp FlutterwaveErrorResponse
		if jsonErr := json.Unmarshal(respBody, &errorResp); jsonErr == nil {
			return nil, fmt.Errorf("flutterwave api error: %s", errorResp.Message)
		}
		return nil, fmt.Errorf("flutterwave api error: status %d", resp.StatusCode)
	}

	var refundResp FlutterwaveRefundResponse
	if err := json.Unmarshal(respBody, &refundResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if refundResp.Status != "success" {
		return nil, fmt.Errorf("failed to create refund: %s", refundResp.Message)
	}

	refundAmount := refundResp.Data.AmountRefunded
	now := time.Now()

	response := &models.PaymentResponse{
		PaymentID:     paymentID,
		Status:        "refunded",
		Amount:        refundAmount,
		CurrencyCode:  paymentResp.CurrencyCode,
		PaymentMethod: "flutterwave",
		TransactionID: refundResp.Data.FlwRef,
		ProcessedAt:   now.Format(time.RFC3339),
		PaymentDetails: map[string]interface{}{
			"refund_reference": refundResp.Data.FlwRef,
			"refund_status":    refundResp.Data.Status,
			"refund_id":        refundResp.Data.ID,
			"reason":           reason,
		},
	}

	return response, nil
}

// HandleWebhook processes webhook events from Flutterwave
func (f *FlutterwaveService) HandleWebhook(ctx context.Context, payload []byte, signature string) (*models.PaymentWebhookPayload, error) {
	// TODO: Implement webhook signature validation
	// Flutterwave uses webhook hash for verification

	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	webhookPayload := &models.PaymentWebhookPayload{
		Provider:   "flutterwave",
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
		case "charge.completed", "charge.failed", "charge.pending":
			if txRef, exists := data["tx_ref"].(string); exists {
				webhookPayload.PaymentID = txRef
				webhookPayload.TransactionID = txRef
			}
			if amount, exists := data["amount"].(float64); exists {
				webhookPayload.Amount = amount
			}
			if currency, exists := data["currency"].(string); exists {
				webhookPayload.Currency = currency
			}
			if status, exists := data["status"].(string); exists {
				webhookPayload.Status = f.ConvertFlutterwaveStatusToPaymentStatus(status)
			}
			if meta, exists := data["meta"].(map[string]interface{}); exists {
				webhookPayload.Metadata = meta
			}
		}
	}

	return webhookPayload, nil
}
