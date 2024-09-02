// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Node interface {
	IsNode()
	GetID() string
}

type SocialMediaContact interface {
	IsSocialMediaContact()
	GetURL() *string
}

type AllowedProductAttributes struct {
	Key      string                    `json:"key"`
	DataType *ProductAttributeDataType `json:"dataType,omitempty"`
	Options  []ProductAttributeValue   `json:"options"`
}

type Category struct {
	ID                string                     `json:"id"`
	Slug              string                     `json:"slug"`
	Title             string                     `json:"title"`
	Description       *string                    `json:"description,omitempty"`
	Parent            *Category                  `json:"parent,omitempty"`
	Children          []*Category                `json:"children"`
	Products          *ProductConnection         `json:"products,omitempty"`
	AllowedAttributes []AllowedProductAttributes `json:"allowedAttributes"`
	Image             *CategoryImage             `json:"image,omitempty"`
	UpdatedAt         *string                    `json:"updatedAt,omitempty"`
	CreatedAt         *string                    `json:"createdAt,omitempty"`
}

func (Category) IsNode()            {}
func (this Category) GetID() string { return this.ID }

type CategoryImage struct {
	URL string `json:"url"`
}

type CreateShopInput struct {
	Title  string `json:"title"`
	Domain string `json:"domain"`
}

type CreateShopPayload struct {
	Shop       *Shop `json:"shop,omitempty"`
	Successful bool  `json:"successful"`
}

type Facebook struct {
	URL    string `json:"url"`
	Handle string `json:"handle"`
}

func (Facebook) IsSocialMediaContact() {}
func (this Facebook) GetURL() *string  { return &this.URL }

type Instagram struct {
	URL    string `json:"url"`
	Handle string `json:"handle"`
}

func (Instagram) IsSocialMediaContact() {}
func (this Instagram) GetURL() *string  { return &this.URL }

type Location struct {
	Address string `json:"address"`
	State   string `json:"state"`
	Country string `json:"country"`
}

type Mutation struct {
}

type PageInfo struct {
	EndCursor       *string `json:"endCursor,omitempty"`
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor,omitempty"`
}

type PhoneNumber struct {
	Number      string `json:"number"`
	CountryCode string `json:"countryCode"`
}

type Product struct {
	ID                string                     `json:"id"`
	Slug              string                     `json:"slug"`
	Title             string                     `json:"title"`
	Price             float64                    `json:"price"`
	Description       string                     `json:"description"`
	Category          *Category                  `json:"category"`
	DefaultVariant    *ProductVariant            `json:"defaultVariant"`
	Variants          []ProductVariant           `json:"variants"`
	AllowedAttributes []AllowedProductAttributes `json:"allowedAttributes"`
	Images            []ProductImage             `json:"images"`
	Status            *ProductStatus             `json:"status,omitempty"`
	UpdatedAt         *string                    `json:"updatedAt,omitempty"`
	CreatedAt         *string                    `json:"createdAt,omitempty"`
}

func (Product) IsNode()            {}
func (this Product) GetID() string { return this.ID }

type ProductAttribute struct {
	Key   string  `json:"key"`
	Value *string `json:"value,omitempty"`
}

type ProductAttributeValue struct {
	IntValue    *int    `json:"intValue,omitempty"`
	StringValue *string `json:"stringValue,omitempty"`
}

type ProductConnection struct {
	Edges    []ProductEdge `json:"edges"`
	PageInfo *PageInfo     `json:"pageInfo"`
}

type ProductEdge struct {
	Cursor string   `json:"cursor"`
	Node   *Product `json:"node"`
}

type ProductImage struct {
	URL string `json:"url"`
}

type ProductVariant struct {
	ID          string              `json:"id"`
	Slug        string              `json:"slug"`
	Title       *string             `json:"title,omitempty"`
	Price       *float64            `json:"price,omitempty"`
	Quantity    *int                `json:"quantity,omitempty"`
	Description *string             `json:"description,omitempty"`
	Attributes  []ProductAttribute  `json:"attributes"`
	StockStatus *ProductStockStatus `json:"stockStatus,omitempty"`
	UpdatedAt   *string             `json:"updatedAt,omitempty"`
	CreatedAt   *string             `json:"createdAt,omitempty"`
}

func (ProductVariant) IsNode()            {}
func (this ProductVariant) GetID() string { return this.ID }

type Query struct {
}

type Shop struct {
	ID             string             `json:"id"`
	Title          string             `json:"title"`
	DefaultDomain  string             `json:"defaultDomain"`
	ContactPhone   *PhoneNumber       `json:"contactPhone,omitempty"`
	ContactEmail   *string            `json:"contactEmail,omitempty"`
	Location       *Location          `json:"location,omitempty"`
	Products       *ProductConnection `json:"products,omitempty"`
	WhatsApp       *WhatsApp          `json:"whatsApp,omitempty"`
	Facebook       *Facebook          `json:"facebook,omitempty"`
	SiteLogoURL    *string            `json:"siteLogoUrl,omitempty"`
	FaviconURL     *string            `json:"faviconUrl,omitempty"`
	CurrencyCode   *string            `json:"currencyCode,omitempty"`
	Status         *ShopStatus        `json:"status,omitempty"`
	About          *string            `json:"about,omitempty"`
	SeoDescription *string            `json:"seoDescription,omitempty"`
	SeoKeywords    []string           `json:"seoKeywords"`
	SeoTitle       *string            `json:"seoTitle,omitempty"`
	UpdatedAt      *string            `json:"updatedAt,omitempty"`
	CreatedAt      *string            `json:"createdAt,omitempty"`
}

func (Shop) IsNode()            {}
func (this Shop) GetID() string { return this.ID }

type WhatsApp struct {
	URL    string       `json:"url"`
	Number *PhoneNumber `json:"number"`
}

func (WhatsApp) IsSocialMediaContact() {}
func (this WhatsApp) GetURL() *string  { return &this.URL }

type ProductAttributeDataType string

const (
	ProductAttributeDataTypeString  ProductAttributeDataType = "STRING"
	ProductAttributeDataTypeInteger ProductAttributeDataType = "INTEGER"
)

var AllProductAttributeDataType = []ProductAttributeDataType{
	ProductAttributeDataTypeString,
	ProductAttributeDataTypeInteger,
}

func (e ProductAttributeDataType) IsValid() bool {
	switch e {
	case ProductAttributeDataTypeString, ProductAttributeDataTypeInteger:
		return true
	}
	return false
}

func (e ProductAttributeDataType) String() string {
	return string(e)
}

func (e *ProductAttributeDataType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProductAttributeDataType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProductAttributeDataType", str)
	}
	return nil
}

func (e ProductAttributeDataType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "DRAFT"
	ProductStatusPublished ProductStatus = "PUBLISHED"
	ProductStatusArchived  ProductStatus = "ARCHIVED"
)

var AllProductStatus = []ProductStatus{
	ProductStatusDraft,
	ProductStatusPublished,
	ProductStatusArchived,
}

func (e ProductStatus) IsValid() bool {
	switch e {
	case ProductStatusDraft, ProductStatusPublished, ProductStatusArchived:
		return true
	}
	return false
}

func (e ProductStatus) String() string {
	return string(e)
}

func (e *ProductStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProductStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProductStatus", str)
	}
	return nil
}

func (e ProductStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ProductStockStatus string

const (
	ProductStockStatusInStock    ProductStockStatus = "IN_STOCK"
	ProductStockStatusOutOfStock ProductStockStatus = "OUT_OF_STOCK"
	ProductStockStatusPreorder   ProductStockStatus = "PREORDER"
)

var AllProductStockStatus = []ProductStockStatus{
	ProductStockStatusInStock,
	ProductStockStatusOutOfStock,
	ProductStockStatusPreorder,
}

func (e ProductStockStatus) IsValid() bool {
	switch e {
	case ProductStockStatusInStock, ProductStockStatusOutOfStock, ProductStockStatusPreorder:
		return true
	}
	return false
}

func (e ProductStockStatus) String() string {
	return string(e)
}

func (e *ProductStockStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProductStockStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProductStockStatus", str)
	}
	return nil
}

func (e ProductStockStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ShopStatus string

const (
	ShopStatusDraft     ShopStatus = "DRAFT"
	ShopStatusPublished ShopStatus = "PUBLISHED"
	ShopStatusArchived  ShopStatus = "ARCHIVED"
)

var AllShopStatus = []ShopStatus{
	ShopStatusDraft,
	ShopStatusPublished,
	ShopStatusArchived,
}

func (e ShopStatus) IsValid() bool {
	switch e {
	case ShopStatusDraft, ShopStatusPublished, ShopStatusArchived:
		return true
	}
	return false
}

func (e ShopStatus) String() string {
	return string(e)
}

func (e *ShopStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ShopStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ShopStatus", str)
	}
	return nil
}

func (e ShopStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}