package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
)

func InventoryRouter(app fiber.Router, repo db.Repository) {
	handler := handlers.NewHandler(repo)

	// Inventory management endpoints
	app.Get("/shops/:shop_id/inventory", handler.GetInventoryReport) // General inventory endpoint
	app.Get("/shops/:shop_id/inventory/low-stock", handler.GetLowStockVariants)
	app.Put("/shops/:shop_id/inventory/variants/:variant_id/stock", handler.UpdateVariantStock)
	app.Post("/shops/:shop_id/inventory/variants/:variant_id/add-stock", handler.AddVariantStock)
	app.Post("/shops/:shop_id/inventory/variants/:variant_id/deduct-stock", handler.DeductVariantStock)
	app.Get("/shops/:shop_id/inventory/report", handler.GetInventoryReport)
	app.Get("/shops/:shop_id/inventory/movements", handler.GetStockMovements)
}
