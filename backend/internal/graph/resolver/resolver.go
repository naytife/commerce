package resolver

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
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

func decodeRelayID(globalID string) (string, *int64, error) {
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
func encodeRelayID(typ string, id string) string {
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
func unmarshalProductAttributes(attributesDB []byte) ([]model.ProductAttribute, error) {
	var attributes []model.ProductAttribute

	// Attempt to unmarshal the attributesDB JSON into the attributes slice
	err := json.Unmarshal(attributesDB, &attributes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal product attributes: %w", err)
	}

	// Return the slice of attributes and nil (no error)
	return attributes, nil
}

func isValidE164(phoneNumber string) bool {
	e164Regex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return e164Regex.MatchString(phoneNumber)
}

func validateProductVariantInput(input model.CreateProductVariantInput, attributesMap map[string]interface{}) error {
	for _, attr := range input.Attributes {
		expectedType, ok := attributesMap[attr.Key]
		if !ok {
			return fmt.Errorf("attribute '%s' is not recognized", attr.Key)
		}

		// Check the expected type
		switch expectedType {
		case "STRING":
			// Any value is acceptable for STRING
		case "INTEGER":
			if _, err := strconv.Atoi(attr.Value); err != nil {
				return fmt.Errorf("attribute '%s' expects an INTEGER value, but got '%s'", attr.Key, attr.Value)
			}
		default:
			return fmt.Errorf("unsupported type '%s' for attribute '%s'", expectedType, attr.Key)
		}
	}
	return nil
}

func formatAttributes(attributes []model.ProductAttributeValueInput) string {
	// Step 1: Sort the attributes by key to ensure consistent order
	sort.Slice(attributes, func(i, j int) bool {
		return attributes[i].Key < attributes[j].Key
	})

	// Step 2: Build the desired string "value key, value key"
	var attributePairs []string
	for _, attr := range attributes {
		attributePairs = append(attributePairs, fmt.Sprintf("%s %s", attr.Value, attr.Key))
	}

	// Step 3: Join the results with commas
	return strings.Join(attributePairs, ", ")
}
