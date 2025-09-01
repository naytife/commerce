package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/services"
	"go.uber.org/zap"
)

// WebhookHandler handles webhook endpoints for all payment providers
type WebhookHandler struct {
	paymentProcessorFactory *services.PaymentProcessorFactory
	repository              db.Repository
}

// NewWebhookHandler creates a new webhook handler with PaymentProcessorFactory
func NewWebhookHandler(paymentProcessorFactory *services.PaymentProcessorFactory, repo db.Repository) *WebhookHandler {
	return &WebhookHandler{
		paymentProcessorFactory: paymentProcessorFactory,
		repository:              repo,
	}
}

// StripeWebhook handles Stripe webhook events
// @Summary Handle Stripe webhooks
// @Description Processes Stripe webhook events for payment updates
// @Tags webhooks
// @Accept json
// @Param shop_id path int true "Shop ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /webhooks/stripe/{shop_id} [post]
func (h *WebhookHandler) StripeWebhook(c *fiber.Ctx) error {
	return h.handleWebhook(c, "stripe")
}

// PayPalWebhook handles PayPal webhook events
// @Summary Handle PayPal webhooks
// @Description Processes PayPal webhook events for payment updates
// @Tags webhooks
// @Accept json
// @Param shop_id path int true "Shop ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /webhooks/paypal/{shop_id} [post]
func (h *WebhookHandler) PayPalWebhook(c *fiber.Ctx) error {
	return h.handleWebhook(c, "paypal")
}

// PaystackWebhook handles Paystack webhook events
// @Summary Handle Paystack webhooks
// @Description Processes Paystack webhook events for payment updates
// @Tags webhooks
// @Accept json
// @Param shop_id path int true "Shop ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /webhooks/paystack/{shop_id} [post]
func (h *WebhookHandler) PaystackWebhook(c *fiber.Ctx) error {
	return h.handleWebhook(c, "paystack")
}

// FlutterwaveWebhook handles Flutterwave webhook events
// @Summary Handle Flutterwave webhooks
// @Description Processes Flutterwave webhook events for payment updates
// @Tags webhooks
// @Accept json
// @Param shop_id path int true "Shop ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /webhooks/flutterwave/{shop_id} [post]
func (h *WebhookHandler) FlutterwaveWebhook(c *fiber.Ctx) error {
	return h.handleWebhook(c, "flutterwave")
}

// handleWebhook processes webhook events for any payment provider
func (h *WebhookHandler) handleWebhook(c *fiber.Ctx, provider string) error {
	shopID, err := strconv.ParseInt(c.Params("shop_id"), 10, 64)
	if err != nil {
		zap.L().Warn("handleWebhook: invalid shop id param", zap.String("shop_id_param", c.Params("shop_id")), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Status:  "error",
			Message: "Invalid shop ID",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Verify shop exists
	_, err = h.repository.GetShop(c.Context(), shopID)
	if err != nil {
		if err == pgx.ErrNoRows {
			zap.L().Warn("handleWebhook: shop not found", zap.Int64("shop_id", shopID))
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Status:  "error",
				Message: "Shop not found",
				Code:    fiber.StatusNotFound,
			})
		}
		zap.L().Error("handleWebhook: failed to fetch shop", zap.Int64("shop_id", shopID), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch shop",
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Get the appropriate payment processor
	processor := h.paymentProcessorFactory.GetProcessor(provider)
	if processor == nil {
		zap.L().Warn("handleWebhook: unsupported payment provider", zap.String("provider", provider))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("Unsupported payment provider: %s", provider),
			Code:    fiber.StatusBadRequest,
		})
	}

	// Get webhook payload and signature
	payload := c.Body()
	signature := c.Get("Stripe-Signature")
	if provider == "paypal" {
		signature = c.Get("PAYPAL-TRANSMISSION-SIG")
	} else if provider == "paystack" {
		signature = c.Get("X-Paystack-Signature")
	} else if provider == "flutterwave" {
		signature = c.Get("verif-hash")
	}

	// Process webhook with the appropriate processor
	webhookPayload, err := processor.HandleWebhook(c.Context(), payload, signature)
	if err != nil {
		zap.L().Error("handleWebhook: failed to process provider webhook", zap.String("provider", provider), zap.Int64("shop_id", shopID), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to process webhook: %v", err),
			Code:    fiber.StatusBadRequest,
		})
	}

	// Process the webhook payload and update payment/order status
	err = h.processWebhookPayload(c.Context(), shopID, webhookPayload)
	if err != nil {
		zap.L().Error("handleWebhook: failed to process webhook payload", zap.Int64("shop_id", shopID), zap.Error(err))
		// Still return 200 to acknowledge receipt to prevent retries
		// but log the error for investigation
	}

	zap.L().Info("handleWebhook: webhook processed",
		zap.String("provider", provider),
		zap.Int64("shop_id", shopID),
		zap.String("event", webhookPayload.EventType),
		zap.String("payment_id", webhookPayload.PaymentID))

	return c.JSON(map[string]string{
		"status": "received",
	})
}

// processWebhookPayload processes the webhook payload and updates the system state
func (h *WebhookHandler) processWebhookPayload(ctx context.Context, shopID int64, payload *models.PaymentWebhookPayload) error {
	zap.L().Info("processWebhookPayload: processing webhook",
		zap.Int64("shop_id", shopID),
		zap.String("provider", payload.Provider),
		zap.String("event", payload.EventType),
		zap.String("payment_id", payload.PaymentID),
		zap.String("transaction_id", payload.TransactionID),
		zap.String("status", payload.Status))

	// Convert payment status to database enum
	var paymentStatus db.PaymentStatusType
	switch payload.Status {
	case "paid", "completed", "successful", "success":
		paymentStatus = db.PaymentStatusTypePaid
	case "failed", "failure":
		paymentStatus = db.PaymentStatusTypeFailed
	case "refunded", "refund":
		paymentStatus = db.PaymentStatusTypeRefunded
	case "partial_refund", "partially_refunded":
		paymentStatus = db.PaymentStatusTypePartialRefund
	case "pending", "processing":
		paymentStatus = db.PaymentStatusTypePending
	default:
		zap.L().Warn("processWebhookPayload: unknown payment status, treating as pending", zap.String("status", payload.Status))
		paymentStatus = db.PaymentStatusTypePending
	}

	// Try to find and update order by transaction ID first
	if payload.TransactionID != "" {
		// Convert string to *string for database parameters
		transactionID := &payload.TransactionID

		order, err := h.repository.GetOrderByTransactionID(ctx, db.GetOrderByTransactionIDParams{
			TransactionID: transactionID,
			ShopID:        shopID,
		})

		if err == nil {
			// Order found by transaction ID, update it
			zap.L().Info("processWebhookPayload: found order by transaction id", zap.Int64("order_id", order.OrderID), zap.String("transaction_id", payload.TransactionID))

			// Determine order status based on payment status
			var orderStatus db.OrderStatusType
			switch paymentStatus {
			case db.PaymentStatusTypePaid:
				orderStatus = db.OrderStatusTypeProcessing // Payment confirmed, ready for fulfillment
			case db.PaymentStatusTypeFailed:
				orderStatus = db.OrderStatusTypeCancelled // Payment failed, cancel order
			case db.PaymentStatusTypeRefunded:
				orderStatus = db.OrderStatusTypeRefunded // Full refund issued
			default:
				orderStatus = order.Status // Keep current status for other payment states
			}

			// Update both order and payment status
			updatedOrder, err := h.repository.UpdateOrderStatusByTransactionID(ctx, db.UpdateOrderStatusByTransactionIDParams{
				TransactionID: transactionID,
				ShopID:        shopID,
				Status:        orderStatus,
				PaymentStatus: paymentStatus,
			})

			if err != nil {
				zap.L().Error("processWebhookPayload: failed to update order status by transaction id", zap.String("transaction_id", payload.TransactionID), zap.Error(err))
				return fmt.Errorf("failed to update order status: %w", err)
			}

			zap.L().Info("processWebhookPayload: successfully updated order status",
				zap.Int64("order_id", updatedOrder.OrderID),
				zap.String("status", string(updatedOrder.Status)),
				zap.String("payment_status", string(updatedOrder.PaymentStatus)))

			return nil
		} else {
			zap.L().Warn("processWebhookPayload: no order found with transaction id", zap.String("transaction_id", payload.TransactionID), zap.Error(err))
		}
	}

	// If we have a PaymentID (like Stripe payment intent ID), try to find order by metadata
	// This is for cases where the transaction ID might not be set yet or differs from PaymentID
	if payload.PaymentID != "" {
		zap.L().Info("processWebhookPayload: attempting to find order by payment id in metadata", zap.String("payment_id", payload.PaymentID))

		// In a more sophisticated implementation, you might:
		// 1. Store payment intent IDs in order metadata
		// 2. Have a separate payment_intents table linking to orders
		// 3. Use a mapping table for external payment IDs

		// For now, log the webhook data for manual investigation
		zap.L().Warn("processWebhookPayload: webhook received for payment id but no order found", zap.String("payment_id", payload.PaymentID))
		zap.L().Info("processWebhookPayload: payment details", zap.Float64("amount", payload.Amount), zap.String("currency", payload.Currency), zap.String("provider", payload.Provider))

		// Store webhook data for later processing or investigation
		// In a production system, you might want to store unmatched webhooks in a separate table
		// for manual review or retry logic
	}

	// If we reach here, we couldn't find a matching order
	// This might happen for:
	// 1. Test payments that don't correspond to real orders
	// 2. Webhooks arriving before the order is created (rare but possible)
	// 3. Payment intents created but not yet associated with an order
	zap.L().Info("processWebhookPayload: webhook processed but no matching order found",
		zap.String("payment_id", payload.PaymentID), zap.String("transaction_id", payload.TransactionID))

	return nil
}
