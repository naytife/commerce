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

func (v XValidator) Validate(data interface{}) []Error {
	validationErrors := []Error{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem Error

			elem.Field = err.Field()   // Export struct field name
			elem.Message = err.Error() // Export struct tag
			elem.Code = "VALIDATION_ERROR"

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

func FormatValidationErrors(errs []Error) string {
	var errMsgs []string
	for _, err := range errs {
		errMsgs = append(errMsgs, fmt.Sprintf("[%s]: '%v' | Needs to implement '%s'", err.Field, err.Message, err.Message))
	}
	return strings.Join(errMsgs, " and ")
}
