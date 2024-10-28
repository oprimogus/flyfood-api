package xerrors

import (
	"errors"

	"log/slog"
	"net/http"

	"github.com/oprimogus/cardapiogo/internal/config"
)

const (
	NOT_FOUND_RECORD      = "Record not found"
	DUPLICATED_RECORD     = "There is a record with this data"
	FOREIGN_KEY_VIOLATION = "Foreign key violation"
	NULL_VIOLATION        = "Null value not allowed for column"
	VALUE_TOO_LONG        = "Input value too long for column"
	INTERNAL_SERVER_ERROR = "Internal Server Error"
	TOO_MANY_VALUES       = "There is more than one record"
	INVALID_VALUES        = "Invalid values for few fields"
	UNKNOWN_ERROR         = "Unknown error"
)

type CustomError struct {
	Status        int         `json:"-"`
	ErrorMessage  string      `json:"error"`
	Details       interface{} `json:"details,omitempty"`
	TransactionID string      `json:"transactionID"`
	Debug         interface{} `json:"debug,omitempty"`
}

func New(transactionID string, status int, message string, details ...interface{}) *CustomError {
	err := &CustomError{
		Status:        status,
		ErrorMessage:  message,
		TransactionID: transactionID,
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

func HandleError(err error, transactionID string) *CustomError {
	if err == nil {
		return nil
	}

	slog.Debug(err.Error())

	if jsonErr := handleJSONError(err, transactionID); jsonErr != nil {
		return jsonErr
	}

	if gocloakErr := handleGocloakError(err, transactionID); gocloakErr != nil {
		return gocloakErr
	}

	if dbErr := handleDatabaseError(err, transactionID); dbErr != nil {
		return dbErr
	}

	if coreErr := handleCoreError(err, transactionID); coreErr != nil {
		return coreErr
	}

	var customError *CustomError
	if errors.As(err, &customError) {
		customError.TransactionID = transactionID
		return customError
	}

	if config.GetInstance().Api.Environment != string(config.Production) {
		return New(transactionID, http.StatusInternalServerError, err.Error(), err)
	}

	return New(transactionID, http.StatusInternalServerError, INTERNAL_SERVER_ERROR)
}
