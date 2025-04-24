package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petrejonn/naytife/internal/db"
)

type Product struct {
	ID             int64                  `json:"product_id"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description"`
	Status         db.ProductStatus       `json:"status"`
	CategoryID     int64                  `json:"category_id"`
	Attributes     []ProductAttribute     `json:"attributes"`
	Variants       []ProductVariant       `json:"variants"`
	DefaultVariant ProductVariant         `json:"default_variant"`
	Images         []ProductImageResponse `json:"images"`
	UpdatedAt      pgtype.Timestamptz     `json:"updated_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	CreatedAt      pgtype.Timestamptz     `json:"created_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
}

type ProductAttribute struct {
	AttributeID       int64   `json:"attribute_id"`
	AttributeOptionID *int64  `json:"attribute_option_id,omitempty"`
	Value             *string `json:"value"`
	AttributeTitle    string  `json:"title"`
}

type ProductVariant struct {
	VariationID       int64              `json:"id"`
	Sku               string             `json:"sku"`
	Slug              string             `json:"slug"`
	Description       string             `json:"description"`
	Price             pgtype.Numeric     `json:"price" swaggertype:"primitive,number"`
	IsDefault         bool               `json:"is_default"`
	AvailableQuantity int64              `json:"available_quantity"`
	SeoDescription    *string            `json:"seo_description"`
	SeoKeywords       []string           `json:"seo_keywords"`
	SeoTitle          *string            `json:"seo_title"`
	Attributes        []ProductAttribute `json:"attributes"`
	CreatedAt         pgtype.Timestamptz `json:"created_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	UpdatedAt         pgtype.Timestamptz `json:"updated_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
}

type ProductCreateParams struct {
	Title       string                                    `json:"title"`
	Description string                                    `json:"description"`
	Attributes  []ProductAttributeValuesBatchUpsertParams `json:"attributes" validate:"required"`
	Variants    []ProductVariantParams                    `json:"variants" validate:"required"`
}

type ProductAttributeValuesBatchUpsertParams struct {
	Value             *string `json:"value,omitempty"`
	AttributeOptionID *int64  `json:"attribute_option_id,omitempty"`
	AttributeID       int64   `json:"attribute_id" validate:"required"`
}

type ProductVariantParams struct {
	// SKU will be generated automatically from product type's sku_substring and variation_id
	Description       string                                    `json:"description"`
	Price             pgtype.Numeric                            `json:"price" validate:"required" swaggertype:"primitive,number"`
	AvailableQuantity int64                                     `json:"available_quantity" validate:"required,min=1"`
	SeoDescription    *string                                   `json:"seo_description"`
	SeoKeywords       []string                                  `json:"seo_keywords"`
	SeoTitle          *string                                   `json:"seo_title"`
	Attributes        []ProductAttributeValuesBatchUpsertParams `json:"attributes"`
	IsDefault         bool                                      `json:"is_default"`
}

type ProductUpdateParams struct {
	Title       *string                                   `json:"title"`
	Description *string                                   `json:"description"`
	Attributes  []ProductAttributeValuesBatchUpsertParams `json:"attributes"`
	Variants    []ProductVariantUpdateParams              `json:"variants"`
}

type ProductVariantUpdateParams struct {
	ID                int64                                     `json:"id,omitempty"`
	Description       string                                    `json:"description"`
	Price             pgtype.Numeric                            `json:"price" validate:"required" swaggertype:"primitive,number"`
	AvailableQuantity int64                                     `json:"available_quantity" validate:"required,min=1"`
	SeoDescription    *string                                   `json:"seo_description"`
	SeoKeywords       []string                                  `json:"seo_keywords"`
	SeoTitle          *string                                   `json:"seo_title"`
	Attributes        []ProductAttributeValuesBatchUpsertParams `json:"attributes"`
	IsDefault         bool                                      `json:"is_default"`
}
