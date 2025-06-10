package handlers

import (
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/services"
)

type Handler struct {
	Repository              db.Repository
	StripeService           *services.StripeService // Keep for backward compatibility
	PaymentProcessorFactory *services.PaymentProcessorFactory
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
