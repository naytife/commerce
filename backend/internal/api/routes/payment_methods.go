package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
)

func PaymentMethodsRouter(app fiber.Router, repo db.Repository) {
	handler := handlers.NewHandler(repo)

	// Payment methods endpoints
	app.Get("/shops/:shop_id/payment-methods", handler.GetShopPaymentMethods)
	app.Put("/shops/:shop_id/payment-methods/:method_type", handler.UpsertShopPaymentMethod)
	app.Patch("/shops/:shop_id/payment-methods/:method_type/status", handler.UpdateShopPaymentMethodStatus)
	app.Delete("/shops/:shop_id/payment-methods/:method_type", handler.DeleteShopPaymentMethod)

	// Payment testing endpoints
	app.Post("/shops/:shop_id/payment-methods/:method_type/test", handler.TestPaymentMethod)
}
