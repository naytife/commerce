package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api"
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
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse request body", nil)
	}
	validator := &models.XValidator{}
	if errs := validator.Validate(&param); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)

		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}
	objDB, err := h.Repository.UpsertUser(c.Context(), db.UpsertUserParams{
		Sub:            param.Email,
		ProviderID:     param.ProviderID,
		Provider:       param.Provider,
		Email:          param.Email,
		Name:           param.Name,
		Locale:         param.Locale,
		ProfilePicture: param.ProfilePicture,
		VerifiedEmail:  param.VerifiedEmail,
	})
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create or update user", nil)
	}
	resp := models.UserResponse{
		UserID:         objDB.UserID,
		Provider:       objDB.Provider,
		Email:          objDB.Email,
		Name:           objDB.Name,
		ProfilePicture: objDB.ProfilePicture,
		CreatedAt:      objDB.CreatedAt,
		LastLogin:      objDB.LastLogin,
		ProviderID:     objDB.ProviderID,
		Locale:         objDB.Locale,
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "User created or updated successfully")
}

// GetMe fetches the currently authenticated user
// @Summary      Fetch the currently authenticated user
// @Description
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}   models.SuccessResponse{data=models.UserResponse} "User fetched successfully"
// @Security     OAuth2AccessCode
// @Security XUserIdAuth
// @Router       /me [get]
func (h *Handler) GetMe(c *fiber.Ctx) error {
	userIDStr, _ := c.Locals("user_id").(string)
	objDB, err := h.Repository.GetUserBySubWithShops(c.Context(), &userIDStr)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get user", nil)
	}
	var shops []models.UserShopResponse
	if err := json.Unmarshal(objDB.Shops, &shops); err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get user shops", nil)
	}
	resp := models.UserResponse{
		UserID:         objDB.UserID,
		Provider:       objDB.Provider,
		Email:          objDB.Email,
		Name:           objDB.Name,
		ProfilePicture: objDB.ProfilePicture,
		CreatedAt:      objDB.CreatedAt,
		LastLogin:      objDB.LastLogin,
		ProviderID:     objDB.ProviderID,
		Locale:         objDB.Locale,
		Shops:          shops,
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "User fetched successfully")

}
