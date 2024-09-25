package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/petrejonn/naytife/config"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/graph"
	"github.com/petrejonn/naytife/internal/middleware"
)

func main() {
	// Load environment configuration
	env, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Initialize database connection
	dbase, err := db.InitDB(env.DATABASE_URL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbase.Close()

	// Initialize repository
	repo := db.NewRepository(dbase)

	// Set up Fiber app
	app := fiber.New()

	// Add some middlewares (CORS, logger)
	app.Use(cors.New())
	app.Use(logger.New())

	// JWT and ShopID middlewares
	api := app.Group("/api", middleware.ShopIDMiddlewareFiber(repo))

	// GraphQL handlers (using the same resolver logic)
	api.All("/query", graph.NewHandler(repo))

	// Playground for testing
	app.Get("/", graph.NewPlaygroundHandler("/api/query"))

	// REST endpoint: Add a store
	api.Post("/store", func(c *fiber.Ctx) error {
		// Your logic to add a store
		return c.JSON(fiber.Map{
			"message": "Store added successfully",
		})
	})
	app.Use(swagger.New())

	// Start the server
	address := ":" + env.PORT
	fmt.Fprintf(os.Stdout, "🚀 Server ready at port %s\n", address)
	if err := app.Listen(address); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
