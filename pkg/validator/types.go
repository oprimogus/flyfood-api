package validator

import "github.com/go-playground/validator/v10"

type Validation struct {
	TagName    string
	Func       validator.Func
	ErrMessage string
}