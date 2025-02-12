package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/models"
)

func ErrorResponse(c *fiber.Ctx, statusCode int, message string, errors []models.Error) error {
	response := models.ErrorResponse{
		Status:  "error",
		Code:    statusCode,
		Message: message,
		Errors:  errors,
		Meta: models.Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
	return c.Status(statusCode).JSON(response)
}

func SuccessResponse(c *fiber.Ctx, statusCode int, data interface{}, message string) error {
	response := models.SuccessResponse{
		Status:  "success",
		Code:    statusCode,
		Data:    data,
		Message: message,
		Meta: models.Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
	return c.Status(statusCode).JSON(response)
}
