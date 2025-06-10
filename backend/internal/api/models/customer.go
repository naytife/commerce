package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petrejonn/naytife/internal/db"
)

// CustomerListResponse represents a paginated list of customers
type CustomerListResponse struct {
	Customers []CustomerResponse `json:"customers"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	Limit     int                `json:"limit"`
}

// UpdateCustomerParams represents the parameters for updating a customer
type UpdateCustomerParams struct {
	Name           *string `json:"name" validate:"omitempty,min=3,max=255"`
	Locale         *string `json:"locale"`
	ProfilePicture *string `json:"profile_picture"`
	VerifiedEmail  *bool   `json:"verified_email"`
}

// OrderResponse represents an order in API responses
type OrderResponse struct {
	OrderID         int64               `json:"order_id"`
	Status          string              `json:"status"`
	Amount          float64             `json:"amount"`
	Discount        float64             `json:"discount"`
	ShippingCost    float64             `json:"shipping_cost"`
	Tax             float64             `json:"tax"`
	ShippingAddress string              `json:"shipping_address"`
	PaymentMethod   string              `json:"payment_method"`
	PaymentStatus   string              `json:"payment_status"`
	ShippingMethod  string              `json:"shipping_method"`
	ShippingStatus  string              `json:"shipping_status"`
	TransactionID   *string             `json:"transaction_id"`
	Username        string              `json:"username"`
	CustomerName    string              `json:"customer_name"`
	CustomerEmail   *string             `json:"customer_email"`
	CustomerPhone   *string             `json:"customer_phone"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	Items           []OrderItemResponse `json:"items"`
}

// OrderItemResponse represents an order item in API responses
type OrderItemResponse struct {
	OrderItemID        int64   `json:"order_item_id"`
	ProductVariationID int64   `json:"product_variation_id"`
	Quantity           int64   `json:"quantity"`
	Price              float64 `json:"price"`
	ProductTitle       string  `json:"product_title"`
	VariantTitle       *string `json:"variant_title"`
}

// CustomerOrdersResponse represents the response for customer orders
type CustomerOrdersResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int             `json:"total"`
}

// Helper function to convert pgtype.Numeric to float64
func NumericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, _ := n.Float64Value()
	return f.Float64
}

// Helper function to convert pgtype.Timestamptz to time.Time
func TimestamptzToTime(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// Helper function to convert db.OrderStatusType to string
func OrderStatusToString(status db.OrderStatusType) string {
	return string(status)
}

// Helper function to convert db.PaymentMethodType to string
func PaymentMethodToString(method db.PaymentMethodType) string {
	return string(method)
}

// Helper function to convert db.PaymentStatusType to string
func PaymentStatusToString(status db.PaymentStatusType) string {
	return string(status)
}

// Helper function to convert db.ShippingStatusType to string
func ShippingStatusToString(status db.ShippingStatusType) string {
	return string(status)
}
