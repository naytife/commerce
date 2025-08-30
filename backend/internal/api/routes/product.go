package routes

import (
	"github.com/gofiber/fiber/v2"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/services"
	"go.uber.org/zap"
)

func ProductRouter(app fiber.Router, repo db.Repository, logger *zap.Logger, retryClient *retryablehttp.Client) {
	storeDeployerClient := services.NewStoreDeployerClient(retryClient, logger)
	handler := handlers.NewHandlerWithStoreDeployerClient(repo, logger, retryClient, storeDeployerClient)

	app.Post("/shops/:shop_id/product-types/:product_type_id/products", handler.CreateProduct)
	app.Get("/shops/:shop_id/product-types/:product_type_id/products", handler.GetProductsByType)
	app.Get("/shops/:shop_id/products", handler.GetProducts)
	app.Get("/shops/:shop_id/products/:product_id", handler.GetProduct)
	app.Put("/shops/:shop_id/products/:product_id", handler.UpdateProduct)
	app.Delete("/shops/:shop_id/products/:product_id", handler.DeleteProduct)

	// Product images routes
	app.Post("/shops/:shop_id/products/:product_id/images", handler.AddProductImage)
	app.Get("/shops/:shop_id/products/:product_id/images", handler.GetProductImages)
	app.Delete("/shops/:shop_id/products/:product_id/images/:image_id", handler.DeleteProductImage)
}
