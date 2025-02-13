package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
)

func ProductTypeRouter(app fiber.Router, repo db.Repository) {
	handler := handlers.NewHandler(repo)

	app.Post("/shops/:shop_id/product-types", handler.CreateProductType)
	app.Get("/shops/:shop_id/product-types", handler.GetProductTypes)
	app.Get("/shops/:shop_id/product-types/:product_type_id", handler.GetProductType)
	app.Put("/shops/:shop_id/product-types/:product_type_id", handler.UpdateProductType)
	app.Delete("/shops/:shop_id/product-types/:product_type_id", handler.DeleteProductType)
}
