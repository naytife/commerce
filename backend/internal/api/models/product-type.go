package models

type ProductType struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	Shippable    bool    `json:"shippable"`
	Digital      bool    `json:"digital"`
	SkuSubstring *string `json:"sku_substring,omitempty"`
}

type ProductTypeCreateParams struct {
	Title        string  `json:"title" example:"Book"`
	Shippable    bool    `json:"shippable"`
	Digital      bool    `json:"digital" example:"false"`
	SkuSubstring *string `json:"sku_substring,omitempty" example:"BK"`
}

type ProductTypeUpdateParams struct {
	Title        *string `json:"title" example:"Book"`
	Shippable    *bool   `json:"shippable"`
	Digital      *bool   `json:"digital" example:"false"`
	SkuSubstring *string `json:"sku_substring,omitempty" example:"BK"`
}
