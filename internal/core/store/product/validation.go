package product

import (
	"github.com/go-playground/validator/v10"
	"github.com/oprimogus/cardapiogo/internal/xvalidator"
)

var validationsMap = map[string]xvalidator.PersonalizedValidation{
	"productType": xvalidator.
		NewPersonalizedValidation("productType", "Tipo de produto inválido", IsValidTypeValidation),
}

func IsValidTypeValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	return IsValidType(value)
}
