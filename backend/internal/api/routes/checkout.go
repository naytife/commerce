package routes

import (
	"github.com/gofiber/fiber/v2"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/services"
	"go.uber.org/zap"
)

func CheckoutRouter(app fiber.Router, repo db.Repository, logger *zap.Logger, retryClient *retryablehttp.Client, paymentProcessorFactory *services.PaymentProcessorFactory) {
	handler := handlers.NewHandlerWithPaymentFactory(repo, retryClient, paymentProcessorFactory)

	// Checkout endpoints
	app.Post("/shops/:shop_id/checkout", handler.InitiateCheckout)
	app.Post("/shops/:shop_id/payment", handler.ProcessPayment)
	app.Post("/shops/:shop_id/payment/intent", handler.CreatePaymentIntent)
}
