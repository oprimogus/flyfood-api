package xerrors

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/oprimogus/flyfood-api/internal/config"
	logger "github.com/oprimogus/flyfood-api/pkg/log"
)

type CustomError struct {
	Status       int    `json:"-"`
	ErrorMessage string `json:"error"`
	Details      any    `json:"details,omitempty"`
	TraceID      string `json:"traceID"`
	Debug        any   `json:"debug,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Input   string `json:"input"`
	Message string `json:"message"`
	Debug   any    `json:"debug,omitempty"`
}

func (e *CustomError) Error() string {
	return e.ErrorMessage
}

func (e *CustomError) StatusCode() int {
	return e.Status
}

func New(ctx context.Context, status int, err error, details ...any) *CustomError {
	data := logger.GetRequestContext(ctx)
	customErr := &CustomError{
		Status:       status,
		ErrorMessage: err.Error(),
		TraceID:      data.TraceID,
	}

	if len(details) > 0 {
		customErr.Details = details
	}

	return customErr
}

func HandleError(ctx context.Context, err error, traceID string) *CustomError {
	if err == nil {
		return nil
	}

	slog.DebugContext(ctx, err.Error())

	if jsonErr := handleJSONError(ctx, err); jsonErr != nil {
		return jsonErr
	}

	if validationErr := HandleValidationError(err); validationErr != nil {
		return validationErr
	}

	if dbError := handleDatabaseError(ctx, err); dbError != nil {
		return dbError
	}

	var customError *CustomError
	if errors.As(err, &customError) {
		customError.TraceID = traceID
		return customError
	}

	if config.GetInstance().Api.Environment != string(config.Production) {
		return New(ctx, http.StatusInternalServerError, err)
	}

	return New(ctx, http.StatusInternalServerError, errors.New("Internal Server Error"))
}
