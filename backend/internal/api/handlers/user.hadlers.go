package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
)

// UpsertUser creates or updates a user
// @Summary      Create or update a user
// @Description
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body models.RegisterUserParams true "User object that needs to be created or updated"
// @Success      200  {object}   models.SuccessResponse{data=models.UserResponse} "User created or updated successfully"
// @Router       /auth/register [post]
func (h *Handler) UpsertUser(c *fiber.Ctx) error {
	var param models.RegisterUserParams
	err := c.BodyParser(&param)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	validator := &models.XValidator{}
	if errs := validator.Validate(&param); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)

		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}
	var userParam db.UpsertUserParams
	copier.Copy(&userParam, &param)
	userParam.Sub = param.Email
	objDB, err := h.Repository.UpsertUser(c.Context(), userParam)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}
	var resp models.UserResponse
	copier.Copy(&resp, &objDB)
	return c.JSON(resp)
}

// GetMe fetches the currently authenticated user
// @Summary      Fetch the currently authenticated user
// @Description
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}   models.SuccessResponse{data=models.UserResponse} "User fetched successfully"
// @Security     OAuth2AccessCode
// @Router       /me [get]
func (h *Handler) GetMe(c *fiber.Ctx) error {
	userIDStr, _ := c.Locals("user_id").(string)
	objDB, err := h.Repository.GetUserBySub(c.Context(), &userIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}
	var resp models.UserResponse
	copier.Copy(&resp, &objDB)
	return c.JSON(resp)
}
