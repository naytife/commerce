package routes

import (
	"github.com/gofiber/fiber/v2"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
	"go.uber.org/zap"
)

func CustomerRouter(app fiber.Router, repo db.Repository, logger *zap.Logger, retryClient *retryablehttp.Client) {
	handler := handlers.NewHandler(repo, retryClient)

	// Customer management endpoints
	app.Get("/shops/:shop_id/customers", handler.GetCustomers)
	app.Get("/shops/:shop_id/customers/search", handler.SearchCustomers)
	app.Get("/shops/:shop_id/customers/:customer_id", handler.GetCustomerById)
	app.Put("/shops/:shop_id/customers/:customer_id", handler.UpdateCustomer)
	app.Delete("/shops/:shop_id/customers/:customer_id", handler.DeleteCustomer)
	app.Get("/shops/:shop_id/customers/:customer_id/orders", handler.GetCustomerOrders)
}
