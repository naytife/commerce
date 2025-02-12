package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ShopCreateParams struct {
	Title        string `json:"title" validate:"required,min=3,max=255"`
	Domain       string `json:"domain" validate:"required,min=3,max=255"`
	CurrencyCode string `json:"currency_code" validate:"required,oneof=USD NGN"`
	Status       string `json:"status" validate:"required,oneof=PUBLISHED DRAFT"`
}

type Shop struct {
	ShopID              int64              `json:"shop_id"`
	Title               string             `json:"title"`
	Domain              string             `json:"domain"`
	Email               string             `json:"email"`
	CurrencyCode        string             `json:"currency_code"`
	Status              string             `json:"status"`
	About               *string            `json:"about"`
	Address             *string            `json:"address"`
	PhoneNumber         *string            `json:"phone_number"`
	WhatsappPhoneNumber *string            `json:"whatsapp_phone_number"`
	WhatsappLink        *string            `json:"whatsapp_link"`
	FacebookLink        *string            `json:"facebook_link"`
	InstagramLink       *string            `json:"instagram_link"`
	SeoDescription      *string            `json:"seo_description"`
	SeoKeywords         []string           `json:"seo_keywords"`
	SeoTitle            *string            `json:"seo_title"`
	UpdatedAt           pgtype.Timestamptz `json:"updated_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	CreatedAt           pgtype.Timestamptz `json:"created_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
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
