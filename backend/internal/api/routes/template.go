package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/services"
)

func TemplateRouter(app fiber.Router, repo db.Repository, retryClient *retryablehttp.Client) {
	// Create both handlers - proxy for microservices, template for local logic
	storeDeployerClient := services.NewStoreDeployerClient(retryClient)
	proxyHandler := handlers.NewProxyHandler(repo)
	// attach retry client and store-deployer client to proxy handler
	proxyHandler.RetryClient = retryClient
	proxyHandler.StoreDeployerClient = storeDeployerClient
	templateHandler := handlers.NewTemplateHandler(repo)
	templateHandler.RetryClient = retryClient

	// Template management endpoints (proxied to template-registry)
	app.Get("/templates", proxyHandler.ProxyListTemplates)
	app.Get("/templates/:name", proxyHandler.ProxyGetTemplate)
	app.Get("/templates/:name/versions", proxyHandler.ProxyGetTemplateVersions)
	app.Get("/templates/:name/latest", proxyHandler.ProxyGetLatestTemplateVersion)
	app.Get("/templates/:name/versions/:version", proxyHandler.ProxyGetTemplateVersion)
	app.Get("/templates/:name/versions/:version/download", proxyHandler.ProxyDownloadTemplate)
	app.Post("/templates/upload", proxyHandler.ProxyUploadTemplate)

	// Store deployment endpoints (proxied to store-deployer)
	app.Post("/shops/:shop_id/deploy", proxyHandler.ProxyDeployStore)
	app.Post("/shops/:shop_id/redeploy", proxyHandler.ProxyRedeployStore)
	app.Get("/shops/:shop_id/deployment-status", templateHandler.GetDeploymentStatus)
	app.Post("/shops/:shop_id/update-data", proxyHandler.ProxyUpdateStoreData)
	app.Delete("/shops/:shop_id/cleanup", proxyHandler.ProxyCleanupStore)

	// Internal deployment callback endpoint
	app.Post("/internal/deployments/:shop_id/complete", templateHandler.CompleteDeploymentCallback)

	// Template version management endpoints (local handlers)
	app.Get("/shops/:shop_id/template/current", templateHandler.GetCurrentTemplate)
	app.Post("/shops/:shop_id/template/update-latest", templateHandler.UpdateToLatestTemplate)

	// Health check for services
	app.Get("/health/services", proxyHandler.ProxyHealthCheck)
}
