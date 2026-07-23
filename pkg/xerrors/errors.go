package xerrors

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/oprimogus/flyfood-api/internal/infra/logger"
)

type CustomError struct {
	Status  int    `json:"-"`
	Message string `json:"error"`
	Details any    `json:"details,omitempty"`
	TraceID string `json:"traceID"`
	Debug   any    `json:"debug,omitempty"`
}

func New(message string) *CustomError {
	if message == "" {
		message = "Ocorreu um erro inesperado. Tente novamente mais tarde."
	}
	return &CustomError{
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}

func (e *CustomError) Error() string {
	return e.Message
}

func (e *CustomError) StatusCode() int {
	return e.Status
}

func (e *CustomError) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("message", e.Message),
		slog.Int("status", e.Status),
		slog.Any("details", e.Details),
		slog.Any("debug", e.Debug),
	)
}

type FieldError struct {
	Field   string `json:"field"`
	Input   string `json:"input"`
	Message string `json:"message"`
	Debug   any    `json:"debug,omitempty"`
}

func NewWithContext(ctx context.Context, err error) *CustomError {
	data := logger.GetRequestContext(ctx)
	customErr := &CustomError{
        Message: err.Error(),
		TraceID: data.TraceID,
		Status: http.StatusInternalServerError,
	}
	return customErr
}

func (e *CustomError) WithStatus(status int) *CustomError {
	e.Status = status
	return e
}

func (e *CustomError) WithStatusBadRequest() *CustomError {
	e.Status = http.StatusBadRequest
	return e
}

func (e *CustomError) WithStatusNotFound() *CustomError {
	e.Status = http.StatusNotFound
	return e
}

func (e *CustomError) WithStatusConflict() *CustomError {
	e.Status = http.StatusConflict
	return e
}

func (e *CustomError) WithStatusForbidden() *CustomError {
	e.Status = http.StatusForbidden
	return e
}

func (e *CustomError) WithContext(ctx context.Context) *CustomError {
	data := logger.GetRequestContext(ctx)
	e.TraceID = data.TraceID
	return e
}

func (e *CustomError) WithDetails(details any) *CustomError {
	e.Details = details
	return e
}