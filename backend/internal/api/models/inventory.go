package models

import (
	"time"
)

// LowStockVariantResponse represents a product variant with low stock
type LowStockVariantResponse struct {
	VariantID    int64   `json:"variant_id"`
	ProductID    int64   `json:"product_id"`
	ProductTitle string  `json:"product_title"`
	VariantTitle *string `json:"variant_title"`
	SKU          *string `json:"sku"`
	CurrentStock int32   `json:"current_stock"`
	ReorderLevel int32   `json:"reorder_level"`
	Price        float64 `json:"price"`
}

// UpdateStockParams represents parameters for updating stock
type UpdateStockParams struct {
	Quantity     int32    `json:"quantity" validate:"required,min=0"`
	MovementType string   `json:"movement_type" validate:"required,oneof=adjustment purchase damage transfer"`
	Reason       *string  `json:"reason,omitempty" validate:"omitempty,min=3,max=255"`
	CostPrice    *float64 `json:"cost_price,omitempty" validate:"omitempty,min=0"`
}

// AddStockParams represents parameters for adding stock
type AddStockParams struct {
	Quantity      int32   `json:"quantity" validate:"required,min=1"`
	Reason        string  `json:"reason" validate:"required,min=3,max=255"`
	ReferenceType string  `json:"reference_type" validate:"required"`
	ReferenceID   *string `json:"reference_id"`
}

// DeductStockParams represents parameters for deducting stock
type DeductStockParams struct {
	Quantity      int32   `json:"quantity" validate:"required,min=1"`
	Reason        string  `json:"reason" validate:"required,min=3,max=255"`
	ReferenceType string  `json:"reference_type" validate:"required"`
	ReferenceID   *string `json:"reference_id"`
}

// VariantStockResponse represents the response after stock operations
type VariantStockResponse struct {
	VariantID int64     `json:"variant_id"`
	Stock     int32     `json:"stock"`
	UpdatedAt time.Time `json:"updated_at"`
}

// InventoryReportResponse represents an inventory report
type InventoryReportResponse struct {
	TotalProducts   int                     `json:"total_products"`
	TotalVariants   int                     `json:"total_variants"`
	TotalStockValue float64                 `json:"total_stock_value"`
	LowStockCount   int                     `json:"low_stock_count"`
	OutOfStockCount int                     `json:"out_of_stock_count"`
	GeneratedAt     time.Time               `json:"generated_at"`
	Items           []InventoryItemResponse `json:"items"`
}

// InventoryItemResponse represents an inventory item in the report
type InventoryItemResponse struct {
	VariantID         int64   `json:"variant_id"`
	ProductID         int64   `json:"product_id"`
	ProductTitle      string  `json:"product_title"`
	VariantTitle      string  `json:"variant_title"`
	SKU               string  `json:"sku"`
	CurrentStock      int64   `json:"current_stock"`
	ReservedStock     int64   `json:"reserved_stock"`
	AvailableStock    int64   `json:"available_stock"`
	LowStockThreshold int64   `json:"low_stock_threshold"`
	Location          *string `json:"location"`
	LastUpdated       string  `json:"last_updated"`
}

// StockMovementResponse represents a stock movement record
type StockMovementResponse struct {
	MovementID    int64     `json:"movement_id"`
	VariantID     int64     `json:"variant_id"`
	MovementType  string    `json:"movement_type"`
	Quantity      int32     `json:"quantity"`
	PreviousStock *int32    `json:"previous_stock"`
	NewStock      *int32    `json:"new_stock"`
	Reason        *string   `json:"reason"`
	ReferenceType *string   `json:"reference_type"`
	ReferenceID   *string   `json:"reference_id"`
	CreatedAt     time.Time `json:"created_at"`
	ProductTitle  string    `json:"product_title"`
	VariantTitle  *string   `json:"variant_title"`
	SKU           string    `json:"sku"`
}

// StockMovementsResponse represents a paginated list of stock movements
type StockMovementsResponse struct {
	Movements []StockMovementResponse `json:"movements"`
	Total     int                     `json:"total"`
	Page      int                     `json:"page"`
	Limit     int                     `json:"limit"`
}
