package xvalidator

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/pt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	ptTranslations "github.com/go-playground/validator/v10/translations/pt"
)

type Validator struct {
	Validator  *validator.Validate
	translator ut.Translator
	locale     string
}

var (
	once             sync.Once
	validatorService *Validator
)

func init() {
	if err := getInstance("pt"); err != nil {
		panic(err)
	}
}

func newValidator(locale string) (*Validator, error) {
	if locale == "" {
		locale = "pt"
	}

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
		if err := enTranslations.RegisterDefaultTranslations(v, translator); err != nil {
			return nil, fmt.Errorf("error registering en translations: %w", err)
		}
	case "pt":
		if err := ptTranslations.RegisterDefaultTranslations(v, translator); err != nil {
			return nil, fmt.Errorf("error registering pt translations: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported locale: %s", locale)
	}

	validatorService = &Validator{
		Validator:  v,
		translator: translator,
		locale:     locale,
	}

	err := AddValidations(personalizedValidations)
	if err != nil {
		return nil, fmt.Errorf("error adding personalized validations: %w", err)
	}

	return validatorService, nil
}

func getInstance(locale string) error {
	var err error
	once.Do(func() {
		validatorService, err = newValidator(locale)
	})
	return err
}

func AddValidations(validations map[string]PersonalizedValidation) error {
	if validatorService == nil {
		err := getInstance("pt")
		if err != nil {
			return err
		}
	}

	for i, validation := range validations {
		err := validatorService.Validator.RegisterValidation(i, validation.ValidationFn)
		if err != nil {
			return fmt.Errorf("could not create validator for type %s: %v", i, err)
		}
		err = validatorService.Validator.RegisterTranslation(i, validatorService.translator,
			validation.RegisterTranslationsFn, validation.TranslationFn)
		if err != nil {
			return fmt.Errorf("could not create translation for type %s: %v", i, err)
		}
	}
	return nil
}

func Validate(i any) error {
	if validatorService == nil {
		err := getInstance("pt")
		if err != nil {
			return err
		}
	}

	if err := validatorService.Validator.Struct(i); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			mapErrFields := make([]FieldError, len(errs))
			for i, value := range errs {
				input, ok := value.Value().(string)
				if !ok {
					input = ""
				}
				mapErrFields[i] = FieldError{
					Field:   value.Field(),
					Input:   input,
					Message: strings.Replace(value.Translate(validatorService.translator), value.Field()+" ", "", 1),
				}
			}
			return NewValidationError(mapErrFields)
		}
	}
	return nil
}
