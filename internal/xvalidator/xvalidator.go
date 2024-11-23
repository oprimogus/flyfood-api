package xvalidator

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/pt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	ptTranslations "github.com/go-playground/validator/v10/translations/pt"
	"strings"
)

var validatorService *Validator

type Validator struct {
	Validator  *validator.Validate
	translator ut.Translator
	locale     string
}

func (v *Validator) AddValidations(validations map[string]PersonalizedValidation) error {
	for i, validation := range validations {
		err := v.Validator.RegisterValidation(i, validation.ValidationFn)
		if err != nil {
			return fmt.Errorf("could not create validator for type %s: %v", i, err)
		}
		err = v.Validator.RegisterTranslation(i, v.translator,
			validation.RegisterTranslationsFn, validation.TranslationFn)
		if err != nil {
			return fmt.Errorf("could not create translation for type %s: %v", i, err)
		}
	}
	return nil
}

func NewValidator(locale string) (*Validator, error) {
	v := validator.New(validator.WithRequiredStructEnabled())

	enLocale := en.New()
	ptLocale := pt.New()
	uni := ut.New(enLocale, ptLocale, enLocale)

	translator, found := uni.GetTranslator(locale)
	if !found {
		return nil, fmt.Errorf("locale %s not found", locale)
	}
	switch locale {
	case "en":
		err := enTranslations.RegisterDefaultTranslations(v, translator)
		if err != nil {
			panic(fmt.Sprintf("Could not register locale %v translation: %v", locale, err))
		}
	case "pt":
		err := ptTranslations.RegisterDefaultTranslations(v, translator)
		if err != nil {
			panic(fmt.Sprintf("Could not register locale %v translation: %v", locale, err))
		}
	default:
		return nil, fmt.Errorf("unsupported locale: %s", locale)
	}

	validatorInstance := &Validator{
		Validator:  v,
		translator: translator,
		locale:     locale,
	}

	err := validatorInstance.AddValidations(personalizedValidations)
	if err != nil {
		return nil, err
	}
	return validatorInstance, nil
}

func GetInstance(locale string) *Validator {
	if validatorService == nil {
		v, err := NewValidator(locale)
		if err != nil {
			panic(err)
		}
		validatorService = v
	}
	if validatorService.locale != locale {
		v, err := NewValidator(locale)
		if err != nil {
			panic(err)
		}
		validatorService = v
	}
	return validatorService
}

func GetPtInstance() *Validator {
	if validatorService == nil {
		v, err := NewValidator("pt")
		if err != nil {
			panic(err)
		}
		validatorService = v
	}
	return validatorService
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.Validator.Struct(i); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			mapErrFields := make([]ErrField, len(errs))
			for i, value := range errs {
				input, ok := value.Value().(string)
				if !ok {
					input = ""
				}
				mapErrFields[i] = ErrField{
					Field:   value.Field(),
					Input:   input,
					Message: strings.Replace(value.Translate(v.translator), value.Field()+" ", "", 1),
				}
			}
			return NewValidationError(mapErrFields)
		}
	}
	return nil
}
