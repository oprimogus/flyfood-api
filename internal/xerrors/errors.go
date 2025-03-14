package xerrors

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/oprimogus/flyfood-api/internal/config"
)

type CustomError struct {
	Status       int         `json:"-"`
	ErrorMessage string      `json:"error"`
	Details      interface{} `json:"details,omitempty"`
	TraceID      string      `json:"traceID"`
	Debug        interface{} `json:"debug,omitempty"`
}

type FieldError struct {
	Field   string      `json:"field"`
	Input   string      `json:"input"`
	Message string      `json:"message"`
	Debug   interface{} `json:"debug,omitempty"`
}

func (e *CustomError) Error() string {
	return e.ErrorMessage
}

func (e *CustomError) StatusCode() int {
	return e.Status
}

func New(traceID string, status int, message string, details ...interface{}) *CustomError {
	err := &CustomError{
		Status:       status,
		ErrorMessage: message,
		TraceID:      traceID,
	}

	if len(details) > 0 {
		err.Details = details
	}

	return err
}

func HandleError(ctx context.Context, err error, traceID string) *CustomError {
	if err == nil {
		return nil
	}

	slog.DebugContext(ctx, err.Error())

	if jsonErr := handleJSONError(err, traceID); jsonErr != nil {
		return jsonErr
	}

	if validationErr := HandleValidationError(err); validationErr != nil {
		return validationErr
	}

	if gocloakErr := handleGocloakError(err, traceID); gocloakErr != nil {
		return gocloakErr
	}

	if customerErr := HandleCustomerError(err); customerErr != nil {
		return customerErr
	}

	if dbError := handleDatabaseError(err, traceID); dbError != nil {
		return dbError
	}

	var customError *CustomError
	if errors.As(err, &customError) {
		customError.TraceID = traceID
		return customError
	}

	if config.GetInstance().Api.Environment != string(config.Production) {
		return New(traceID, http.StatusInternalServerError, err.Error())
	}

	return New(traceID, http.StatusInternalServerError, "Internal Server Error")
}
