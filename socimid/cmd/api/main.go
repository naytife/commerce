package main

import (
	"github.com/gofiber/fiber/v2"
	fb "github.com/huandu/facebook/v2"
)

// @title Naytife API Docs
// @version 1.0
// @description This is the Naytife API documentation
// @host localhost:8000
// @BasePath /api/v1
// @schemes http
// securityDefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://auth.naytife.com/oauth2/token
// @authorizationUrl https://auth.naytife.com/oauth2/auth
// @securityDefinitions.apikey XUserIdAuth
// @in header
// @name X-User-Id
func main() {
	app := fiber.New()

	app.Post("/upload-photo", func(c *fiber.Ctx) error {
		// Assume file is uploaded and get the file URL or data
		fileURL := "https://loremflickr.com/320/240" // Replace with actual file handling
		accessToken := "your-access-token"
		pageID := "your-page-id"

		// Initialize Facebook session
		session := fb.New("", "").Session(accessToken)

		// Upload photo (first step)
		params := fb.Params{
			"url":       fileURL,
			"published": "false",
		}
		res, err := session.Post(pageID+"/photos", params)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		photoID := res["id"].(string)
		// Publish as story (second step)
		storyParams := fb.Params{
			"photo_id": photoID,
		}
		_, err = session.Post(pageID+"/photo_stories", storyParams)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"success": true, "message": "Photo posted to story"})
	})

	app.Listen(":8003")
}
