package xvalidator

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/pt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	ptTranslations "github.com/go-playground/validator/v10/translations/pt"

	xerrors "github.com/oprimogus/cardapiogo/internal/errors"
)

var personalizedValidations = map[string]func(fl validator.FieldLevel) bool{
	"role":          isValidUserRole,
	"cpf":           IsValidCpf,
	"cnpj":          IsValidCnpj,
	"cpfCnpj":       IsValidCpfOrCnpj,
	"shopType":      IsValidShopType,
	"paymentMethod": IsValidPaymentMethod,
	"phone":         IsValidPhone,
	"businessHour":  IsValidBusinessHour,
	"weekDay":       isValidWeekDay,
	"itemType":      IsValidItemType,
}

type Validator struct {
	Validator  *validator.Validate
	translator ut.Translator
	locale     string
}

func NewValidator(locale string) (*Validator, error) {
	v := validator.New(validator.WithRequiredStructEnabled())

	for i, validation := range personalizedValidations {
		err := v.RegisterValidation(i, validation)
		if err != nil {
			panic(fmt.Sprintf("Could not create validator for type %s: %v", i, err))
		}
	}

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

	return &Validator{
		Validator:  v,
		translator: translator,
		locale:     locale,
	}, nil
}

func (v *Validator) Validate(traceID string, i interface{}) error {
	out := make(map[string]string)

	err := v.Validator.Struct(i)
	if err != nil {
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			out["error"] = "Unknown validation error"
			return xerrors.New(traceID, http.StatusBadRequest, out["error"])
		}

		for _, e := range errs {
			_, isPersonalized := personalizedValidations[e.Tag()]
			if isPersonalized {
				out[e.StructField()] = errorPersonalized(v.locale, e.Tag())
			} else {
				out[e.StructField()] = e.Translate(v.translator)
			}
		}
	}

	if len(out) > 0 {
		return xerrors.InvalidInput(traceID, out)
	}
	return nil
}

func errorPersonalized(locale string, tag string) string {
	if locale == "pt" {
		return "Valor inválido"
	}
	return fmt.Sprintf("Invalid value for field %v", tag)
}
