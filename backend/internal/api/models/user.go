package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserResponse struct {
	// TODO: add sub fields
	UserID         uuid.UUID        `json:"user_id"`
	Provider       *string          `json:"provider"`
	Email          *string          `json:"email"`
	Name           *string          `json:"name"`
	ProfilePicture *string          `json:"profile_picture"`
	CreatedAt      pgtype.Timestamp `json:"created_at" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	LastLogin      pgtype.Timestamp `json:"last_login" swaggertype:"primitive,string" format:"date-time" example:"2025-02-09T09:38:25Z"`
	ProviderID     *string          `json:"provider_id"`
	Locale         *string          `json:"locale"`
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
