package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/services"
)

func CheckoutRouter(app fiber.Router, repo db.Repository, paymentProcessorFactory *services.PaymentProcessorFactory) {
	handler := handlers.NewHandlerWithPaymentFactory(repo, paymentProcessorFactory)

	// Checkout endpoints
	app.Post("/shops/:shop_id/checkout", handler.InitiateCheckout)
	app.Post("/shops/:shop_id/payment", handler.ProcessPayment)
	app.Post("/shops/:shop_id/payment/intent", handler.CreatePaymentIntent)
}
