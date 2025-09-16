package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/petrejonn/naytife/config"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/routes"
	"github.com/petrejonn/naytife/internal/db"
	publicgraph "github.com/petrejonn/naytife/internal/gql/public"
	"github.com/petrejonn/naytife/internal/middleware"
	"github.com/petrejonn/naytife/internal/services"
	"go.uber.org/zap"

	// Import generated swagger docs package if present
	docs "github.com/petrejonn/naytife/docs"
)

// @title Naytife API Docs
// @version 1.0
// @description This is the Naytife API documentation
// @servers.url http://127.0.0.1:8080/v1
// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenurl http://127.0.0.1:8080/oauth2/token
// @authorizationurl http://127.0.0.1:8080/oauth2/auth
// @in header

func main() {
	env, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}
	defer logger.Sync()

	dbase, err := db.InitDB(env.DATABASE_URL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbase.Close()

	repo := db.NewRepository(dbase)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.RetryWaitMin = 500 * time.Millisecond
	retryClient.RetryWaitMax = 5 * time.Second
	retryClient.HTTPClient = &http.Client{
		Timeout: 15 * time.Second, // hard timeout per attempt
	}
	retryClient.Logger = nil // silence retryablehttp’s default logs (we use zap)

	// Initialize Redis client for publish functionality
	var redisClient *redis.Client
	if env.REDIS_URL != "" {
		// Parse Redis URL if provided
		opt, err := redis.ParseURL(env.REDIS_URL)
		if err != nil {
			logger.Fatal("Failed to parse Redis URL", zap.Error(err))
		} else {
			redisClient = redis.NewClient(opt)
			// Test Redis connection
			if _, err := redisClient.Ping(redisClient.Context()).Result(); err != nil {
				logger.Warn("Failed to connect to Redis", zap.Error(err))
				redisClient = nil
			} else {
				logger.Info("✅ Connected to Redis successfully!")
			}
		}
	}

	// Initialize services
	stripeService := services.NewStripeService(repo)
	paypalService := services.NewPayPalService(repo)
	paystackService := services.NewPaystackService(repo)
	flutterwaveService := services.NewFlutterwaveService(repo)

	// Wire retry client into services that perform outbound HTTP calls
	paypalService.RetryClient = retryClient
	paystackService.RetryClient = retryClient
	flutterwaveService.RetryClient = retryClient

	// Initialize payment processor factory
	paymentProcessorFactory := services.NewPaymentProcessorFactory(
		stripeService,
		paypalService,
		paystackService,
		flutterwaveService,
	)

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

	// Note: CORS is handled by Oathkeeper, so we don't need the CORS middleware here
	// Custom logging middleware using Zap
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next() // process request
		stop := time.Since(start)

		logger.Info("incoming request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", stop),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
		)

		return err
	})

	// Health check endpoints for Kubernetes
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"service":   "naytife-backend",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	app.Get("/ready", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ready",
			"service":   "naytife-backend",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	app.Get("/v1/docs/swagger.json", func(c *fiber.Ctx) error {
		// Start from generated swagger (by swag init)
		raw := docs.SwaggerInfo.ReadDoc()
		var spec map[string]any
		if err := json.Unmarshal([]byte(raw), &spec); err != nil {
			// Fallback to serving the static file if parsing fails
			return c.SendFile("docs/swagger.json")
		}

		// Ensure servers reflect current environment
		serverURL := fmt.Sprintf("%s/v1", env.API_URL)
		envName := env.ENV
		if envName == "" {
			envName = "default"
		}
		spec["servers"] = []map[string]string{{
			"url":         serverURL,
			"description": envName,
		}}

		// Update OAuth2 authorization/token URLs if available
		authBase := env.AUTH_URL
		if authBase == "" {
			// sensible default: use API_URL root if AUTH_URL not provided
			authBase = env.API_URL
		}

		// Navigate to components.securitySchemes to adjust OAuth flows
		if compAny, ok := spec["components"]; ok {
			if comp, ok := compAny.(map[string]any); ok {
				if secAny, ok := comp["securitySchemes"]; ok {
					if sec, ok := secAny.(map[string]any); ok {
						// Try common keys produced by swag for oauth2 accessCode/authorizationCode
						for _, key := range []string{"OAuth2AccessCode", "OAuth2", "oauth2"} {
							if schemeAny, ok := sec[key]; ok {
								if scheme, ok := schemeAny.(map[string]any); ok {
									if flowsAny, ok := scheme["flows"]; ok {
										if flows, ok := flowsAny.(map[string]any); ok {
											// authorizationCode flow is typical for access code
											if acAny, ok := flows["authorizationCode"]; ok {
												if ac, ok := acAny.(map[string]any); ok {
													ac["authorizationUrl"] = fmt.Sprintf("%s/oauth2/auth", authBase)
													ac["tokenUrl"] = fmt.Sprintf("%s/oauth2/token", authBase)
												}
											}
											// legacy accessCode (OpenAPI v2) mapping sometimes appears
											if acAny, ok := flows["accessCode"]; ok {
												if ac, ok := acAny.(map[string]any); ok {
													ac["authorizationUrl"] = fmt.Sprintf("%s/oauth2/auth", authBase)
													ac["tokenUrl"] = fmt.Sprintf("%s/oauth2/token", authBase)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}

		out, err := json.MarshalIndent(spec, "", "  ")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("failed to encode swagger spec")
		}
		c.Set("Content-Type", "application/json")
		return c.Send(out)
	})
	app.Get("/v1/docs/*", swagger.New(swagger.Config{
		URL:         fmt.Sprintf("%s/v1/docs/swagger.json", env.API_URL),
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "list",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:      "Naytife API",
			ClientId:     "d39beaaa-9c53-48e7-b82a-37ff52127473",
			ClientSecret: "-tzS7OuCyHjTZUxtfx5TxGR1f.",
			Scopes:       []string{"openid", "offline", "profile", "email", "offline_access"},
			UseBasicAuthenticationWithAccessCodeGrant: true,
			AdditionalQueryStringParams: map[string]string{
				"app_type": "dashboard",
			},
		},
		PersistAuthorization: true,
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: fmt.Sprintf("%s/v1/docs/oauth2-redirect.html", env.API_URL),
	}))

	v1 := app.Group("/v1")
	api := v1.Group("/", middleware.WebMiddlewareFiber())
	routes.AuthRouter(v1, repo, retryClient)
	routes.ShopRouter(api, repo, retryClient)
	routes.ProductTypeRouter(api, repo, retryClient)
	routes.ProductRouter(api, repo, retryClient)
	routes.AttributeRouter(api, repo, retryClient)
	routes.UserRouter(api, repo, retryClient)
	routes.CheckoutRouter(api, repo, retryClient, paymentProcessorFactory)
	routes.PaymentRouter(api, repo, paymentProcessorFactory)
	routes.PaymentMethodsRouter(api, repo, retryClient)
	routes.OrderRouter(api, repo, retryClient)
	routes.CustomerRouter(api, repo, retryClient)
	routes.InventoryRouter(api, repo, retryClient)
	routes.AnalyticsRouter(api, repo)
	routes.TemplateRouter(api, repo, retryClient)
	routes.WebhookRouter(v1, repo, paymentProcessorFactory)

	app.Get("/graph", publicgraph.NewPlaygroundHandler("/query"))

	graphql := app.Group("/query", middleware.ShopIDMiddlewareFiber(repo))
	graphql.Post("/", publicgraph.NewHandler(repo)) // public

	address := ":" + env.PORT
	fmt.Fprintf(os.Stdout, "🚀 Server ready at port %s\n", address)
	if err := app.Listen(address); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
