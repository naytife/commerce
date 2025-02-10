package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserResponse struct {
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
