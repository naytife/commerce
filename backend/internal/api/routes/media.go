package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
)

func MediaRouter(app fiber.Router, repo db.Repository) {
	handler := handlers.NewHandler(repo)

	app.Post("/images/presigned_url", handler.CreateShop)
	// app.Delete("/shops/:id", handler.DeleteShop)
}
