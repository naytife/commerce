package models

type ProductType struct {
	ProductTypeID int64  `json:"product_type_id"`
	Title         string `json:"title"`
	Shippable     bool   `json:"shippable"`
	Digital       bool   `json:"digital"`
	ShopID        int64  `json:"shop_id"`
}

type ProductTypeCreate struct {
	ProductTypeID int64  `json:"product_type_id"`
	Title         string `json:"title"`
	Shippable     bool   `json:"shippable"`
	Digital       bool   `json:"digital"`
	ShopID        int64  `json:"shop_id"`
}

type ProductTypeUpdate struct {
	ProductTypeID int64  `json:"product_type_id"`
	Title         string `json:"title"`
	Shippable     bool   `json:"shippable"`
	Digital       bool   `json:"digital"`
	ShopID        int64  `json:"shop_id"`
}
