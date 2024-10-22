package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/db"
)

// Handler struct with the repository
// GetCategory handler for REST API
func (h *Handler) UpsertUser(c *fiber.Ctx) error {
	var param db.UpsertUserParams
	c.BodyParser(&param)
	objDB, err := h.Repository.UpsertUser(c.Context(), param)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create user",
		})
	}

	return c.JSON(objDB)
}
