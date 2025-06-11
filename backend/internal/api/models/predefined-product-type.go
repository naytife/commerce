package models

// PredefinedProductType represents a template for creating product types
type PredefinedProductType struct {
	ID           string                        `json:"id"`
	Title        string                        `json:"title"`
	Description  string                        `json:"description"`
	SkuSubstring string                        `json:"sku_substring"`
	Shippable    bool                          `json:"shippable"`
	Digital      bool                          `json:"digital"`
	Category     string                        `json:"category"`
	Icon         string                        `json:"icon"`
	Attributes   []PredefinedAttributeTemplate `json:"attributes"`
}

// PredefinedAttributeTemplate represents an attribute template for predefined product types
type PredefinedAttributeTemplate struct {
	Title     string                      `json:"title"`
	DataType  string                      `json:"data_type"` // Text, Number, Date, Option, Color
	Unit      *string                     `json:"unit,omitempty"`
	Required  bool                        `json:"required"`
	AppliesTo string                      `json:"applies_to"` // Product, ProductVariation
	Options   []PredefinedAttributeOption `json:"options,omitempty"`
}

// PredefinedAttributeOption represents an option for predefined attributes
type PredefinedAttributeOption struct {
	Value string `json:"value"`
}

// CreateProductTypeFromTemplateParams represents the request to create a product type from a template
type CreateProductTypeFromTemplateParams struct {
	TemplateID string `json:"template_id" validate:"required"`
}

// ProductTypeWithTemplateResponse represents the response when creating from template
type ProductTypeWithTemplateResponse struct {
	ProductType ProductType `json:"product_type"`
	Attributes  []Attribute `json:"attributes"`
}
