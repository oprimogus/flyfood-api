package xvalidator

import "fmt"

var (
	MessageInvalidInput = "Há um ou mais erros com os dados enviados"
)

type ErrField struct {
	Field   string      `json:"field"`
	Input   string      `json:"input"`
	Message string      `json:"message"`
	Debug   interface{} `json:"debug,omitempty"`
}

func (e *ErrField) Error() string {
	return fmt.Sprintf("o valor %s é inválido para o campo %s: %s", e.Input, e.Field, e.Message)
}

type ValidationError struct {
	Message string     `json:"message"`
	Fields  []ErrField `json:"fields"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s", MessageInvalidInput)
}

func NewValidationError(fields []ErrField) *ValidationError {
	return &ValidationError{
		Message: MessageInvalidInput,
		Fields:  fields,
	}
}
