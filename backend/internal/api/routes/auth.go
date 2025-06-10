package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
)

func AuthRouter(app fiber.Router, repo db.Repository) {
	handler := handlers.NewHandler(repo)
	app.Post("/auth/register", handler.UpsertUser)
	app.Post("/auth/register-customer", handler.UpsertCustomer)
}
