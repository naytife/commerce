package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
)

func CartRouter(app fiber.Router, repo db.Repository) {
	handler := handlers.NewHandler(repo)

	app.Post("/cart", handler.CreateShop)
	app.Get("/cart", handler.GetShops)
	app.Delete("/cart/:id", handler.DeleteShop)
}
