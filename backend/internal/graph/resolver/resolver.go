package resolver

import (
	"encoding/base64"
	"errors"
	"strconv"
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

func fromGlobalID(globalID string) (string, *int64, error) {
	bytes, err := base64.StdEncoding.DecodeString(globalID)
	if err != nil {
		return "", nil, err
	}
	parts := strings.SplitN(string(bytes), ":", 2)
	if len(parts) != 2 {
		return "", nil, errors.New("invalid global ID")
	}
	key, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", nil, err
	}
	keyInt64 := int64(key)
	return parts[0], &keyInt64, nil
}
func toGlobalID(typ string, id int64) string {
	return base64.StdEncoding.EncodeToString([]byte(typ + ":" + strconv.Itoa(int(id))))
}

func pgTextFromStringPointer(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}
