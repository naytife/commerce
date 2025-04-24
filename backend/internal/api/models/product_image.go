package models

// ProductImage represents a product image
type ProductImage struct {
	ID        int64  `json:"product_image_id"`
	URL       string `json:"url"`
	Alt       string `json:"alt"`
	ProductID int64  `json:"product_id"`
}

// ProductImageCreateParams contains parameters for creating a product image
type ProductImageCreateParams struct {
	URL string `json:"url" validate:"required,url"`
	Alt string `json:"alt" validate:"required"`
}

// ProductImageResponse represents a product image response
type ProductImageResponse struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
	Alt string `json:"alt"`
}

// ProductImagesResponse represents a list of product images
type ProductImagesResponse struct {
	Images []ProductImageResponse `json:"images"`
}
