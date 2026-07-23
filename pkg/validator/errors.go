package validator

import (
	"fmt"
	"log/slog"
)

const (
	MessageInvalidInput = "Há um ou mais erros com os dados enviados"
)

type FieldError struct {
	Field   string `json:"field"`
	Input   string `json:"input"`
	Message string `json:"message"`
	Debug   any    `json:"debug,omitempty"`
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("o valor %s é inválido para o campo %s: %s", e.Input, e.Field, e.Message)
}

func (e *FieldError) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("field", e.Field),
		slog.String("input", e.Input),
		slog.String("message", e.Message),
		slog.Any("debug", e.Debug),
	)
}

type ValidationError struct {
	Message string       `json:"message"`
	Fields  []FieldError `json:"fields"`
}

func (e *ValidationError) Error() string {
	return MessageInvalidInput
}

func (e *ValidationError) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("message", e.Message),
		slog.Any("fields", e.Fields),
	)
}

func NewValidationError(fields []FieldError) *ValidationError {
	return &ValidationError{
		Message: MessageInvalidInput,
		Fields:  fields,
	}
}

func NewFieldError(field, value string) *ValidationError {
	return &ValidationError{
		Message: MessageInvalidInput,
		Fields: []FieldError{
			{
				Field:   field,
				Input:   value,
				Message: "Valor inválido para o campo",
			},
		},
	}
}
