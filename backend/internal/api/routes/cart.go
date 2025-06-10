package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
)

func CartRouter(app fiber.Router, repo db.Repository) {
	handler := handlers.NewHandler(repo)

	// Cart management endpoints
	app.Get("/cart", handler.GetCart)
	app.Delete("/cart", handler.ClearCart)

	// Cart item management endpoints
	app.Post("/cart/items", handler.AddToCart)
	app.Put("/cart/items/:item_id", handler.UpdateCartItem)
	app.Delete("/cart/items/:item_id", handler.RemoveFromCart)
}
