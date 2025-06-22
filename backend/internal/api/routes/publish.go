package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
)

func PublishRouter(app fiber.Router, repo db.Repository, redisClient *redis.Client) {
	var handler *handlers.Handler

	// Create handler with Redis if available, otherwise create basic handler
	if redisClient != nil {
		handler = handlers.NewHandlerWithRedis(repo, redisClient)
	} else {
		handler = handlers.NewHandler(repo)
	}

	// Publish endpoints
	app.Post("/shops/:shop_id/publish", handler.TriggerBuild)
	app.Get("/shops/:shop_id/publish/status", handler.GetBuildStatus)
	app.Get("/shops/:shop_id/publish/history", handler.GetBuildHistory)
}
