package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/services"
)

// WebhookRouter sets up webhook routes for all payment providers
func WebhookRouter(app fiber.Router, repo db.Repository, paymentProcessorFactory *services.PaymentProcessorFactory) {
	// Create webhook handler
	webhookHandler := handlers.NewWebhookHandler(paymentProcessorFactory, repo)

	// Webhook routes group
	webhooks := app.Group("/webhooks")

	// Stripe webhook endpoint
	webhooks.Post("/stripe/:shop_id", webhookHandler.StripeWebhook)

	// PayPal webhook endpoint
	webhooks.Post("/paypal/:shop_id", webhookHandler.PayPalWebhook)

	// Paystack webhook endpoint
	webhooks.Post("/paystack/:shop_id", webhookHandler.PaystackWebhook)

	// Flutterwave webhook endpoint
	webhooks.Post("/flutterwave/:shop_id", webhookHandler.FlutterwaveWebhook)
}
