package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
)

func ProductRouter(app fiber.Router, repo db.Repository) {
	handler := handlers.NewHandler(repo)

	app.Post("/shops/:shop_id/product-types/:product_type_id/products", handler.CreateProduct)
	app.Get("/shops/:shop_id/products", handler.GetProducts)
	app.Get("/shops/:shop_id/products/:product_id", handler.GetProduct)
	app.Put("/shops/:shop_id/products/:product_id", handler.UpdateProduct)
	app.Delete("/shops/:shop_id/products/:product_id", handler.DeleteProduct)
}
