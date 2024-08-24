package api

import (
	"fmt"
	"strings"

	"github.com/devphasex/cedar-bank-api/util"
	"github.com/go-playground/validator/v10"
)

// Custom error handler for Gin
func prettyValidateError(err error) FailedResponse {
	errs := err.(validator.ValidationErrors)
	var errMessages []string
	for _, e := range errs {
		switch e.Tag() {
		case "currency":
			currencyErr := &util.UnsupportedCurrencyError{Currency: e.Value().(string)}
			errMessages = append(errMessages, currencyErr.Error())
		default:
			errMessages = append(errMessages, formatValidationError(e))
		}
	}
	return errorResponse(errMessages)
}

// Helper function to format validation errors
func formatValidationError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())
	switch e.Tag() {
	case "required":
		return field + " is required"
	// Add more cases for other validation tags as needed
	default:
		return fmt.Sprintf("%s failed on the '%s' tag", field, e.Tag())
	}
}
