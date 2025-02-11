package middleware

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/db"
)

func ShopIDMiddlewareFiber(repo db.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		subdomain := strings.Split(c.Hostname(), ".")[0]
		ctx := c.UserContext()

		shopID, err := repo.GetShopIDByDomain(ctx, subdomain)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid shop",
			})
		}

		if err := repo.SetShopIDInSession(ctx, shopID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to set shop_id",
			})
		}

		c.Locals("shop_id", shopID)

		return c.Next()
	}
}

func WebMiddlewareFiber() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Get("X-User-Id")
		c.Locals("user_id", userID)
		return c.Next()
	}
}
