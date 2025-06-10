package resolver

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// RelayID represents a globally unique identifier in the GraphQL API
type RelayID struct {
	Type     string
	IntID    *int64
	StringID *string
	UUID     *uuid.UUID
}

// NewIntID creates a new RelayID with an integer ID
func NewIntID(typeName string, id int64) RelayID {
	return RelayID{
		Type:  typeName,
		IntID: &id,
	}
}

// NewStringID creates a new RelayID with a string ID
func NewStringID(typeName string, id string) RelayID {
	return RelayID{
		Type:     typeName,
		StringID: &id,
	}
}

// NewUUIDID creates a new RelayID with a UUID
func NewUUIDID(typeName string, id uuid.UUID) RelayID {
	return RelayID{
		Type: typeName,
		UUID: &id,
	}
}

// EncodeRelayID encodes a RelayID to a base64 string
func EncodeRelayID(id RelayID) string {
	var rawID string

	if id.IntID != nil {
		rawID = strconv.FormatInt(*id.IntID, 10)
	} else if id.StringID != nil {
		rawID = *id.StringID
	} else if id.UUID != nil {
		rawID = id.UUID.String()
	}

	return base64.StdEncoding.EncodeToString([]byte(id.Type + ":" + rawID))
}

// DecodeRelayID decodes a base64 string into a RelayID
func DecodeRelayID(globalID string) (RelayID, error) {
	bytes, err := base64.StdEncoding.DecodeString(globalID)
	if err != nil {
		return RelayID{}, err
	}

	parts := strings.SplitN(string(bytes), ":", 2)
	if len(parts) != 2 {
		return RelayID{}, errors.New("invalid global ID format")
	}

	typeName := parts[0]
	rawID := parts[1]

	// Try to parse as integer first
	if intID, err := strconv.ParseInt(rawID, 10, 64); err == nil {
		return NewIntID(typeName, intID), nil
	}

	// Try to parse as UUID
	if u, err := uuid.Parse(rawID); err == nil {
		return NewUUIDID(typeName, u), nil
	}

	// Default to string ID
	return NewStringID(typeName, rawID), nil
}

// EncodeIntID is a convenient wrapper for encoding integer IDs
func EncodeIntID(typeName string, id int64) string {
	return EncodeRelayID(NewIntID(typeName, id))
}

// EncodeStringID is a convenient wrapper for encoding string IDs
func EncodeStringID(typeName string, id string) string {
	return EncodeRelayID(NewStringID(typeName, id))
}

// EncodeUUIDID is a convenient wrapper for encoding UUID IDs
func EncodeUUIDID(typeName string, id uuid.UUID) string {
	return EncodeRelayID(NewUUIDID(typeName, id))
}

// ParsePgUUID converts a pgtype.UUID to a regular uuid.UUID
func ParsePgUUID(pgUUID pgtype.UUID) *uuid.UUID {
	if !pgUUID.Valid {
		return nil
	}

	u, err := uuid.FromBytes(pgUUID.Bytes[:])
	if err != nil {
		return nil
	}

	return &u
}

// GetPgUUID converts a uuid.UUID to pgtype.UUID
func GetPgUUID(id uuid.UUID) pgtype.UUID {
	var pgUUID pgtype.UUID
	pgUUID.Scan(id)
	return pgUUID
}

// IsUUIDType returns true if the type typically uses UUIDs as IDs
func IsUUIDType(typeName string) bool {
	// Add all types that use UUIDs here
	return typeName == "User"
}

// decodeRelayID is a helper that extracts the type and integer ID from a base64 relay ID
func decodeRelayID(globalID string) (string, *int64, error) {
	relayID, err := DecodeRelayID(globalID)
	if err != nil {
		return "", nil, err
	}
	return relayID.Type, relayID.IntID, nil
}

// encodeRelayID is a helper that encodes a type and ID string to a base64 relay ID
func encodeRelayID(typeName string, id string) string {
	return EncodeStringID(typeName, id)
}
