package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/taufiqdp/go-simplebank/utils"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	currency, ok := fl.Field().Interface().(string)
	if ok {
		return utils.IsSupportedCurrency(currency)
	}

	return false
}
