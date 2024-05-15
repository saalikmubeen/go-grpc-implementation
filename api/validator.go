package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/saalikmubeen/backend-masterclass-go/utils"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return utils.IsSupportedCurrency(currency)
	}
	return false
}
