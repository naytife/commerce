package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/db"
)

func (h *Handler) CreateShop(c *fiber.Ctx) error {
	fakeAuthSub := "user@example.com"
	owner, err := h.Repository.GetUser(c.Context(), fakeAuthSub)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}
	var param db.CreateShopParams
	c.BodyParser(&param)
	param.OwnerID = owner.UserID
	objDB, err := h.Repository.CreateShop(c.Context(), param)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create shop",
		})
	}

	return c.JSON(objDB)
}
