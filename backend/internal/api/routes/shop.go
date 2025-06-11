package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/handlers"
	"github.com/petrejonn/naytife/internal/db"
)

func ShopRouter(app fiber.Router, repo db.Repository) {
	handler := handlers.NewHandler(repo)

	app.Post("/shops", handler.CreateShop)
	app.Get("/shops", handler.GetShops)
	app.Delete("/shops/:shop_id", handler.DeleteShop)
	app.Get("/shops/:shop_id", handler.GetShop)
	app.Put("/shops/:shop_id", handler.UpdateShop)
	app.Put("/shops/:shop_id/images", handler.UpdateShopImages)
	app.Get("/shops/subdomain/:subdomain", handler.GetShopBySubDomain)
	app.Get("/shops/check-subdomain/:subdomain", handler.CheckSubdomainAvailability)
	app.Get("/customerinfo", handler.GetCustomerByEmail)
}
