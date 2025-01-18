package store

import (
	"github.com/go-playground/validator/v10"
	//"github.com/oprimogus/cardapiogo/internal/core/product"
	"github.com/oprimogus/cardapiogo/internal/xvalidator"
)

var validationsMap = map[string]xvalidator.PersonalizedValidation{
	"storeType": xvalidator.
		NewPersonalizedValidation("storeType", "Tipo de loja inválido", IsValidTypeValidation),
	//"productType": xvalidator.
	//	NewPersonalizedValidation("productType", "Tipo de produto inválido", IsValidProductTypeValidation),
	"paymentMethod": xvalidator.
		NewPersonalizedValidation("paymentMethod", "Forma de pagamento inválida", IsValidPaymentMethodValidation),
	"businessHour": xvalidator.
		NewPersonalizedValidation("businessHour", "Hora inválida", IsValidHourValidation),
}

//func IsValidProductTypeValidation(fl validator.FieldLevel) bool {
//	value := fl.Field().String()
//	if value == "" {
//		return true
//	}
//	return product.IsValidType(value)
//}

func IsValidPaymentMethodValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	return IsValidPaymentMethod(value)
}

func IsValidTypeValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	return IsValidType(value)
}

func IsValidHourValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return isValidHour(value)
}
