package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
)

func AnalyticsRouter(app fiber.Router, repo db.Repository) {
	handler := handlers.NewAnalyticsHandler(repo)

	app.Get("/shops/:shop_id/analytics/sales-summary", handler.GetSalesSummary)
	app.Get("/shops/:shop_id/analytics/orders-over-time", handler.GetOrdersOverTime)
	app.Get("/shops/:shop_id/analytics/top-products", handler.GetTopProducts)
	app.Get("/shops/:shop_id/analytics/customers-summary", handler.GetCustomerSummary)
	app.Get("/shops/:shop_id/analytics/low-stock", handler.GetLowStockProducts)
}
