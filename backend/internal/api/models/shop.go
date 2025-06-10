package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// ShopImagesData holds the actual image URLs (renamed from the original ShopImagesResponse to avoid conflict)
type ShopImagesData struct {
	ID                int64   `json:"shop_image_id,omitempty"`
	FaviconUrl        *string `json:"favicon_url,omitempty"`
	LogoUrl           *string `json:"logo_url,omitempty"`
	LogoUrlDark       *string `json:"logo_url_dark,omitempty"`
	BannerUrl         *string `json:"banner_url,omitempty"`
	BannerUrlDark     *string `json:"banner_url_dark,omitempty"`
	CoverImageUrl     *string `json:"cover_image_url,omitempty"`
	CoverImageUrlDark *string `json:"cover_image_url_dark,omitempty"`
}

type ShopCreateParams struct {
	Title        string `json:"title" validate:"required,min=3,max=255"`
	Subdomain    string `json:"subdomain" validate:"required,min=3,max=255"`
	CurrencyCode string `json:"currency_code" validate:"required,oneof=USD NGN"`
	Status       string `json:"status" validate:"required,oneof=PUBLISHED DRAFT"`
}

type Shop struct {
	ID                  int64              `json:"shop_id"`
	Title               string             `json:"title"`
	Subdomain           string             `json:"subdomain"`
	CustomDomain        string             `json:"custom_domain"`
	Email               string             `json:"email,omitempty"`
	CurrencyCode        string             `json:"currency_code"`
	Status              string             `json:"status"`
	About               *string            `json:"about,omitempty"`
	Address             *string            `json:"address,omitempty"`
	PhoneNumber         *string            `json:"phone_number,omitempty"`
	WhatsappPhoneNumber *string            `json:"whatsapp_phone_number,omitempty"`
	WhatsappLink        *string            `json:"whatsapp_link,omitempty"`
	FacebookLink        *string            `json:"facebook_link,omitempty"`
	InstagramLink       *string            `json:"instagram_link,omitempty"`
	SeoDescription      *string            `json:"seo_description,omitempty"`
	SeoKeywords         []string           `json:"seo_keywords,omitempty"`
	SeoTitle            *string            `json:"seo_title,omitempty"`
	UpdatedAt           pgtype.Timestamptz `json:"updated_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	CreatedAt           pgtype.Timestamptz `json:"created_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	Images              *ShopImagesData    `json:"images,omitempty"`
}

type ShopUpdateParams struct {
	Title               *string  `json:"title"`
	CurrencyCode        *string  `json:"currency_code"`
	About               *string  `json:"about"`
	Status              *string  `json:"status"`
	PhoneNumber         *string  `json:"phone_number"`
	WhatsappLink        *string  `json:"whatsapp_link"`
	WhatsappPhoneNumber *string  `json:"whatsapp_phone_number"`
	FacebookLink        *string  `json:"facebook_link"`
	InstagramLink       *string  `json:"instagram_link"`
	SeoDescription      *string  `json:"seo_description"`
	SeoKeywords         []string `json:"seo_keywords"`
	SeoTitle            *string  `json:"seo_title"`
	Address             *string  `json:"address"`
	Email               *string  `json:"email"`
}

type ShopImagesUpdateParams struct {
	FaviconUrl        *string `json:"favicon_url"`
	LogoUrl           *string `json:"logo_url"`
	LogoUrlDark       *string `json:"logo_url_dark"`
	BannerUrl         *string `json:"banner_url"`
	BannerUrlDark     *string `json:"banner_url_dark"`
	CoverImageUrl     *string `json:"cover_image_url"`
	CoverImageUrlDark *string `json:"cover_image_url_dark"`
}
