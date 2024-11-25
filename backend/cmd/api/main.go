package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/petrejonn/naytife/config"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/api/routes"
	"github.com/petrejonn/naytife/internal/db"
	admingraph "github.com/petrejonn/naytife/internal/gql/admin"
	publicgraph "github.com/petrejonn/naytife/internal/gql/public"
	"github.com/petrejonn/naytife/internal/middleware"
)

func main() {
	env, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	dbase, err := db.InitDB(env.DATABASE_URL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbase.Close()

	repo := db.NewRepository(dbase)

	app := fiber.New(fiber.Config{
		// Global custom error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
			})
		},
	})

	app.Use(cors.New())
	app.Use(logger.New())

	api := app.Group("/api/v1", middleware.WebMiddlewareFiber())
	routes.ShopRouter(api, repo)
	routes.UserRouter(api, repo)
	routes.CartRouter(api, repo)

	app.Get("/api/graph", publicgraph.NewPlaygroundHandler("/api/query"))
	app.Get("/api/admin/graph", admingraph.NewPlaygroundHandler("/api/admin/query"))

	graphql := app.Group("/api", middleware.ShopIDMiddlewareFiber(repo))
	graphql.Post("/query", publicgraph.NewHandler(repo))      // public
	graphql.Post("/admin/query", admingraph.NewHandler(repo)) // admin

	address := ":" + env.PORT
	fmt.Fprintf(os.Stdout, "🚀 Server ready at port %s\n", address)
	if err := app.Listen(address); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
