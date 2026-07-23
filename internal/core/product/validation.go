package product

import (
	"github.com/go-playground/validator/v10"
	xvalidator "github.com/oprimogus/flyfood-api/pkg/validator"
)

func IsValidTypeValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	return IsValidType(value)
}

func init() {
	err := xvalidator.AddNewValidation(xvalidator.NewValidation("productType", "Tipo de produto inválido", IsValidTypeValidation))
	if err != nil {
		panic(err)
	}
}
