package models

type ProductType struct {
	ProductTypeID int64  `json:"product_type_id"`
	Title         string `json:"title"`
	Shippable     bool   `json:"shippable"`
	Digital       bool   `json:"digital"`
}

type ProductTypeCreateParams struct {
	Title     string `json:"title" example:"Book"`
	Shippable bool   `json:"shippable"`
	Digital   bool   `json:"digital" example:"false"`
}

type ProductTypeUpdateParams struct {
	Title     string `json:"title" example:"Book"`
	Shippable bool   `json:"shippable"`
	Digital   bool   `json:"digital" example:"false"`
}
