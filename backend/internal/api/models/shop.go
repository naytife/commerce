package models

import "github.com/google/uuid"

type ShopCreate struct {
	OwnerID      uuid.UUID `json:"owner_id" validate:"required,uuid"`
	Title        string    `json:"title" validate:"required,min=3,max=255"`
	Domain       string    `json:"domain" validate:"required,min=3,max=255"`
	CurrencyCode string    `json:"currency_code" validate:"required,oneof=USD NGN"`
	Status       string    `json:"status" validate:"required,oneof=PUBLISHED DRAFT"`
}
