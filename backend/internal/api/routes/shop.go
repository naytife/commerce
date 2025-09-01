package routes

import (
	"github.com/gofiber/fiber/v2"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/services"
	"go.uber.org/zap"
)

func ShopRouter(app fiber.Router, repo db.Repository, logger *zap.Logger, retryClient *retryablehttp.Client) {
	storeDeployerClient := services.NewStoreDeployerClient(retryClient)
	handler := handlers.NewHandlerWithStoreDeployerClient(repo, retryClient, storeDeployerClient)

	app.Post("/shops", handler.CreateShop)
	app.Get("/shops", handler.GetShops)
	app.Delete("/shops/:shop_id", handler.DeleteShop)
	app.Get("/shops/:shop_id", handler.GetShop)
	app.Put("/shops/:shop_id", handler.UpdateShop)
	app.Put("/shops/:shop_id/images", handler.UpdateShopImages)
	app.Get("/subdomains/:subdomain", handler.GetShopBySubDomain)
	app.Get("/subdomains/:subdomain/check", handler.CheckSubdomainAvailability)
	app.Get("/customerinfo", handler.GetCustomerByEmail)
}
