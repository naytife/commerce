package resolver

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestEncodeRelayID(t *testing.T) {
	// Test integer ID
	intID := int64(123)
	intRelayID := NewIntID("Product", intID)
	intEncodedID := EncodeRelayID(intRelayID)
	assert.NotEmpty(t, intEncodedID)

	// Test string ID
	stringID := "abc123"
	stringRelayID := NewStringID("Token", stringID)
	stringEncodedID := EncodeRelayID(stringRelayID)
	assert.NotEmpty(t, stringEncodedID)

	// Test UUID ID
	uuidID := uuid.New()
	uuidRelayID := NewUUIDID("User", uuidID)
	uuidEncodedID := EncodeRelayID(uuidRelayID)
	assert.NotEmpty(t, uuidEncodedID)
}

func TestDecodeRelayID(t *testing.T) {
	// Test integer ID
	intID := int64(123)
	intEncodedID := EncodeIntID("Product", intID)
	intRelayID, err := DecodeRelayID(intEncodedID)
	assert.NoError(t, err)
	assert.Equal(t, "Product", intRelayID.Type)
	assert.NotNil(t, intRelayID.IntID)
	assert.Equal(t, intID, *intRelayID.IntID)
	assert.Nil(t, intRelayID.StringID)
	assert.Nil(t, intRelayID.UUID)

	// Test string ID
	stringID := "abc123"
	stringEncodedID := EncodeStringID("Token", stringID)
	stringRelayID, err := DecodeRelayID(stringEncodedID)
	assert.NoError(t, err)
	assert.Equal(t, "Token", stringRelayID.Type)
	assert.Nil(t, stringRelayID.IntID)
	assert.NotNil(t, stringRelayID.StringID)
	assert.Equal(t, stringID, *stringRelayID.StringID)
	assert.Nil(t, stringRelayID.UUID)

	// Test UUID ID
	uuidID := uuid.New()
	uuidEncodedID := EncodeUUIDID("User", uuidID)
	uuidRelayID, err := DecodeRelayID(uuidEncodedID)
	assert.NoError(t, err)
	assert.Equal(t, "User", uuidRelayID.Type)
	assert.Nil(t, uuidRelayID.IntID)
	assert.Nil(t, uuidRelayID.StringID)
	assert.NotNil(t, uuidRelayID.UUID)
	assert.Equal(t, uuidID, *uuidRelayID.UUID)

	// Test invalid ID
	_, err = DecodeRelayID("invalid-base64")
	assert.Error(t, err)

	// Test invalid format (missing colon)
	invalidID := "invalid-format"
	invalidEncodedID := EncodeStringID("", invalidID) // This will create an ID with no type
	_, err = DecodeRelayID(invalidEncodedID)
	assert.Error(t, err)
}

func TestConvenienceFunctions(t *testing.T) {
	// Test EncodeIntID
	intID := int64(123)
	intEncodedID := EncodeIntID("Product", intID)
	assert.NotEmpty(t, intEncodedID)
	intRelayID, err := DecodeRelayID(intEncodedID)
	assert.NoError(t, err)
	assert.Equal(t, "Product", intRelayID.Type)
	assert.Equal(t, intID, *intRelayID.IntID)

	// Test EncodeStringID
	stringID := "abc123"
	stringEncodedID := EncodeStringID("Token", stringID)
	assert.NotEmpty(t, stringEncodedID)
	stringRelayID, err := DecodeRelayID(stringEncodedID)
	assert.NoError(t, err)
	assert.Equal(t, "Token", stringRelayID.Type)
	assert.Equal(t, stringID, *stringRelayID.StringID)

	// Test EncodeUUIDID
	uuidID := uuid.New()
	uuidEncodedID := EncodeUUIDID("User", uuidID)
	assert.NotEmpty(t, uuidEncodedID)
	uuidRelayID, err := DecodeRelayID(uuidEncodedID)
	assert.NoError(t, err)
	assert.Equal(t, "User", uuidRelayID.Type)
	assert.Equal(t, uuidID, *uuidRelayID.UUID)
}

func TestPgUUIDConversion(t *testing.T) {
	// Test convert regular UUID to pgtype.UUID
	original := uuid.New()
	pgUUID := GetPgUUID(original)
	assert.True(t, pgUUID.Valid)

	// Test convert pgtype.UUID back to regular UUID
	converted := ParsePgUUID(pgUUID)
	assert.NotNil(t, converted)
	assert.Equal(t, original, *converted)

	// Test invalid pgtype.UUID
	invalidPg := pgtype.UUID{Valid: false}
	invalidConverted := ParsePgUUID(invalidPg)
	assert.Nil(t, invalidConverted)
}

func TestIsUUIDType(t *testing.T) {
	assert.True(t, IsUUIDType("User"))
	assert.False(t, IsUUIDType("Product"))
	assert.False(t, IsUUIDType("Order"))
	assert.False(t, IsUUIDType("Shop"))
}
