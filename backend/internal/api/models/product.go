package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petrejonn/naytife/internal/db"
)

type Product struct {
	ProductID   int64              `json:"product_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Status      db.ProductStatus   `json:"status"`
	CategoryID  int64              `json:"category_id"`
	Attributes  []ProductAttribute `json:"attributes"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	CreatedAt   pgtype.Timestamptz `json:"created_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
}

type ProductAttribute struct {
	AttributeID       int64   `json:"attribute_id"`
	AttributeOptionID *int64  `json:"attribute_option_id"`
	Value             *string `json:"value"`
}

type ProductCreateParams struct {
	Title       string                                    `json:"title"`
	Description string                                    `json:"description"`
	Attributes  []ProductAttributeValuesBatchUpsertParams `json:"attributes"`
}

type ProductAttributeValuesBatchUpsertParams struct {
	Value             *string `json:"value,omitempty"`
	AttributeOptionID *int64  `json:"attribute_option_id,omitempty"`
	AttributeID       int64   `json:"attribute_id"`
}

type ProductUpdateParams struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
