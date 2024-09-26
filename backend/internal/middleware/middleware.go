package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/db"
)

func ShopIDMiddlewareFiber(repo db.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract the host from the Fiber context
		host := c.Hostname()
		ctx := c.UserContext()

		// Query the database to get the shop_id using the domain (host)
		shopID, err := repo.GetShopIDByDomain(ctx, host)
		if err != nil {
			// Handle error (return unauthorized response)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid shop",
			})
		}

		// Set the shop_id in the session or handle accordingly
		if err := repo.SetShopIDInSession(ctx, shopID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to set shop_id",
			})
		}

		// Store the shopID in the context (Fiber uses c.Locals for local variables)
		c.Locals("shop_id", shopID)

		// Proceed to the next middleware/handler
		return c.Next()
	}
}
