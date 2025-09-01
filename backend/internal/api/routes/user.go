package routes

import (
	"github.com/gofiber/fiber/v2"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
	"go.uber.org/zap"
)

func UserRouter(app fiber.Router, repo db.Repository, logger *zap.Logger, retryClient *retryablehttp.Client) {
	handler := handlers.NewHandler(repo, retryClient)
	app.Get("/me", handler.GetMe)
	app.Get("/userinfo", handler.GetUser)
}
