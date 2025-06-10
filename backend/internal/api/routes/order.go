package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
)

func OrderRouter(app fiber.Router, repo db.Repository) {
	handler := handlers.NewHandler(repo)

	// Order management endpoints
	app.Post("/shops/:shop_id/orders", handler.CreateOrder)
	app.Get("/shops/:shop_id/orders", handler.GetOrders)
	app.Get("/shops/:shop_id/orders/:order_id", handler.GetOrder)
	app.Put("/shops/:shop_id/orders/:order_id", handler.UpdateOrder)
	app.Patch("/shops/:shop_id/orders/:order_id/status", handler.UpdateOrderStatus)
	app.Delete("/shops/:shop_id/orders/:order_id", handler.DeleteOrder)
}
