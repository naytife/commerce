package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
)

func AttributeRouter(app fiber.Router, repo db.Repository) {
	handler := handlers.NewHandler(repo)

	app.Post("/shops/:shop_id/product-types/:product_type_id/attributes", handler.CreateAttribute)
	app.Get("/shops/:shop_id/product-types/:product_type_id/attributes", handler.GetAttributes)
	app.Get("/shops/:shop_id/attributes/:attribute_id", handler.GetAttribute)
	app.Put("/shops/:shop_id/attributes/:attribute_id", handler.UpdateAttribute)
	app.Delete("/shops/:shop_id/attributes/:attribute_id", handler.DeleteAttribute)
	app.Post("/shops/:shop_id/attributes/:attribute_id/options", handler.CreateAttributeOption)
	app.Get("/shops/:shop_id/attributes/:attribute_id/options", handler.GetAttributeOptions)
	app.Put("/shops/:shop_id/attribute-options/:attribute_option_id", handler.UpdateAttributeOption)
	app.Delete("/shops/:shop_id/attribute-options/:attribute_option_id", handler.DeleteAttributeOption)
}
