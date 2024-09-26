package resolver

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repository db.Repository
}

func DecodeRelayID(globalID string) (string, *int64, error) {
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
func EncodeRelayID(typ string, id string) string {
	return base64.StdEncoding.EncodeToString([]byte(typ + ":" + id))
}

func pgTextFromStringPointer(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func unmarshalCategoryAttributes(attributesDB []byte) ([]model.AllowedCategoryAttributes, error) {
	// Unmarshal JSONB ([]byte) into a Go map
	var attributesMap map[string]interface{}
	if err := json.Unmarshal(attributesDB, &attributesMap); err != nil {
		return nil, err
	}

	// Create a list to hold the allowed category attributes
	attributes := make([]model.AllowedCategoryAttributes, 0, len(attributesMap))

	// Iterate over the map and populate the attribute list
	for title, dataType := range attributesMap {
		attributes = append(attributes, model.AllowedCategoryAttributes{
			Title:    title,                                             // The map key is the title (string)
			DataType: model.ProductAttributeDataType(dataType.(string)), // The map value is the data type
		})
	}

	return attributes, nil
}

func unmarshalAllowedProductAttributes(attributesDB []byte) ([]model.AllowedProductAttributes, error) {
	// Unmarshal JSONB ([]byte) into a Go map
	var attributesMap map[string]interface{}
	if err := json.Unmarshal(attributesDB, &attributesMap); err != nil {
		return nil, err
	}

	// Create a list to hold the allowed category attributes
	attributes := make([]model.AllowedProductAttributes, 0, len(attributesMap))

	// Iterate over the map and populate the attribute list
	for title, dataType := range attributesMap {
		attributes = append(attributes, model.AllowedProductAttributes{
			Title:    title,                                             // The map key is the title (string)
			DataType: model.ProductAttributeDataType(dataType.(string)), // The map value is the data type
		})
	}

	return attributes, nil
}
