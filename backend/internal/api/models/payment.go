package models

import (
	"time"
)

// PaymentRequest represents a payment request from the client
type PaymentRequest struct {
	CheckoutSessionID string                 `json:"checkout_session_id" validate:"required"`
	PaymentMethod     string                 `json:"payment_method" validate:"required,oneof=stripe paypal paystack flutterwave cash_on_delivery"`
	PaymentDetails    map[string]interface{} `json:"payment_details,omitempty"`
	SavePaymentMethod bool                   `json:"save_payment_method,omitempty"`
}

// PaymentResponse represents a payment response to the client
type PaymentResponse struct {
	PaymentID      string                 `json:"payment_id"`
	OrderID        int64                  `json:"order_id"`
	Status         string                 `json:"status"`
	Amount         float64                `json:"amount"`
	CurrencyCode   string                 `json:"currency_code"`
	PaymentMethod  string                 `json:"payment_method"`
	TransactionID  string                 `json:"transaction_id,omitempty"`
	ProcessedAt    string                 `json:"processed_at,omitempty"`
	PaymentDetails map[string]interface{} `json:"payment_details,omitempty"`
	NextAction     *PaymentNextAction     `json:"next_action,omitempty"`
}

// PaymentNextAction represents additional actions required for payment completion
type PaymentNextAction struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// PaymentMethodInfo represents available payment method information
type PaymentMethodInfo struct {
	Type        string                 `json:"type"`
	DisplayName string                 `json:"display_name"`
	Enabled     bool                   `json:"enabled"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// PaymentWebhookPayload represents a webhook payload from payment providers
type PaymentWebhookPayload struct {
	Provider      string                 `json:"provider"`
	EventType     string                 `json:"event_type"`
	PaymentID     string                 `json:"payment_id"`
	TransactionID string                 `json:"transaction_id,omitempty"`
	Status        string                 `json:"status"`
	Amount        float64                `json:"amount,omitempty"`
	Currency      string                 `json:"currency,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	RawPayload    interface{}            `json:"raw_payload,omitempty"`
	ReceivedAt    time.Time              `json:"received_at"`
}

// PaymentIntentRequest represents a request to create a payment intent
type PaymentIntentRequest struct {
	Amount            float64                `json:"amount" validate:"required,min=0.01"`
	CurrencyCode      string                 `json:"currency_code" validate:"required"`
	PaymentMethod     string                 `json:"payment_method" validate:"required,oneof=stripe paypal paystack flutterwave cash_on_delivery"`
	CheckoutSessionID string                 `json:"checkout_session_id" validate:"required"`
	PaymentMethodID   string                 `json:"payment_method_id,omitempty"`
	CustomerID        string                 `json:"customer_id,omitempty"`
	Description       string                 `json:"description,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// PaymentIntentResponse represents a payment intent response
type PaymentIntentResponse struct {
	PaymentIntentID string             `json:"payment_intent_id"`
	ClientSecret    string             `json:"client_secret,omitempty"`
	Status          string             `json:"status"`
	Amount          float64            `json:"amount"`
	CurrencyCode    string             `json:"currency_code"`
	NextAction      *PaymentNextAction `json:"next_action,omitempty"`
}
