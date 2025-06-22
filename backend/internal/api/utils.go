package api

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/models"
)

// ErrorCategory represents different types of errors
type ErrorCategory string

const (
	ErrorCategoryValidation   ErrorCategory = "VALIDATION_ERROR"
	ErrorCategoryNotFound     ErrorCategory = "NOT_FOUND_ERROR"
	ErrorCategoryUnauthorized ErrorCategory = "UNAUTHORIZED_ERROR"
	ErrorCategoryBusiness     ErrorCategory = "BUSINESS_LOGIC_ERROR"
	ErrorCategorySystem       ErrorCategory = "SYSTEM_ERROR"
	ErrorCategoryExternal     ErrorCategory = "EXTERNAL_SERVICE_ERROR"
)

// LogContext provides structured logging context
type LogContext struct {
	RequestID string
	UserID    string
	ShopID    string
	Action    string
	Resource  string
}

// ExtractLogContext extracts logging context from fiber context
func ExtractLogContext(c *fiber.Ctx) LogContext {
	return LogContext{
		RequestID: c.Get("X-Request-ID", "unknown"),
		UserID:    c.Get("X-User-ID", "unknown"),
		ShopID:    c.Params("shop_id", "unknown"),
		Action:    c.Method(),
		Resource:  c.Path(),
	}
}

// LogError logs errors with structured context
func LogError(ctx context.Context, logCtx LogContext, category ErrorCategory, err error, message string) {
	log.Printf("[ERROR] Category: %s | RequestID: %s | UserID: %s | ShopID: %s | Action: %s %s | Message: %s | Error: %v",
		category, logCtx.RequestID, logCtx.UserID, logCtx.ShopID, logCtx.Action, logCtx.Resource, message, err)
}

// ValidationErrorResponse handles validation errors consistently
func ValidationErrorResponse(c *fiber.Ctx, validationErrors []models.Error) error {
	logCtx := ExtractLogContext(c)
	LogError(c.Context(), logCtx, ErrorCategoryValidation, nil, "Validation failed")

	return ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
}

// NotFoundErrorResponse handles not found errors
func NotFoundErrorResponse(c *fiber.Ctx, resource string) error {
	logCtx := ExtractLogContext(c)
	LogError(c.Context(), logCtx, ErrorCategoryNotFound, nil, fmt.Sprintf("%s not found", resource))

	return ErrorResponse(c, fiber.StatusNotFound, fmt.Sprintf("%s not found", resource), nil)
}

// UnauthorizedErrorResponse handles unauthorized errors
func UnauthorizedErrorResponse(c *fiber.Ctx, message string) error {
	logCtx := ExtractLogContext(c)
	LogError(c.Context(), logCtx, ErrorCategoryUnauthorized, nil, message)

	if message == "" {
		message = "Authentication required"
	}
	return ErrorResponse(c, fiber.StatusUnauthorized, message, nil)
}

// BusinessLogicErrorResponse handles business logic errors
func BusinessLogicErrorResponse(c *fiber.Ctx, message string) error {
	logCtx := ExtractLogContext(c)
	LogError(c.Context(), logCtx, ErrorCategoryBusiness, nil, message)

	return ErrorResponse(c, fiber.StatusBadRequest, message, nil)
}

// SystemErrorResponse handles system/internal errors
func SystemErrorResponse(c *fiber.Ctx, err error, message string) error {
	logCtx := ExtractLogContext(c)
	LogError(c.Context(), logCtx, ErrorCategorySystem, err, message)

	// Don't expose internal error details to clients
	return ErrorResponse(c, fiber.StatusInternalServerError, message, nil)
}

// ExternalServiceErrorResponse handles external service errors
func ExternalServiceErrorResponse(c *fiber.Ctx, service string, err error) error {
	logCtx := ExtractLogContext(c)
	message := fmt.Sprintf("%s service temporarily unavailable", service)
	LogError(c.Context(), logCtx, ErrorCategoryExternal, err, message)

	return ErrorResponse(c, fiber.StatusServiceUnavailable, message, nil)
}

// ParseIDParameter safely parses ID parameters with proper error handling
func ParseIDParameter(c *fiber.Ctx, paramName string, resourceName string) (int64, error) {
	idStr := c.Params(paramName)
	if idStr == "" {
		return 0, BusinessLogicErrorResponse(c, fmt.Sprintf("%s ID is required", resourceName))
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, BusinessLogicErrorResponse(c, fmt.Sprintf("Invalid %s ID", resourceName))
	}

	return id, nil
}

// ValidateRequest validates request body and returns proper error response
func ValidateRequest(c *fiber.Ctx, data interface{}) error {
	validator := &models.XValidator{}
	if errs := validator.Validate(data); len(errs) > 0 {
		return ValidationErrorResponse(c, errs)
	}
	return nil
}

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

// Response optimization helpers for Phase 3B

// PaginatedResponse represents a standardized paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
	HasNext    bool        `json:"has_next"`
	HasPrev    bool        `json:"has_prev"`
}

// PaginatedSuccessResponse creates a standardized paginated success response
func PaginatedSuccessResponse(c *fiber.Ctx, statusCode int, data interface{}, total int64, page, limit int, message string) error {
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := models.SuccessResponse{
		Status: "success",
		Code:   statusCode,
		Data: PaginatedResponse{
			Data:       data,
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
		Message: message,
		Meta: models.Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
	return c.Status(statusCode).JSON(response)
}

// ParsePaginationParams safely parses pagination parameters from query
func ParsePaginationParams(c *fiber.Ctx) (limit, offset int, err error) {
	limitStr := c.Query("limit", "20")
	offsetStr := c.Query("offset", "0")

	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	offset, err = strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	return limit, offset, nil
}

// ResponseMapper helps standardize response transformations
type ResponseMapper struct {
	transformers map[string]func(interface{}) interface{}
}

// NewResponseMapper creates a new response mapper
func NewResponseMapper() *ResponseMapper {
	return &ResponseMapper{
		transformers: make(map[string]func(interface{}) interface{}),
	}
}

// Register adds a transformer for a specific type
func (rm *ResponseMapper) Register(typeName string, transformer func(interface{}) interface{}) {
	rm.transformers[typeName] = transformer
}

// Transform applies the registered transformer to data
func (rm *ResponseMapper) Transform(typeName string, data interface{}) interface{} {
	if transformer, exists := rm.transformers[typeName]; exists {
		return transformer(data)
	}
	return data
}

// Common response transformers
var CommonResponseMapper = NewResponseMapper()

func init() {
	// Register common transformations
	CommonResponseMapper.Register("customer", func(data interface{}) interface{} {
		// Add any common customer transformations here
		return data
	})

	CommonResponseMapper.Register("order", func(data interface{}) interface{} {
		// Add any common order transformations here
		return data
	})

	CommonResponseMapper.Register("product", func(data interface{}) interface{} {
		// Add any common product transformations here
		return data
	})
}

// CacheableSuccessResponse creates a success response with cache headers
func CacheableSuccessResponse(c *fiber.Ctx, statusCode int, data interface{}, message string, maxAge int) error {
	c.Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
	c.Set("ETag", generateETag(data))

	return SuccessResponse(c, statusCode, data, message)
}

// generateETag creates a simple ETag for response caching
func generateETag(data interface{}) string {
	hash := fmt.Sprintf("%x", data)
	if len(hash) > 16 {
		hash = hash[:16]
	}
	return fmt.Sprintf(`"%s"`, hash)
}

// BatchResponse handles batch operation responses
type BatchResponse struct {
	Success []interface{} `json:"success"`
	Failed  []BatchError  `json:"failed"`
	Total   int           `json:"total"`
}

type BatchError struct {
	Index   int    `json:"index"`
	ID      string `json:"id,omitempty"`
	Message string `json:"message"`
}

// BatchSuccessResponse creates a standardized batch operation response
func BatchSuccessResponse(c *fiber.Ctx, statusCode int, success []interface{}, failed []BatchError, message string) error {
	response := models.SuccessResponse{
		Status: "success",
		Code:   statusCode,
		Data: BatchResponse{
			Success: success,
			Failed:  failed,
			Total:   len(success) + len(failed),
		},
		Message: message,
		Meta: models.Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
	return c.Status(statusCode).JSON(response)
}
