package xerrors

import (
	"errors"

	"log/slog"
	"net/http"

	"github.com/oprimogus/cardapiogo/internal/config"
)

const (
	NotFoundRecord      = "Record not found"
	DuplicatedRecord    = "There is a record with this data"
	InternalServerError = "Internal Server Error"
	TooManyValues       = "There is more than one record"
	InvalidValues       = "Invalid values for few fields"
	UnknownError        = "Unknown error"
)

type CustomError struct {
	Status       int         `json:"-"`
	ErrorMessage string      `json:"error"`
	Details      interface{} `json:"details,omitempty"`
	TraceID      string      `json:"traceID"`
	Debug        interface{} `json:"debug,omitempty"`
}

func New(traceID string, status int, message string, details ...interface{}) *CustomError {
	err := &CustomError{
		Status:       status,
		ErrorMessage: message,
		TraceID:      traceID,
	}

	if len(details) > 0 {
		err.Details = details[0]
	}

	return err
}

func (e *CustomError) Error() string {
	return e.ErrorMessage
}

func (e *CustomError) StatusCode() int {
	return e.Status
}

func HandleError(err error, traceID string) *CustomError {
	if err == nil {
		return nil
	}

	slog.Debug(err.Error())

	if jsonErr := handleJSONError(err, traceID); jsonErr != nil {
		return jsonErr
	}

	if gocloakErr := handleGocloakError(err, traceID); gocloakErr != nil {
		return gocloakErr
	}

	if dbErr := handleDatabaseError(err, traceID); dbErr != nil {
		return dbErr
	}

	if coreErr := handleCoreError(err, traceID); coreErr != nil {
		return coreErr
	}

	var customError *CustomError
	if errors.As(err, &customError) {
		customError.TraceID = traceID
		return customError
	}

	if config.GetInstance().Api.Environment != string(config.Production) {
		return New(traceID, http.StatusInternalServerError, err.Error(), err)
	}

	return New(traceID, http.StatusInternalServerError, InternalServerError)
}
