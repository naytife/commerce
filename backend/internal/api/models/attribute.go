package models

import "github.com/petrejonn/naytife/internal/db"

type Attribute struct {
	Title     string                `json:"title"`
	DataType  db.AttributeDataType  `json:"data_type"`
	Unit      db.NullAttributeUnit  `json:"unit"`
	Required  bool                  `json:"required"`
	AppliesTo db.AttributeAppliesTo `json:"applies_to"`
}
