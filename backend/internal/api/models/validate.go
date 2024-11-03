package models

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type XValidator struct{}

func (v XValidator) Validate(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

func FormatValidationErrors(errs []ErrorResponse) string {
	var errMsgs []string
	for _, err := range errs {
		errMsgs = append(errMsgs, fmt.Sprintf("[%s]: '%v' | Needs to implement '%s'", err.FailedField, err.Value, err.Tag))
	}
	return strings.Join(errMsgs, " and ")
}
