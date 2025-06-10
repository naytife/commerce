package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserResponse struct {
	// TODO: add sub fields
	UserID         uuid.UUID          `json:"user_id"`
	Provider       *string            `json:"provider"`
	Email          *string            `json:"email"`
	Name           *string            `json:"name"`
	ProfilePicture *string            `json:"profile_picture"`
	CreatedAt      pgtype.Timestamp   `json:"created_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	LastLogin      pgtype.Timestamp   `json:"last_login" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	ProviderID     *string            `json:"provider_id"`
	Locale         *string            `json:"locale"`
	Shops          []UserShopResponse `json:"shops"`
}

type RegisterUserParams struct {
	Sub            *string `json:"sub"`
	ProviderID     *string `json:"id" validate:"required"`
	Provider       *string `json:"provider"`
	Email          *string `json:"email" validate:"required,email"`
	Name           *string `json:"name" validate:"required,min=3,max=255"`
	Locale         *string `json:"locale"`
	ProfilePicture *string `json:"picture"`
	VerifiedEmail  *bool   `json:"verified_email"`
}

type UserShopResponse struct {
	ID           int64              `json:"shop_id"`
	Title        string             `json:"title"`
	Subdomain    string             `json:"subdomain"`
	CustomDomain string             `json:"custom_domain"`
	Status       string             `json:"status"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	CreatedAt    pgtype.Timestamptz `json:"created_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
}

type RegisterCustomerParams struct {
	CustomerID     *uuid.UUID `json:"customer_id"`
	ShopID         int64      `json:"shop_id" validate:"required"`
	Email          *string    `json:"email" validate:"required,email"`
	Name           *string    `json:"name" validate:"required,min=3,max=255"`
	Locale         *string    `json:"locale"`
	ProfilePicture *string    `json:"picture"`
	VerifiedEmail  *bool      `json:"verified_email"`
	AuthProvider   *string    `json:"auth_provider"`
	AuthProviderID *string    `json:"auth_provider_id"`
}

type CustomerResponse struct {
	CustomerID     uuid.UUID          `json:"customer_id"`
	ShopID         int64              `json:"shop_id"`
	Email          *string            `json:"email"`
	Name           *string            `json:"name"`
	Locale         *string            `json:"locale"`
	ProfilePicture *string            `json:"profile_picture"`
	CreatedAt      pgtype.Timestamptz `json:"created_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	LastLogin      pgtype.Timestamptz `json:"last_login" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	VerifiedEmail  *bool              `json:"verified_email"`
	AuthProvider   *string            `json:"auth_provider"`
	AuthProviderID *string            `json:"auth_provider_id"`
}
