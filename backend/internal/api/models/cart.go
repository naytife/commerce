package models

import "time"

// Cart represents a shopping cart
type Cart struct {
	Items []CartItem `json:"items"`
	Total float64    `json:"total"`
}

// CartItem represents an item in the shopping cart
type CartItem struct {
	ID                 string    `json:"id"`
	ProductVariationID string    `json:"product_variation_id"`
	Quantity           int       `json:"quantity"`
	Title              string    `json:"title"`
	Price              float64   `json:"price"`
	Image              *string   `json:"image,omitempty"`
	Slug               *string   `json:"slug,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// AddToCartRequest represents a request to add an item to cart
type AddToCartRequest struct {
	ProductVariationID string  `json:"product_variation_id" validate:"required"`
	Quantity           int     `json:"quantity" validate:"required,min=1"`
	Title              string  `json:"title" validate:"required"`
	Price              float64 `json:"price" validate:"required,min=0"`
	Image              *string `json:"image,omitempty"`
	Slug               *string `json:"slug,omitempty"`
}

// UpdateCartItemRequest represents a request to update cart item quantity
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}

// CartTotals represents cart total calculations
type CartTotals struct {
	Subtotal   float64 `json:"subtotal"`
	Tax        float64 `json:"tax"`
	Shipping   float64 `json:"shipping"`
	Discount   float64 `json:"discount"`
	Total      float64 `json:"total"`
	ItemCount  int     `json:"item_count"`
	TotalItems int     `json:"total_items"` // Sum of all quantities
}
