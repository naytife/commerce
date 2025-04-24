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
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/routes"
	"github.com/petrejonn/naytife/internal/db"
	publicgraph "github.com/petrejonn/naytife/internal/gql/public"
	"github.com/petrejonn/naytife/internal/middleware"
)

// @title Naytife API Docs
// @version 1.0
// @description This is the Naytife API documentation
// @servers.url http://localhost:8000/v1
// securityDefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://auth.naytife.com/oauth2/token
// @authorizationUrl https://auth.naytife.com/oauth2/auth
// @securityDefinitions.apikey XUserIdAuth
// @in header
// @name X-User-Id
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
		ReadBufferSize: 8192,
		// Global custom error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default error response
			statusCode := fiber.StatusInternalServerError
			message := "An unexpected error occurred"

			// Handle specific error types
			if e, ok := err.(*fiber.Error); ok {
				statusCode = e.Code
				switch statusCode {
				case fiber.StatusBadRequest:
					message = "Invalid input data"
				case fiber.StatusNotFound:
					message = "Resource not found"
				case fiber.StatusUnauthorized:
					message = "Authentication required"
				case fiber.StatusForbidden:
					message = "Insufficient permissions"
				}
			}

			return api.ErrorResponse(c, statusCode, message, nil)
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,OPTIONS,DELETE,PUT,PATCH",
		AllowHeaders: "Content-Type,Authorization,X-User-Id",
	}))
	app.Use(logger.New())
	app.Get("/v1/docs/swagger.json", func(c *fiber.Ctx) error {
		return c.SendFile("docs/swagger.json")
	})
	app.Get("/v1/docs/*", swagger.New(swagger.Config{
		URL:         fmt.Sprintf("%s/v1/docs/swagger.json", env.API_URL),
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:      "Naytife API",
			ClientId:     "bcb8b621-1519-4127-b4d2-6187b48eba99",
			ClientSecret: "Zr1-2LBQgiR5.0SzF~AT8.rPut",
			Scopes:       []string{"openid", "offline", "hydra.openid", "introspect"},
			UseBasicAuthenticationWithAccessCodeGrant: true,
		},
		PersistAuthorization: true,
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: fmt.Sprintf("%s/v1/docs/oauth2-redirect.html", env.API_URL),
	}))

	v1 := app.Group("/v1")
	api := v1.Group("/", middleware.WebMiddlewareFiber())
	routes.AuthRouter(v1, repo)
	routes.ShopRouter(api, repo)
	routes.ProductTypeRouter(api, repo)
	routes.ProductRouter(api, repo)
	routes.AttributeRouter(api, repo)
	routes.UserRouter(api, repo)
	routes.CartRouter(api, repo)

	app.Get("/graph", publicgraph.NewPlaygroundHandler("/query"))

	graphql := app.Group("/query", middleware.ShopIDMiddlewareFiber(repo))
	graphql.Post("/", publicgraph.NewHandler(repo)) // public

	address := ":" + env.PORT
	fmt.Fprintf(os.Stdout, "🚀 Server ready at port %s\n", address)
	if err := app.Listen(address); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
