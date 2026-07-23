package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var (
    v *xvalidator
)

type xvalidator struct {
    *validator.Validate
}

func NewValidator() *xvalidator{
    validator := validator.New(validator.WithRequiredStructEnabled())
    return &xvalidator{validator}
}

func init() {
    v = NewValidator()
    AddNewValidation(NewValidation("cpf", "CPF inválido", IsValidCpf))
    AddNewValidation(NewValidation("cnpj", "CNPJ inválido", IsValidCnpj))
    AddNewValidation(NewValidation("cpfCnpj", "CPF/CNPJ inválido", IsValidCpfOrCnpj))
    AddNewValidation(NewValidation("phone", "número de telefone inválido", IsValidPhone))
    AddNewValidation(NewValidation("week", "dia de semana inválido", isValidWeekDay))
}

func Validate(data any) error {
	err := v.Validate.Struct(data)
	if err != nil {
		if errs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			mapErrFields := make([]FieldError, len(errs))
			for i, v := range errs {
				input, ok := v.Value().(string)
				if !ok {
					input = ""
				}
				mapErrFields[i] = FieldError{
					Field:   v.Field(),
					Input:   input,
					Message: v.Error(),
				}
			}
			return NewValidationError(mapErrFields)
		}
		return err
	}
	return nil
}

func NewValidation(tagName string, errMessage string, fn validator.Func) Validation {
	return Validation{tagName, fn, errMessage}
}

func AddNewValidation(validation Validation) error {
	if err := v.RegisterValidation(validation.TagName, validation.Func); err != nil {
		return err
	}
	return nil
}

