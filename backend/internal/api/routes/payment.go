package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/services"
)

// PaymentRouter sets up payment routes
func PaymentRouter(app fiber.Router, repo db.Repository, paymentProcessorFactory *services.PaymentProcessorFactory) {
	// Create payment handler
	paymentHandler := handlers.NewPaymentHandler(paymentProcessorFactory, repo)

	// Payment routes group
	payments := app.Group("/payments")

	// Payment session endpoints
	payments.Post("/checkout", paymentHandler.CreateCheckoutSession)
	payments.Post("/:shop_id/confirm", paymentHandler.ConfirmPayment)
	payments.Get("/:shop_id/status/:payment_intent_id", paymentHandler.GetPaymentStatus)
}
