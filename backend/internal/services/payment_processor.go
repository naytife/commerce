package services

import (
	"context"

	"github.com/petrejonn/naytife/internal/api/models"
)

// PaymentProcessor defines the interface that all payment service providers must implement
type PaymentProcessor interface {
	// ProcessPayment processes a payment request and returns the payment response
	ProcessPayment(ctx context.Context, shopID int64, req models.PaymentRequest, amount float64, currencyCode string) (*models.PaymentResponse, error)

	// CreatePaymentIntent creates a payment intent for deferred payment processing
	CreatePaymentIntent(ctx context.Context, shopID int64, req models.PaymentIntentRequest) (*models.PaymentIntentResponse, error)

	// ConfirmPayment confirms a payment intent or transaction
	ConfirmPayment(ctx context.Context, shopID int64, paymentID string) (*models.PaymentResponse, error)

	// GetPaymentStatus retrieves the current status of a payment
	GetPaymentStatus(ctx context.Context, shopID int64, paymentID string) (*models.PaymentResponse, error)

	// RefundPayment processes a refund for a completed payment
	RefundPayment(ctx context.Context, shopID int64, paymentID string, amount float64, reason string) (*models.PaymentResponse, error)

	// HandleWebhook processes webhook events from the payment provider
	HandleWebhook(ctx context.Context, payload []byte, signature string) (*models.PaymentWebhookPayload, error)

	// ValidateConfig validates the payment provider configuration
	ValidateConfig(config map[string]interface{}) error
}

// PaymentProcessorFactory creates payment processors based on provider type
type PaymentProcessorFactory struct {
	stripe      PaymentProcessor
	paypal      PaymentProcessor
	paystack    PaymentProcessor
	flutterwave PaymentProcessor
}

// NewPaymentProcessorFactory creates a new payment processor factory
func NewPaymentProcessorFactory(
	stripe PaymentProcessor,
	paypal PaymentProcessor,
	paystack PaymentProcessor,
	flutterwave PaymentProcessor,
) *PaymentProcessorFactory {
	return &PaymentProcessorFactory{
		stripe:      stripe,
		paypal:      paypal,
		paystack:    paystack,
		flutterwave: flutterwave,
	}
}

// GetProcessor returns the appropriate payment processor for the given provider
func (f *PaymentProcessorFactory) GetProcessor(provider string) PaymentProcessor {
	switch provider {
	case "stripe":
		return f.stripe
	case "paypal":
		return f.paypal
	case "paystack":
		return f.paystack
	case "flutterwave":
		return f.flutterwave
	default:
		return nil
	}
}

// GetSupportedProviders returns a list of supported payment providers
func (f *PaymentProcessorFactory) GetSupportedProviders() []string {
	return []string{"stripe", "paypal", "paystack", "flutterwave"}
}
