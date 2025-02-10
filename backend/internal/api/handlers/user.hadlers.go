package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/petrejonn/naytife/internal/db"
)

// UpsertUser creates or updates a user
// @Summary      Create or update a user
// @Description
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body db.UpsertUserParams true "User object that needs to be created or updated"
// @Success      200  {object}   models.ResponseHTTP{data=models.UserResponse} "User created or updated successfully"
// @Router       /auth/register [post]
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

// GetMe fetches the currently authenticated user
// @Summary      Fetch the currently authenticated user
// @Description
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}   models.ResponseHTTP{data=models.UserResponse} "User fetched successfully"
// @Security     OAuth2AccessCode
// @Router       /me [get]
func (h *Handler) GetMe(c *fiber.Ctx) error {
	userIDStr, _ := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	objDB, err := h.Repository.GetUserById(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.JSON(objDB)
}
