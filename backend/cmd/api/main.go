package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/petrejonn/naytife/config"
	"github.com/petrejonn/naytife/internal/api/routes"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/graph"
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

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	graphql := app.Group("/query", middleware.ShopIDMiddlewareFiber(repo))
	graphql.All("/", graph.NewHandler(repo))

	rest := app.Group("/v1")
	routes.ShopRouter(rest, repo)

	app.Get("/", graph.NewPlaygroundHandler("/api/query"))

	address := ":" + env.PORT
	fmt.Fprintf(os.Stdout, "🚀 Server ready at port %s\n", address)
	if err := app.Listen(address); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
