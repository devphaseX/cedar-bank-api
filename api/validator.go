package api

import (
	"github.com/devphasex/cedar-bank-api/util"
	"github.com/go-playground/validator/v10"
)

var currencyValidator validator.Func = func(fl validator.FieldLevel) bool {
	return util.IsCurrencySupported(fl.Field().String())
}

// Register custom validators
func registerCustomValidators(v *validator.Validate) {
	v.RegisterValidation("currency", currencyValidator)
}
