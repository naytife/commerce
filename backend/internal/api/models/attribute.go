package models

import "github.com/petrejonn/naytife/internal/db"

type Attribute struct {
	ID            int64                 `json:"attribute_id"`
	Title         string                `json:"title"`
	DataType      db.AttributeDataType  `json:"data_type"`
	Unit          db.AttributeUnit      `json:"unit,omitempty"`
	Required      bool                  `json:"required"`
	AppliesTo     db.AttributeAppliesTo `json:"applies_to"`
	ProductTypeID int64                 `json:"product_type_id"`
	Options       []AttributeOption     `json:"options,omitempty"`
}

type AttributeCreateParams struct {
	Title     string                `json:"title" validate:"required,min=3,max=255" example:"Size"`
	DataType  db.AttributeDataType  `json:"data_type" validate:"required,oneof=Text Number Date" example:"Text"`
	Unit      db.AttributeUnit      `json:"unit,omitempty" example:"KG"`
	Required  bool                  `json:"required" validate:"boolean"`
	AppliesTo db.AttributeAppliesTo `json:"applies_to" validate:"required,oneof=Product ProductVariation"`
}

type AttributeUpdateParams struct {
	Title     *string               `json:"title,omitempty" validate:"omitempty,min=3,max=255" example:"Size"`
	DataType  db.AttributeDataType  `json:"data_type,omitempty" validate:"omitempty,oneof=Text Number Date" example:"Text"`
	Unit      db.AttributeUnit      `json:"unit,omitempty" example:"KG"`
	Required  *bool                 `json:"required,omitempty"`
	AppliesTo db.AttributeAppliesTo `json:"applies_to,omitempty" validate:"omitempty,oneof=Product ProductVariation"`
}

type AttributeOption struct {
	ID          int64  `json:"attribute_option_id"`
	Value       string `json:"value"`
	AttributeID int64  `json:"attribute_id,omitempty"`
}

type AttributeOptionCreateParams struct {
	Value string `json:"value" validate:"required,min=1,max=255" example:"XL"`
}

type AttributeOptionUpdateParams struct {
	Value *string `json:"value,omitempty" validate:"omitempty,min=1,max=255" example:"XL"`
}
