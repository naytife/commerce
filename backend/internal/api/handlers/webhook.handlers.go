package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/services"
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
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Status:  "error",
			Message: "Invalid shop ID",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Verify shop exists
	_, err = h.repository.GetShop(c.Context(), shopID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Status:  "error",
			Message: "Shop not found",
			Code:    fiber.StatusNotFound,
		})
	}

	// Get the appropriate payment processor
	processor := h.paymentProcessorFactory.GetProcessor(provider)
	if processor == nil {
		log.Printf("Unsupported payment provider: %s", provider)
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
		log.Printf("Failed to process %s webhook for shop %d: %v", provider, shopID, err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to process webhook: %v", err),
			Code:    fiber.StatusBadRequest,
		})
	}

	// Process the webhook payload and update payment/order status
	err = h.processWebhookPayload(c.Context(), shopID, webhookPayload)
	if err != nil {
		log.Printf("Failed to process webhook payload for shop %d: %v", shopID, err)
		// Still return 200 to acknowledge receipt to prevent retries
		// but log the error for investigation
	}

	log.Printf("Successfully processed %s webhook for shop %d, event: %s, payment: %s",
		provider, shopID, webhookPayload.EventType, webhookPayload.PaymentID)

	return c.JSON(map[string]string{
		"status": "received",
	})
}

// processWebhookPayload processes the webhook payload and updates the system state
func (h *WebhookHandler) processWebhookPayload(ctx context.Context, shopID int64, payload *models.PaymentWebhookPayload) error {
	log.Printf("Processing webhook for shop %d: provider=%s, event=%s, payment_id=%s, transaction_id=%s, status=%s",
		shopID, payload.Provider, payload.EventType, payload.PaymentID, payload.TransactionID, payload.Status)

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
		log.Printf("Unknown payment status '%s', treating as pending", payload.Status)
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
			log.Printf("Found order %d by transaction ID %s", order.OrderID, payload.TransactionID)

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
				log.Printf("Failed to update order status by transaction ID %s: %v", payload.TransactionID, err)
				return fmt.Errorf("failed to update order status: %w", err)
			}

			log.Printf("Successfully updated order %d: status=%s, payment_status=%s",
				updatedOrder.OrderID, updatedOrder.Status, updatedOrder.PaymentStatus)

			return nil
		} else {
			log.Printf("No order found with transaction ID %s: %v", payload.TransactionID, err)
		}
	}

	// If we have a PaymentID (like Stripe payment intent ID), try to find order by metadata
	// This is for cases where the transaction ID might not be set yet or differs from PaymentID
	if payload.PaymentID != "" {
		log.Printf("Attempting to find order by payment ID %s in metadata", payload.PaymentID)

		// In a more sophisticated implementation, you might:
		// 1. Store payment intent IDs in order metadata
		// 2. Have a separate payment_intents table linking to orders
		// 3. Use a mapping table for external payment IDs

		// For now, log the webhook data for manual investigation
		log.Printf("Webhook received for payment ID %s but no order found by transaction ID", payload.PaymentID)
		log.Printf("Payment details: amount=%.2f %s, provider=%s", payload.Amount, payload.Currency, payload.Provider)

		// Store webhook data for later processing or investigation
		// In a production system, you might want to store unmatched webhooks in a separate table
		// for manual review or retry logic
	}

	// If we reach here, we couldn't find a matching order
	// This might happen for:
	// 1. Test payments that don't correspond to real orders
	// 2. Webhooks arriving before the order is created (rare but possible)
	// 3. Payment intents created but not yet associated with an order
	log.Printf("Webhook processed but no matching order found for payment %s / transaction %s",
		payload.PaymentID, payload.TransactionID)

	return nil
}
