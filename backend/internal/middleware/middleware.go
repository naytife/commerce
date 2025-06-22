package middleware

import (
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/db"
)

func ShopIDMiddlewareFiber(repo db.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var subdomain string
		ctx := c.UserContext()

		// Check for cloud-build specific header first
		if cloudBuildSubdomain := c.Get("X-Shop-Subdomain"); cloudBuildSubdomain != "" {
			subdomain = cloudBuildSubdomain
		} else if shopIDParam := c.Query("shop_id"); shopIDParam != "" {
			// Handle direct shop_id parameter (for cloud-build)
			shopID := parseInt64(shopIDParam)
			if shopID > 0 {
				if err := repo.SetShopIDInSession(ctx, shopID); err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Failed to set shop_id",
					})
				}
				c.Locals("shop_id", shopID)
				return c.Next()
			}
		} else {
			// Extract subdomain from hostname (original behavior)
			subdomain = strings.Split(c.Hostname(), ".")[0]
		}

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

// Helper function to parse int64 safely
func parseInt64(s string) int64 {
	if val, err := strconv.ParseInt(s, 10, 64); err == nil {
		return val
	}
	return 0
}

func WebMiddlewareFiber() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Get("X-User-Id")

		// Development mode bypass - allow testing without full OAuth stack
		if userID == "" {
			// Check if we're in development mode
			devMode := os.Getenv("DEV_MODE")
			if devMode == "true" {
				// Use a test user email for development (sub field in users table)
				userID = "test@example.com"
				c.Set("X-User-Id", userID)
			}
		}

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
