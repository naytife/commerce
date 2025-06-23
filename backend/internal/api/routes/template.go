package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
)

func TemplateRouter(app fiber.Router, repo db.Repository) {
	// Create both handlers - proxy for microservices, template for local logic
	proxyHandler := handlers.NewProxyHandler(repo)
	templateHandler := handlers.NewTemplateHandler(repo)

	// Template management endpoints (proxied to template-registry)
	app.Get("/templates", proxyHandler.ProxyListTemplates)
	app.Get("/templates/:name", proxyHandler.ProxyGetTemplate)
	app.Get("/templates/:name/versions", proxyHandler.ProxyGetTemplateVersions)
	app.Get("/templates/:name/latest", proxyHandler.ProxyGetLatestTemplateVersion)
	app.Get("/templates/:name/versions/:version", proxyHandler.ProxyGetTemplateVersion)
	app.Get("/templates/:name/versions/:version/download", proxyHandler.ProxyDownloadTemplate)
	app.Post("/templates/upload", proxyHandler.ProxyUploadTemplate)

	// Keep local template handlers for additional functionality if needed
	app.Post("/templates/build", templateHandler.BuildTemplate)

	// Store deployment endpoints (proxied to store-deployer)
	app.Post("/shops/:shop_id/deploy", proxyHandler.ProxyDeployStore)
	app.Post("/shops/:shop_id/redeploy", proxyHandler.ProxyRedeployStore)
	app.Get("/shops/:shop_id/deployment-status", proxyHandler.ProxyDeploymentStatus)
	app.Post("/shops/:shop_id/update-data", proxyHandler.ProxyUpdateStoreData)
	app.Delete("/shops/:shop_id/cleanup", proxyHandler.ProxyCleanupStore)

	// Health check for services
	app.Get("/health/services", proxyHandler.ProxyHealthCheck)
}
