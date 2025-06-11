package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
)

func ProductTypeRouter(app fiber.Router, repo db.Repository) {
	handler := handlers.NewHandler(repo)

	// Predefined product type templates (public endpoints)
	app.Get("/predefined-product-types", handler.GetPredefinedProductTypes)
	app.Get("/predefined-product-types/:template_id", handler.GetPredefinedProductType)

	// Shop-specific product type endpoints
	app.Post("/shops/:shop_id/product-types", handler.CreateProductType)
	app.Post("/shops/:shop_id/product-types/from-template", handler.CreateProductTypeFromTemplate)
	app.Get("/shops/:shop_id/product-types", handler.GetProductTypes)
	app.Get("/shops/:shop_id/product-types/:product_type_id", handler.GetProductType)
	app.Put("/shops/:shop_id/product-types/:product_type_id", handler.UpdateProductType)
	app.Delete("/shops/:shop_id/product-types/:product_type_id", handler.DeleteProductType)
}
