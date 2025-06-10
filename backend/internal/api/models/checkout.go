package models

import "time"

// CheckoutRequest represents a checkout initiation request
type CheckoutRequest struct {
	Items           []CartItem       `json:"items" validate:"required,min=1"`
	CustomerInfo    CustomerInfo     `json:"customer_info" validate:"required"`
	ShippingAddress ShippingAddress  `json:"shipping_address" validate:"required"`
	BillingAddress  *ShippingAddress `json:"billing_address,omitempty"`
	ShippingCost    float64          `json:"shipping_cost" validate:"min=0"`
	TaxRate         float64          `json:"tax_rate" validate:"min=0,max=1"`
	Discount        float64          `json:"discount" validate:"min=0"`
	CouponCode      *string          `json:"coupon_code,omitempty"`
	Notes           *string          `json:"notes,omitempty"`
}

// CheckoutResponse represents the checkout response
type CheckoutResponse struct {
	ShopID          int64               `json:"shop_id"`
	ShopName        string              `json:"shop_name"`
	CurrencyCode    string              `json:"currency_code"`
	Subtotal        float64             `json:"subtotal"`
	Tax             float64             `json:"tax"`
	Shipping        float64             `json:"shipping"`
	Discount        float64             `json:"discount"`
	Total           float64             `json:"total"`
	PaymentMethods  []PaymentMethodInfo `json:"payment_methods"`
	CustomerInfo    CustomerInfo        `json:"customer_info"`
	ShippingAddress ShippingAddress     `json:"shipping_address"`
	BillingAddress  *ShippingAddress    `json:"billing_address,omitempty"`
	Items           []CartItem          `json:"items"`
	SessionID       *string             `json:"session_id,omitempty"`
	ExpiresAt       time.Time           `json:"expires_at"`
}

// CustomerInfo represents customer information for checkout
type CustomerInfo struct {
	Email     string  `json:"email" validate:"required,email"`
	FirstName string  `json:"first_name" validate:"required"`
	LastName  string  `json:"last_name" validate:"required"`
	Phone     *string `json:"phone,omitempty"`
}

// ShippingAddress represents a shipping/billing address
type ShippingAddress struct {
	FirstName    string  `json:"first_name" validate:"required"`
	LastName     string  `json:"last_name" validate:"required"`
	Company      *string `json:"company,omitempty"`
	AddressLine1 string  `json:"address_line_1" validate:"required"`
	AddressLine2 *string `json:"address_line_2,omitempty"`
	City         string  `json:"city" validate:"required"`
	State        *string `json:"state,omitempty"`
	PostalCode   string  `json:"postal_code" validate:"required"`
	Country      string  `json:"country" validate:"required"`
	Phone        *string `json:"phone,omitempty"`
}
