package routes

import (
	"github.com/gofiber/fiber/v2"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
	"go.uber.org/zap"
)

func OrderRouter(app fiber.Router, repo db.Repository, logger *zap.Logger, retryClient *retryablehttp.Client) {
	handler := handlers.NewHandler(repo, logger, retryClient)

	// Order management endpoints
	app.Post("/shops/:shop_id/orders", handler.CreateOrder)
	app.Get("/shops/:shop_id/orders", handler.GetOrders)
	app.Get("/shops/:shop_id/orders/:order_id", handler.GetOrder)
	app.Put("/shops/:shop_id/orders/:order_id", handler.UpdateOrder)
	app.Patch("/shops/:shop_id/orders/:order_id/status", handler.UpdateOrderStatus)
	app.Delete("/shops/:shop_id/orders/:order_id", handler.DeleteOrder)
}
