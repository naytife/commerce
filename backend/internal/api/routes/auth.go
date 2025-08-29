package routes

import (
	"github.com/gofiber/fiber/v2"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
	"go.uber.org/zap"
)

func AuthRouter(app fiber.Router, repo db.Repository, logger *zap.Logger, retryClient *retryablehttp.Client) {
	handler := handlers.NewHandler(repo, logger, retryClient)
	app.Post("/auth/register", handler.UpsertUser)
	app.Post("/auth/register-customer", handler.UpsertCustomer)
}
