package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"go.uber.org/zap"
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
	if err := c.BodyParser(&param); err != nil {
		zap.L().Warn("UpsertUser: failed to parse request body", zap.Error(err))
		return api.BusinessLogicErrorResponse(c, "Failed to parse request body")
	}

	if err := api.ValidateRequest(c, &param); err != nil {
		return err
	}

	objDB, err := h.Repository.UpsertUser(c.Context(), db.UpsertUserParams{
		Sub:            param.Email,
		AuthProviderID: param.ProviderID,
		AuthProvider:   param.Provider,
		Email:          param.Email,
		Name:           param.Name,
		Locale:         param.Locale,
		ProfilePicture: param.ProfilePicture,
		VerifiedEmail:  param.VerifiedEmail,
	})
	if err != nil {
		userSub := ""
		if param.Email != nil {
			userSub = *param.Email
		}
		zap.L().Error("UpsertUser: failed to upsert user", zap.Error(err), zap.String("user_sub", userSub))
		return api.SystemErrorResponse(c, err, "Failed to create or update user")
	}
	resp := models.UserResponse{
		UserID:         objDB.UserID,
		Provider:       objDB.AuthProvider,
		Email:          objDB.Email,
		Name:           objDB.Name,
		ProfilePicture: objDB.ProfilePicture,
		CreatedAt:      objDB.CreatedAt,
		LastLogin:      objDB.LastLogin,
		ProviderID:     objDB.AuthProviderID,
		Locale:         objDB.Locale,
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "User created or updated successfully")
}

// GetMe fetches the currently authenticated user
// @Summary      Fetch the currently authenticated user
// @Description
// @Tags         user
// @Produce      json
// @Success      200  {object}   models.SuccessResponse{data=models.UserResponse} "User fetched successfully"
// @Security     OAuth2AccessCode
// @Router       /me [get]
func (h *Handler) GetMe(c *fiber.Ctx) error {
	userIDStr, _ := c.Locals("user_id").(string)
	objDB, err := h.Repository.GetUserBySubWithShops(c.Context(), &userIDStr)
	if err != nil {
		zap.L().Error("GetMe: failed to get user by sub", zap.String("user_sub", userIDStr), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get user", nil)
	}
	var shops []models.UserShopResponse
	if err := json.Unmarshal(objDB.Shops, &shops); err != nil {
		zap.L().Error("GetMe: failed to unmarshal user shops", zap.String("user_sub", userIDStr), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get user shops", nil)
	}
	resp := models.UserResponse{
		UserID:         objDB.UserID,
		Provider:       objDB.AuthProvider,
		Email:          objDB.Email,
		Name:           objDB.Name,
		ProfilePicture: objDB.ProfilePicture,
		CreatedAt:      objDB.CreatedAt,
		LastLogin:      objDB.LastLogin,
		ProviderID:     objDB.AuthProviderID,
		Locale:         objDB.Locale,
		Shops:          shops,
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "User fetched successfully")

}

// GetUser fetches a user by email
// @Summary      Fetch a user by email
// @Description
// @Tags         user
// @Produce      json
// @Param        email query string true "User email"
// @Success      200  {object}   models.SuccessResponse{data=models.UserResponse} "User fetched successfully"
// @Router       /userinfo [get]
func (h *Handler) GetUser(c *fiber.Ctx) error {
	email := c.Query("email")
	objDB, err := h.Repository.GetUser(c.Context(), &email)
	if err != nil {
		if err == pgx.ErrNoRows {
			zap.L().Warn("GetUser: user not found", zap.String("email", email))
			return api.ErrorResponse(c, fiber.StatusNotFound, "User not found", nil)
		}
		zap.L().Error("GetUser: failed to get user", zap.String("email", email), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get user", nil)
	}
	resp := models.UserResponse{
		UserID:         objDB.UserID,
		Provider:       objDB.AuthProvider,
		Email:          objDB.Email,
		Name:           objDB.Name,
		ProfilePicture: objDB.ProfilePicture,
		CreatedAt:      objDB.CreatedAt,
		LastLogin:      objDB.LastLogin,
		ProviderID:     objDB.AuthProviderID,
		Locale:         objDB.Locale,
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "User fetched successfully")
}
