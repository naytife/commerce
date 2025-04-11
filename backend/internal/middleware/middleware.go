package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/db"
)

func ShopIDMiddlewareFiber(repo db.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		subdomain := strings.Split(c.Hostname(), ".")[0]
		ctx := c.UserContext()

		shopID, err := repo.GetShopIDBySubDomain(ctx, subdomain)
		if err != nil {
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
		fmt.Println("User ID from header:", userID)
		c.Locals("user_id", userID)
		return c.Next()
	}
}

func GlobalErrorHandler(c *fiber.Ctx) error {
	// Proceed to the next middleware/handler
	err := c.Next()

	// Check if an error occurred
	if err != nil {
		// Default error response
		statusCode := fiber.StatusInternalServerError
		message := "An unexpected error occurred"

		// Handle specific error types
		if e, ok := err.(*fiber.Error); ok {
			statusCode = e.Code
			switch statusCode {
			case fiber.StatusBadRequest:
				message = "Invalid input data"
			case fiber.StatusNotFound:
				message = "Resource not found"
			case fiber.StatusUnauthorized:
				message = "Authentication required"
			case fiber.StatusForbidden:
				message = "Insufficient permissions"
			}
		}

		// Send the error response
		return api.ErrorResponse(c, statusCode, message, nil)
	}

	return nil
}
