package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
)

func CustomerRouter(app fiber.Router, repo db.Repository) {
	handler := handlers.NewHandler(repo)

	// Customer management endpoints
	app.Get("/shops/:shop_id/customers", handler.GetCustomers)
	app.Get("/shops/:shop_id/customers/search", handler.SearchCustomers)
	app.Get("/shops/:shop_id/customers/:customer_id", handler.GetCustomerById)
	app.Put("/shops/:shop_id/customers/:customer_id", handler.UpdateCustomer)
	app.Delete("/shops/:shop_id/customers/:customer_id", handler.DeleteCustomer)
	app.Get("/shops/:shop_id/customers/:customer_id/orders", handler.GetCustomerOrders)
}
