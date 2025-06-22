package handlers

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/services"
)

type Handler struct {
	Repository              db.Repository
	StripeService           *services.StripeService // Keep for backward compatibility
	PaymentProcessorFactory *services.PaymentProcessorFactory
	PublishHandler          *PublishHandler // Add publish handler
}

func NewHandler(repo db.Repository) *Handler {
	return &Handler{
		Repository: repo,
	}
}

func NewHandlerWithServices(repo db.Repository, stripeService *services.StripeService) *Handler {
	return &Handler{
		Repository:    repo,
		StripeService: stripeService,
	}
}

func NewHandlerWithPaymentFactory(repo db.Repository, paymentFactory *services.PaymentProcessorFactory) *Handler {
	return &Handler{
		Repository:              repo,
		PaymentProcessorFactory: paymentFactory,
	}
}

func NewHandlerWithRedis(repo db.Repository, redisClient *redis.Client) *Handler {
	return &Handler{
		Repository:     repo,
		PublishHandler: NewPublishHandler(repo, redisClient),
	}
}

// Publish methods that delegate to the PublishHandler
func (h *Handler) TriggerBuild(c *fiber.Ctx) error {
	if h.PublishHandler == nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, "Publish service not available")
	}
	return h.PublishHandler.TriggerSiteBuild(c)
}

func (h *Handler) GetBuildStatus(c *fiber.Ctx) error {
	if h.PublishHandler == nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, "Publish service not available")
	}
	return h.PublishHandler.GetPublishStatus(c)
}

func (h *Handler) GetBuildHistory(c *fiber.Ctx) error {
	if h.PublishHandler == nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, "Publish service not available")
	}
	return h.PublishHandler.GetPublishHistory(c)
}
