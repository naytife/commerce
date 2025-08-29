package handlers

import (
	"github.com/go-redis/redis/v8"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/services"
	"go.uber.org/zap"
)

type Handler struct {
	Repository              db.Repository
	PaymentProcessorFactory *services.PaymentProcessorFactory
	RetryClient             *retryablehttp.Client
	Logger                  *zap.Logger
}

func NewHandler(repo db.Repository, logger *zap.Logger, retryClient *retryablehttp.Client) *Handler {
	return &Handler{
		Repository:  repo,
		Logger:      logger,
		RetryClient: retryClient,
	}
}

// func NewHandlerWithServices(repo db.Repository, stripeService *services.StripeService) *Handler {
// 	return &Handler{
// 		Repository:    repo,
// 		StripeService: stripeService,
// 	}
// }

func NewHandlerWithPaymentFactory(repo db.Repository, logger *zap.Logger, retryClient *retryablehttp.Client, paymentFactory *services.PaymentProcessorFactory) *Handler {
	return &Handler{
		Repository:              repo,
		Logger:                  logger,
		RetryClient:             retryClient,
		PaymentProcessorFactory: paymentFactory,
	}
}

func NewHandlerWithRedis(repo db.Repository, redisClient *redis.Client) *Handler {
	return &Handler{
		Repository: repo,
		// PublishHandler: NewPublishHandler(repo, redisClient),
	}
}

// Publish methods that delegate to the PublishHandler
// func (h *Handler) TriggerBuild(c *fiber.Ctx) error {
// 	if h.PublishHandler == nil {
// 		return fiber.NewError(fiber.StatusServiceUnavailable, "Publish service not available")
// 	}
// 	return h.PublishHandler.TriggerSiteBuild(c)
// }

// func (h *Handler) GetBuildStatus(c *fiber.Ctx) error {
// 	if h.PublishHandler == nil {
// 		return fiber.NewError(fiber.StatusServiceUnavailable, "Publish service not available")
// 	}
// 	return h.PublishHandler.GetPublishStatus(c)
// }

// func (h *Handler) GetBuildHistory(c *fiber.Ctx) error {
// 	if h.PublishHandler == nil {
// 		return fiber.NewError(fiber.StatusServiceUnavailable, "Publish service not available")
// 	}
// 	return h.PublishHandler.GetPublishHistory(c)
// }
