package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/db"
)

func (h *Handler) UpsertUser(c *fiber.Ctx) error {
	var param db.UpsertUserParams
	err := c.BodyParser(&param)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	objDB, err := h.Repository.UpsertUser(c.Context(), param)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.JSON(objDB)
}
