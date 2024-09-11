package resolver

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petrejonn/naytife/internal/db"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repository db.Repository
}

func fromGlobalID(globalID string) (string, string, error) {
	bytes, err := base64.StdEncoding.DecodeString(globalID)
	if err != nil {
		return "", "", err
	}
	parts := strings.SplitN(string(bytes), ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid global ID")
	}
	return parts[0], parts[1], nil
}
func toGlobalID(typ, id string) string {
	return base64.StdEncoding.EncodeToString([]byte(typ + ":" + id))
}

func pgTextFromStringPointer(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}
