package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/petrejonn/naytife/config"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/api/routes"
	"github.com/petrejonn/naytife/internal/db"
	admingraph "github.com/petrejonn/naytife/internal/gql/admin"
	publicgraph "github.com/petrejonn/naytife/internal/gql/public"
	"github.com/petrejonn/naytife/internal/middleware"
)

// @title Naytife API Docs
// @version 1.0
// @description This is the Naytife API documentation
// @host 127.0.0.1:8080
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl http://127.0.0.1:8080/oauth2/token
// @authorizationUrl http://127.0.0.1:8080/oauth2/auth
// @scope.openid "OpenID scope"
// @scope.profile "Profile scope"
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

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,OPTIONS",
		AllowHeaders: "Content-Type,Authorization",
	}))
	app.Use(logger.New())
	app.Get("/api/v1/docs/swagger.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json") // or your own handler logic
	})
	app.Get("/api/v1/docs/*", swagger.New(swagger.Config{ // custom
		URL:         "http://127.0.0.1:8080/api/v1/docs/swagger.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:      "OAuth Provider",
			ClientId:     "f748b964-7de3-4159-9419-ae6e78502dc1",
			ClientSecret: "_6EHdptai8.SwLKodQxhpv3SKm",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:8080/api/v1/docs/oauth2-redirect.html",
	}))

	v1 := app.Group("/api/v1/")
	api := v1.Group("/*", middleware.WebMiddlewareFiber())
	routes.ShopRouter(api, repo)
	routes.UserRouter(api, repo)
	routes.CartRouter(api, repo)

	app.Get("/api/graph", publicgraph.NewPlaygroundHandler("/api/query"))
	app.Get("/api/admin/graph", admingraph.NewPlaygroundHandler("/api/admin/query"))

	graphql := app.Group("/api/query", middleware.ShopIDMiddlewareFiber(repo))
	graphql.Post("/", publicgraph.NewHandler(repo))           // public
	graphql.Post("/admin/query", admingraph.NewHandler(repo)) // admin

	address := ":" + env.PORT
	fmt.Fprintf(os.Stdout, "🚀 Server ready at port %s\n", address)
	if err := app.Listen(address); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
