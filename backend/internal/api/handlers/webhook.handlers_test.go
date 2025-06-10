package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/mocks"
	"github.com/petrejonn/naytife/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockPaymentProcessor is a mock implementation of PaymentProcessor interface
type MockPaymentProcessor struct {
	handleWebhookFunc func(ctx context.Context, payload []byte, signature string) (*models.PaymentWebhookPayload, error)
}

func (m *MockPaymentProcessor) ProcessPayment(ctx context.Context, shopID int64, req models.PaymentRequest, amount float64, currencyCode string) (*models.PaymentResponse, error) {
	return nil, nil
}

func (m *MockPaymentProcessor) CreatePaymentIntent(ctx context.Context, shopID int64, req models.PaymentIntentRequest) (*models.PaymentIntentResponse, error) {
	return nil, nil
}

func (m *MockPaymentProcessor) ConfirmPayment(ctx context.Context, shopID int64, paymentID string) (*models.PaymentResponse, error) {
	return nil, nil
}

func (m *MockPaymentProcessor) GetPaymentStatus(ctx context.Context, shopID int64, paymentID string) (*models.PaymentResponse, error) {
	return nil, nil
}

func (m *MockPaymentProcessor) RefundPayment(ctx context.Context, shopID int64, paymentID string, amount float64, reason string) (*models.PaymentResponse, error) {
	return nil, nil
}

func (m *MockPaymentProcessor) HandleWebhook(ctx context.Context, payload []byte, signature string) (*models.PaymentWebhookPayload, error) {
	if m.handleWebhookFunc != nil {
		return m.handleWebhookFunc(ctx, payload, signature)
	}
	return nil, nil
}

func (m *MockPaymentProcessor) ValidateConfig(config map[string]interface{}) error {
	return nil
}

// Create a mock factory that implements the same interface as the real factory
func createMockPaymentProcessorFactory(processors map[string]services.PaymentProcessor) *services.PaymentProcessorFactory {
	stripe := processors["stripe"]
	paypal := processors["paypal"]
	paystack := processors["paystack"]
	flutterwave := processors["flutterwave"]

	return services.NewPaymentProcessorFactory(stripe, paypal, paystack, flutterwave)
}

func TestWebhookHandler_StripeWebhook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockProcessor := &MockPaymentProcessor{}

	processors := map[string]services.PaymentProcessor{
		"stripe": mockProcessor,
	}
	mockProcessorFactory := createMockPaymentProcessorFactory(processors)

	handler := NewWebhookHandler(mockProcessorFactory, mockRepo)

	app := fiber.New()
	app.Post("/webhooks/stripe/:shop_id", handler.StripeWebhook)

	t.Run("successful webhook processing", func(t *testing.T) {
		shopID := int64(123)
		transactionID := "txn_123456"

		// Setup mock expectations
		mockRepo.EXPECT().
			GetShop(gomock.Any(), shopID).
			Return(db.Shop{ShopID: shopID}, nil)

		mockRepo.EXPECT().
			GetOrderByTransactionID(gomock.Any(), db.GetOrderByTransactionIDParams{
				TransactionID: &transactionID,
				ShopID:        shopID,
			}).
			Return(db.Order{
				OrderID: 1,
				Status:  db.OrderStatusTypePending,
			}, nil)

		mockRepo.EXPECT().
			UpdateOrderStatusByTransactionID(gomock.Any(), db.UpdateOrderStatusByTransactionIDParams{
				TransactionID: &transactionID,
				ShopID:        shopID,
				Status:        db.OrderStatusTypeProcessing,
				PaymentStatus: db.PaymentStatusTypePaid,
			}).
			Return(db.Order{
				OrderID:       1,
				Status:        db.OrderStatusTypeProcessing,
				PaymentStatus: db.PaymentStatusTypePaid,
			}, nil)

		// Setup webhook payload mock
		webhookPayload := &models.PaymentWebhookPayload{
			Provider:      "stripe",
			EventType:     "payment_intent.succeeded",
			PaymentID:     "pi_123456",
			TransactionID: transactionID,
			Status:        "paid",
			Amount:        100.00,
			Currency:      "USD",
			ReceivedAt:    time.Now(),
		}

		mockProcessor.handleWebhookFunc = func(ctx context.Context, payload []byte, signature string) (*models.PaymentWebhookPayload, error) {
			return webhookPayload, nil
		}

		// Make request
		reqBody := `{"type": "payment_intent.succeeded", "data": {"object": {"id": "pi_123456"}}}`
		req := httptest.NewRequest("POST", "/webhooks/stripe/123", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Stripe-Signature", "t=1234567890,v1=test_signature")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Verify response
		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "received", response["status"])
	})

	t.Run("invalid shop ID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/webhooks/stripe/invalid", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)

		var errorResp models.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid shop ID", errorResp.Message)
	})

	t.Run("shop not found", func(t *testing.T) {
		shopID := int64(999)

		mockRepo.EXPECT().
			GetShop(gomock.Any(), shopID).
			Return(db.Shop{}, errors.New("shop not found"))

		req := httptest.NewRequest("POST", "/webhooks/stripe/999", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)

		var errorResp models.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		require.NoError(t, err)
		assert.Equal(t, "Shop not found", errorResp.Message)
	})

	t.Run("webhook signature validation fails", func(t *testing.T) {
		shopID := int64(123)

		mockRepo.EXPECT().
			GetShop(gomock.Any(), shopID).
			Return(db.Shop{ShopID: shopID}, nil)

		mockProcessor.handleWebhookFunc = func(ctx context.Context, payload []byte, signature string) (*models.PaymentWebhookPayload, error) {
			return nil, errors.New("invalid signature")
		}

		req := httptest.NewRequest("POST", "/webhooks/stripe/123", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Stripe-Signature", "invalid_signature")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)

		var errorResp models.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		require.NoError(t, err)
		assert.Contains(t, errorResp.Message, "Failed to process webhook")
	})
}

func TestWebhookHandler_PayPalWebhook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockProcessor := &MockPaymentProcessor{}

	processors := map[string]services.PaymentProcessor{
		"paypal": mockProcessor,
	}
	mockProcessorFactory := createMockPaymentProcessorFactory(processors)

	handler := NewWebhookHandler(mockProcessorFactory, mockRepo)

	app := fiber.New()
	app.Post("/webhooks/paypal/:shop_id", handler.PayPalWebhook)

	t.Run("successful PayPal webhook processing", func(t *testing.T) {
		shopID := int64(123)
		transactionID := "PAY-123456"

		mockRepo.EXPECT().
			GetShop(gomock.Any(), shopID).
			Return(db.Shop{ShopID: shopID}, nil)

		mockRepo.EXPECT().
			GetOrderByTransactionID(gomock.Any(), db.GetOrderByTransactionIDParams{
				TransactionID: &transactionID,
				ShopID:        shopID,
			}).
			Return(db.Order{
				OrderID: 1,
				Status:  db.OrderStatusTypePending,
			}, nil)

		mockRepo.EXPECT().
			UpdateOrderStatusByTransactionID(gomock.Any(), db.UpdateOrderStatusByTransactionIDParams{
				TransactionID: &transactionID,
				ShopID:        shopID,
				Status:        db.OrderStatusTypeProcessing,
				PaymentStatus: db.PaymentStatusTypePaid,
			}).
			Return(db.Order{
				OrderID:       1,
				Status:        db.OrderStatusTypeProcessing,
				PaymentStatus: db.PaymentStatusTypePaid,
			}, nil)

		webhookPayload := &models.PaymentWebhookPayload{
			Provider:      "paypal",
			EventType:     "PAYMENT.CAPTURE.COMPLETED",
			PaymentID:     "PAYID-123456",
			TransactionID: transactionID,
			Status:        "completed",
			Amount:        50.00,
			Currency:      "USD",
			ReceivedAt:    time.Now(),
		}

		mockProcessor.handleWebhookFunc = func(ctx context.Context, payload []byte, signature string) (*models.PaymentWebhookPayload, error) {
			return webhookPayload, nil
		}

		reqBody := `{"event_type": "PAYMENT.CAPTURE.COMPLETED"}`
		req := httptest.NewRequest("POST", "/webhooks/paypal/123", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("PAYPAL-TRANSMISSION-SIG", "test_signature")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestWebhookHandler_PaystackWebhook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockProcessor := &MockPaymentProcessor{}

	processors := map[string]services.PaymentProcessor{
		"paystack": mockProcessor,
	}
	mockProcessorFactory := createMockPaymentProcessorFactory(processors)

	handler := NewWebhookHandler(mockProcessorFactory, mockRepo)

	app := fiber.New()
	app.Post("/webhooks/paystack/:shop_id", handler.PaystackWebhook)

	t.Run("successful Paystack webhook processing", func(t *testing.T) {
		shopID := int64(123)
		transactionID := "ref_123456"

		mockRepo.EXPECT().
			GetShop(gomock.Any(), shopID).
			Return(db.Shop{ShopID: shopID}, nil)

		mockRepo.EXPECT().
			GetOrderByTransactionID(gomock.Any(), db.GetOrderByTransactionIDParams{
				TransactionID: &transactionID,
				ShopID:        shopID,
			}).
			Return(db.Order{
				OrderID: 1,
				Status:  db.OrderStatusTypePending,
			}, nil)

		mockRepo.EXPECT().
			UpdateOrderStatusByTransactionID(gomock.Any(), db.UpdateOrderStatusByTransactionIDParams{
				TransactionID: &transactionID,
				ShopID:        shopID,
				Status:        db.OrderStatusTypeProcessing,
				PaymentStatus: db.PaymentStatusTypePaid,
			}).
			Return(db.Order{
				OrderID:       1,
				Status:        db.OrderStatusTypeProcessing,
				PaymentStatus: db.PaymentStatusTypePaid,
			}, nil)

		webhookPayload := &models.PaymentWebhookPayload{
			Provider:      "paystack",
			EventType:     "charge.success",
			PaymentID:     "ch_123456",
			TransactionID: transactionID,
			Status:        "success",
			Amount:        25.00,
			Currency:      "NGN",
			ReceivedAt:    time.Now(),
		}

		mockProcessor.handleWebhookFunc = func(ctx context.Context, payload []byte, signature string) (*models.PaymentWebhookPayload, error) {
			return webhookPayload, nil
		}

		reqBody := `{"event": "charge.success"}`
		req := httptest.NewRequest("POST", "/webhooks/paystack/123", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Paystack-Signature", "test_signature")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestWebhookHandler_FlutterwaveWebhook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockProcessor := &MockPaymentProcessor{}

	processors := map[string]services.PaymentProcessor{
		"flutterwave": mockProcessor,
	}
	mockProcessorFactory := createMockPaymentProcessorFactory(processors)

	handler := NewWebhookHandler(mockProcessorFactory, mockRepo)

	app := fiber.New()
	app.Post("/webhooks/flutterwave/:shop_id", handler.FlutterwaveWebhook)

	t.Run("successful Flutterwave webhook processing", func(t *testing.T) {
		shopID := int64(123)
		transactionID := "fw_123456"

		mockRepo.EXPECT().
			GetShop(gomock.Any(), shopID).
			Return(db.Shop{ShopID: shopID}, nil)

		mockRepo.EXPECT().
			GetOrderByTransactionID(gomock.Any(), db.GetOrderByTransactionIDParams{
				TransactionID: &transactionID,
				ShopID:        shopID,
			}).
			Return(db.Order{
				OrderID: 1,
				Status:  db.OrderStatusTypePending,
			}, nil)

		mockRepo.EXPECT().
			UpdateOrderStatusByTransactionID(gomock.Any(), db.UpdateOrderStatusByTransactionIDParams{
				TransactionID: &transactionID,
				ShopID:        shopID,
				Status:        db.OrderStatusTypeProcessing,
				PaymentStatus: db.PaymentStatusTypePaid,
			}).
			Return(db.Order{
				OrderID:       1,
				Status:        db.OrderStatusTypeProcessing,
				PaymentStatus: db.PaymentStatusTypePaid,
			}, nil)

		webhookPayload := &models.PaymentWebhookPayload{
			Provider:      "flutterwave",
			EventType:     "charge.completed",
			PaymentID:     "flw_123456",
			TransactionID: transactionID,
			Status:        "successful",
			Amount:        75.00,
			Currency:      "USD",
			ReceivedAt:    time.Now(),
		}

		mockProcessor.handleWebhookFunc = func(ctx context.Context, payload []byte, signature string) (*models.PaymentWebhookPayload, error) {
			return webhookPayload, nil
		}

		reqBody := `{"event": "charge.completed"}`
		req := httptest.NewRequest("POST", "/webhooks/flutterwave/123", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("verif-hash", "test_signature")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestWebhookHandler_processWebhookPayload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	processors := map[string]services.PaymentProcessor{
		"stripe": &MockPaymentProcessor{},
	}
	mockProcessorFactory := createMockPaymentProcessorFactory(processors)
	handler := NewWebhookHandler(mockProcessorFactory, mockRepo)

	ctx := context.Background()
	shopID := int64(123)

	t.Run("payment status conversion", func(t *testing.T) {
		testCases := []struct {
			name                string
			paymentStatus       string
			expectedPaymentEnum db.PaymentStatusType
			expectedOrderEnum   db.OrderStatusType
		}{
			{"paid status", "paid", db.PaymentStatusTypePaid, db.OrderStatusTypeProcessing},
			{"completed status", "completed", db.PaymentStatusTypePaid, db.OrderStatusTypeProcessing},
			{"successful status", "successful", db.PaymentStatusTypePaid, db.OrderStatusTypeProcessing},
			{"success status", "success", db.PaymentStatusTypePaid, db.OrderStatusTypeProcessing},
			{"failed status", "failed", db.PaymentStatusTypeFailed, db.OrderStatusTypeCancelled},
			{"failure status", "failure", db.PaymentStatusTypeFailed, db.OrderStatusTypeCancelled},
			{"refunded status", "refunded", db.PaymentStatusTypeRefunded, db.OrderStatusTypeRefunded},
			{"refund status", "refund", db.PaymentStatusTypeRefunded, db.OrderStatusTypeRefunded},
			{"partial_refund status", "partial_refund", db.PaymentStatusTypePartialRefund, db.OrderStatusTypePending},
			{"partially_refunded status", "partially_refunded", db.PaymentStatusTypePartialRefund, db.OrderStatusTypePending},
			{"pending status", "pending", db.PaymentStatusTypePending, db.OrderStatusTypePending},
			{"processing status", "processing", db.PaymentStatusTypePending, db.OrderStatusTypePending},
			{"unknown status", "unknown", db.PaymentStatusTypePending, db.OrderStatusTypePending},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				transactionID := fmt.Sprintf("txn_%s", tc.name)

				mockRepo.EXPECT().
					GetOrderByTransactionID(gomock.Any(), db.GetOrderByTransactionIDParams{
						TransactionID: &transactionID,
						ShopID:        shopID,
					}).
					Return(db.Order{
						OrderID: 1,
						Status:  db.OrderStatusTypePending,
					}, nil)

				mockRepo.EXPECT().
					UpdateOrderStatusByTransactionID(gomock.Any(), db.UpdateOrderStatusByTransactionIDParams{
						TransactionID: &transactionID,
						ShopID:        shopID,
						Status:        tc.expectedOrderEnum,
						PaymentStatus: tc.expectedPaymentEnum,
					}).
					Return(db.Order{
						OrderID:       1,
						Status:        tc.expectedOrderEnum,
						PaymentStatus: tc.expectedPaymentEnum,
					}, nil)

				payload := &models.PaymentWebhookPayload{
					Provider:      "stripe",
					EventType:     "payment_intent.succeeded",
					PaymentID:     "pi_123456",
					TransactionID: transactionID,
					Status:        tc.paymentStatus,
					Amount:        100.00,
					Currency:      "USD",
					ReceivedAt:    time.Now(),
				}

				err := handler.processWebhookPayload(ctx, shopID, payload)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("order not found by transaction ID", func(t *testing.T) {
		transactionID := "nonexistent_txn"

		mockRepo.EXPECT().
			GetOrderByTransactionID(gomock.Any(), db.GetOrderByTransactionIDParams{
				TransactionID: &transactionID,
				ShopID:        shopID,
			}).
			Return(db.Order{}, errors.New("order not found"))

		payload := &models.PaymentWebhookPayload{
			Provider:      "stripe",
			EventType:     "payment_intent.succeeded",
			PaymentID:     "pi_123456",
			TransactionID: transactionID,
			Status:        "paid",
			Amount:        100.00,
			Currency:      "USD",
			ReceivedAt:    time.Now(),
		}

		err := handler.processWebhookPayload(ctx, shopID, payload)
		assert.NoError(t, err) // Should not return error, just log
	})

	t.Run("database update fails", func(t *testing.T) {
		transactionID := "txn_fail"

		mockRepo.EXPECT().
			GetOrderByTransactionID(gomock.Any(), db.GetOrderByTransactionIDParams{
				TransactionID: &transactionID,
				ShopID:        shopID,
			}).
			Return(db.Order{
				OrderID: 1,
				Status:  db.OrderStatusTypePending,
			}, nil)

		mockRepo.EXPECT().
			UpdateOrderStatusByTransactionID(gomock.Any(), db.UpdateOrderStatusByTransactionIDParams{
				TransactionID: &transactionID,
				ShopID:        shopID,
				Status:        db.OrderStatusTypeProcessing,
				PaymentStatus: db.PaymentStatusTypePaid,
			}).
			Return(db.Order{}, errors.New("database error"))

		payload := &models.PaymentWebhookPayload{
			Provider:      "stripe",
			EventType:     "payment_intent.succeeded",
			PaymentID:     "pi_123456",
			TransactionID: transactionID,
			Status:        "paid",
			Amount:        100.00,
			Currency:      "USD",
			ReceivedAt:    time.Now(),
		}

		err := handler.processWebhookPayload(ctx, shopID, payload)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update order status")
	})

	t.Run("empty transaction ID with payment ID", func(t *testing.T) {
		payload := &models.PaymentWebhookPayload{
			Provider:      "stripe",
			EventType:     "payment_intent.succeeded",
			PaymentID:     "pi_123456",
			TransactionID: "", // Empty transaction ID
			Status:        "paid",
			Amount:        100.00,
			Currency:      "USD",
			ReceivedAt:    time.Now(),
		}

		err := handler.processWebhookPayload(ctx, shopID, payload)
		assert.NoError(t, err) // Should not return error, just log unmatched webhook
	})
}

func TestWebhookHandler_handleWebhook_UnsupportedProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	processors := map[string]services.PaymentProcessor{
		"stripe": &MockPaymentProcessor{},
	}
	mockProcessorFactory := createMockPaymentProcessorFactory(processors)
	handler := NewWebhookHandler(mockProcessorFactory, mockRepo)

	app := fiber.New()
	app.Post("/test/:shop_id", func(c *fiber.Ctx) error {
		return handler.handleWebhook(c, "unsupported_provider")
	})

	shopID := int64(123)

	mockRepo.EXPECT().
		GetShop(gomock.Any(), shopID).
		Return(db.Shop{ShopID: shopID}, nil)

	req := httptest.NewRequest("POST", "/test/123", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	var errorResp models.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	require.NoError(t, err)
	assert.Contains(t, errorResp.Message, "Unsupported payment provider")
}

func TestWebhookHandler_SignatureHeaders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	processors := map[string]services.PaymentProcessor{
		"stripe":      &MockPaymentProcessor{},
		"paypal":      &MockPaymentProcessor{},
		"paystack":    &MockPaymentProcessor{},
		"flutterwave": &MockPaymentProcessor{},
	}
	mockProcessorFactory := createMockPaymentProcessorFactory(processors)

	testCases := []struct {
		provider    string
		headerName  string
		headerValue string
		expectedSig string
	}{
		{"stripe", "Stripe-Signature", "t=1234567890,v1=test_signature", "t=1234567890,v1=test_signature"},
		{"paypal", "PAYPAL-TRANSMISSION-SIG", "paypal_signature", "paypal_signature"},
		{"paystack", "X-Paystack-Signature", "paystack_signature", "paystack_signature"},
		{"flutterwave", "verif-hash", "flutterwave_hash", "flutterwave_hash"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s signature header", tc.provider), func(t *testing.T) {
			// Get the processor for this provider from the factory
			processor, ok := processors[tc.provider]
			require.True(t, ok, "Processor not found for provider: %s", tc.provider)

			mockProcessor, ok := processor.(*MockPaymentProcessor)
			require.True(t, ok, "Expected MockPaymentProcessor")

			handler := NewWebhookHandler(mockProcessorFactory, mockRepo)

			app := fiber.New()
			app.Post("/test/:shop_id", func(c *fiber.Ctx) error {
				return handler.handleWebhook(c, tc.provider)
			})

			shopID := int64(123)

			mockRepo.EXPECT().
				GetShop(gomock.Any(), shopID).
				Return(db.Shop{ShopID: shopID}, nil)

			var capturedSignature string
			mockProcessor.handleWebhookFunc = func(ctx context.Context, payload []byte, signature string) (*models.PaymentWebhookPayload, error) {
				capturedSignature = signature
				return &models.PaymentWebhookPayload{
					Provider:   tc.provider,
					EventType:  "test_event",
					PaymentID:  "test_payment",
					Status:     "paid",
					ReceivedAt: time.Now(),
				}, nil
			}

			req := httptest.NewRequest("POST", "/test/123", bytes.NewReader([]byte("{}")))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set(tc.headerName, tc.headerValue)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, tc.expectedSig, capturedSignature)
		})
	}
}
