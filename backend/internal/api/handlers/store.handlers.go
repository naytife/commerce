package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/db"
)

// Handler struct with the repository
type Handler struct {
	Repository db.Repository
}

// NewHandler returns a handler with a repository
func NewHandler(repo db.Repository) *Handler {
	return &Handler{Repository: repo}
}

// GetCategory handler for REST API
func (h *Handler) CreateShop(c *fiber.Ctx) error {
	fakeAuthSub := "9vgPO5K5ipI424xe84HUrtqQJMWT3e7f@clients"
	owner, err := h.Repository.GetUser(c.Context(), &fakeAuthSub)
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
