package handlers

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	// Set up gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Use the generated mock from gomock
	mockRepo := mocks.NewMockRepository(ctrl)
	handler := &Handler{Repository: mockRepo}

	app := fiber.New()
	app.Post("/shops/:shop_id/product-types/:product_type_id/products", handler.CreateProduct)

	tests := []struct {
		name           string
		payload        models.ProductCreateParams
		expectedStatus int
		mockSetup      func()
	}{
		{
			name: "Valid request",
			payload: models.ProductCreateParams{
				Title:       "Test Product",
				Description: "Test Description",
				Variants: []models.ProductVariantParams{
					{
						Sku:               "TEST-SKU",
						Price:             pgtype.Numeric{Int: big.NewInt(10000), Valid: true},
						AvailableQuantity: 10,
					},
				},
				Attributes: []models.ProductAttributeValuesBatchUpsertParams{},
			},
			expectedStatus: fiber.StatusCreated,
			mockSetup: func() {
				mockRepo.
					EXPECT().
					WithTx(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "Invalid request body",
			payload: models.ProductCreateParams{
				Title:       "",
				Description: "",
				Variants:    []models.ProductVariantParams{},
			},
			expectedStatus: fiber.StatusBadRequest,
			mockSetup:      func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/shops/1/product-types/1/products", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
