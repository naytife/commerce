package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petrejonn/naytife/internal/db"
)

// Order represents the API order response
type Order struct {
	ID              int64                 `json:"order_id"`
	Status          db.OrderStatusType    `json:"status"`
	CreatedAt       pgtype.Timestamptz    `json:"created_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	UpdatedAt       pgtype.Timestamptz    `json:"updated_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	CustomerID      uuid.UUID             `json:"customer_id"`
	Amount          pgtype.Numeric        `json:"amount" swaggertype:"primitive,number"`
	Discount        pgtype.Numeric        `json:"discount" swaggertype:"primitive,number"`
	ShippingCost    pgtype.Numeric        `json:"shipping_cost" swaggertype:"primitive,number"`
	Tax             pgtype.Numeric        `json:"tax" swaggertype:"primitive,number"`
	ShippingAddress string                `json:"shipping_address"`
	PaymentMethod   db.PaymentMethodType  `json:"payment_method"`
	PaymentStatus   db.PaymentStatusType  `json:"payment_status"`
	ShippingMethod  string                `json:"shipping_method"`
	ShippingStatus  db.ShippingStatusType `json:"shipping_status"`
	TransactionID   *string               `json:"transaction_id"`
	Username        string                `json:"username"`
	ShopID          int64                 `json:"shop_id"`
	Items           []OrderItem           `json:"items,omitempty"`
	// Customer contact information
	CustomerName  string  `json:"customer_name"`
	CustomerEmail *string `json:"customer_email,omitempty"`
	CustomerPhone *string `json:"customer_phone,omitempty"`
}

// OrderItem represents the API order item response
type OrderItem struct {
	ID                 int64              `json:"order_item_id"`
	Quantity           int64              `json:"quantity"`
	Price              pgtype.Numeric     `json:"price" swaggertype:"primitive,number"`
	CreatedAt          pgtype.Timestamptz `json:"created_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	UpdatedAt          pgtype.Timestamptz `json:"updated_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	ProductVariationID int64              `json:"product_variation_id"`
	OrderID            int64              `json:"order_id"`
	ShopID             int64              `json:"shop_id"`
}

// CreateOrderParams represents the request body for creating an order
type CreateOrderParams struct {
	Status          db.OrderStatusType      `json:"status" validate:"required"`
	Amount          pgtype.Numeric          `json:"amount" validate:"required" swaggertype:"primitive,number"`
	Discount        pgtype.Numeric          `json:"discount" swaggertype:"primitive,number"`
	ShippingCost    pgtype.Numeric          `json:"shipping_cost" swaggertype:"primitive,number"`
	Tax             pgtype.Numeric          `json:"tax" swaggertype:"primitive,number"`
	ShippingAddress string                  `json:"shipping_address" validate:"required"`
	PaymentMethod   db.PaymentMethodType    `json:"payment_method" validate:"required"`
	PaymentStatus   db.PaymentStatusType    `json:"payment_status" validate:"required"`
	ShippingMethod  string                  `json:"shipping_method" validate:"required"`
	ShippingStatus  db.ShippingStatusType   `json:"shipping_status" validate:"required"`
	TransactionID   *string                 `json:"transaction_id"`
	Username        string                  `json:"username" validate:"required"`
	Items           []CreateOrderItemParams `json:"items" validate:"required,min=1"`
}

// CreateOrderItemParams represents the request body for creating an order item
type CreateOrderItemParams struct {
	Quantity           int64          `json:"quantity" validate:"required,min=1"`
	Price              pgtype.Numeric `json:"price" validate:"required" swaggertype:"primitive,number"`
	ProductVariationID int64          `json:"product_variation_id" validate:"required"`
}

// UpdateOrderParams represents the request body for updating an order
type UpdateOrderParams struct {
	Status          db.OrderStatusType    `json:"status" validate:"required"`
	Amount          pgtype.Numeric        `json:"amount" validate:"required" swaggertype:"primitive,number"`
	Discount        pgtype.Numeric        `json:"discount" swaggertype:"primitive,number"`
	ShippingCost    pgtype.Numeric        `json:"shipping_cost" swaggertype:"primitive,number"`
	Tax             pgtype.Numeric        `json:"tax" swaggertype:"primitive,number"`
	ShippingAddress string                `json:"shipping_address" validate:"required"`
	PaymentMethod   db.PaymentMethodType  `json:"payment_method" validate:"required"`
	PaymentStatus   db.PaymentStatusType  `json:"payment_status" validate:"required"`
	ShippingMethod  string                `json:"shipping_method" validate:"required"`
	ShippingStatus  db.ShippingStatusType `json:"shipping_status" validate:"required"`
	TransactionID   *string               `json:"transaction_id"`
	Username        string                `json:"username" validate:"required"`
	// Customer contact information
	CustomerName  string  `json:"customer_name"`
	CustomerEmail *string `json:"customer_email,omitempty"`
	CustomerPhone *string `json:"customer_phone,omitempty"`
}

// UpdateOrderStatusParams represents the request body for updating an order status
type UpdateOrderStatusParams struct {
	Status db.OrderStatusType `json:"status" validate:"required"`
}

// CreateOrderRequest represents the API request for creating an order from cart
type CreateOrderRequest struct {
	CustomerID      *string                  `json:"customer_id,omitempty"`
	CustomerName    string                   `json:"customer_name" validate:"required"`
	CustomerEmail   *string                  `json:"customer_email,omitempty"`
	CustomerPhone   *string                  `json:"customer_phone,omitempty"`
	ShippingAddress string                   `json:"shipping_address" validate:"required"`
	ShippingMethod  string                   `json:"shipping_method" validate:"required"`
	PaymentMethod   string                   `json:"payment_method" validate:"required,oneof=flutterwave paystack paypal stripe"`
	TransactionID   *string                  `json:"transaction_id,omitempty"`
	Discount        float64                  `json:"discount,omitempty"`
	ShippingCost    float64                  `json:"shipping_cost,omitempty"`
	Tax             float64                  `json:"tax,omitempty"`
	Items           []CreateOrderRequestItem `json:"items" validate:"required,min=1"`
}

// CreateOrderRequestItem represents an item in the order creation request
type CreateOrderRequestItem struct {
	ProductVariationID string  `json:"product_variation_id" validate:"required"`
	Quantity           int     `json:"quantity" validate:"required,min=1"`
	Price              float64 `json:"price" validate:"required,min=0"`
}
