package store

import (
	"github.com/go-playground/validator/v10"
	xvalidator "github.com/oprimogus/flyfood-api/pkg/validator"
)

func init() {
	err := xvalidator.AddNewValidation(xvalidator.NewValidation("storeType", "Tipo de loja inválido", IsValidTypeValidation))
	if err != nil {
		panic(err)
	}
	err = xvalidator.AddNewValidation(xvalidator.NewValidation("paymentMethod", "Forma de pagamento inválida", IsValidPaymentMethodValidation))
	if err != nil {
		panic(err)
	}
	err = xvalidator.AddNewValidation(xvalidator.NewValidation("businessHour", "Hora inválida", IsValidHourValidation))
	if err != nil {
		panic(err)
	}
}

func IsValidPaymentMethodValidation(fl validator.FieldLevel) bool {
	value := fl.Field().Int()
	if value == 0 {
		return true
	}
	return IsValidPaymentMethod(int(value))
}

func IsValidTypeValidation(fl validator.FieldLevel) bool {
	value := fl.Field().Int()
	if value == 0 {
		return true
	}
	return IsValidType(int(value))
}

func IsValidHourValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return isValidHour(value)
}
