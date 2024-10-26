package xerrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v13"

	"github.com/oprimogus/cardapiogo/internal/config"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	"github.com/oprimogus/cardapiogo/internal/core/user"
)

var (
	environment = config.GetInstance().Api.Environment
)

type CustomError struct {
	Status        int         `json:"-"`
	ErrorMessage  string      `json:"error"`
	Details       interface{} `json:"details"`
	TransactionID string      `json:"transactionID"`
	Debug         interface{} `json:"debug,omitempty"`
}

func New(transactionID string, status int, message string, details ...interface{}) *CustomError {
	return &CustomError{
		Status:        status,
		ErrorMessage:  message,
		Details:       details,
		TransactionID: transactionID,
	}
}

func (e *CustomError) Error() string {
	return e.ErrorMessage
}

func (e *CustomError) StatusCode() int {
	return e.Status
}

func HandleError(err error, transactionID string) *CustomError {
	slog.Debug(err.Error())

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		er := err.(*json.UnmarshalTypeError)
		return &CustomError{
			Status: http.StatusBadRequest,
			ErrorMessage: fmt.Sprintf(
				"Invalid JSON: field %s is not valid for type %s",
				er.Field,
				er.Value),
			Details:       er.Struct,
			TransactionID: transactionID,
		}
	}

	var jsonSyntaxError *json.SyntaxError
	if errors.As(err, &jsonSyntaxError) {
		er := err.(*json.SyntaxError)
		return &CustomError{
			Status:        http.StatusBadRequest,
			ErrorMessage:  fmt.Sprintf("Invalid JSON: %s", er),
			Details:       er.Offset,
			TransactionID: transactionID,
		}
	}

	var gocloakApiError *gocloak.APIError
	if errors.As(err, &gocloakApiError) {
		er := err.(*gocloak.APIError)
		messages := strings.Split(er.Message, ":")
		if environment != string(config.Production) {
			if len(messages) == 2 {
				return &CustomError{
					Status:        er.Code,
					ErrorMessage:  strings.TrimSpace(messages[len(messages)-1]),
					TransactionID: transactionID,
				}
			}
			return &CustomError{
				Status:        er.Code,
				ErrorMessage:  strings.TrimSpace(messages[len(messages)-1]),
				Details:       strings.TrimSpace(messages[len(messages)-2]),
				TransactionID: transactionID,
			}
		}
		return &CustomError{
			Status:        er.Code,
			ErrorMessage:  "Occurred an error when you request was processed.",
			TransactionID: transactionID,
		}
	}

	handledCoreError := handleCoreError(err, transactionID)
	if handledCoreError != nil {
		return handledCoreError
	}

	var customError *CustomError
	if errors.As(err, &customError) {
		customError.TransactionID = transactionID
		return customError
	}

	if isDatabaseError(err) {
		return handleDatabaseErrors(err, transactionID)
	}

	if environment != string(config.Production) {
		return &CustomError{
			Status:        http.StatusInternalServerError,
			ErrorMessage:  err.Error(),
			Details:       err,
			TransactionID: transactionID,
		}
	}

	return &CustomError{
		Status:        http.StatusInternalServerError,
		ErrorMessage:  INTERNAL_SERVER_ERROR,
		TransactionID: transactionID,
	}

}

func handleCoreError(err error, transactionID string) *CustomError {
	switch err {
	case user.ErrExistUserWithDocument,
		user.ErrExistUserWithEmail,
		user.ErrExistUserWithPhone:
		return &CustomError{
			Status:        http.StatusConflict,
			ErrorMessage:  err.Error(),
			TransactionID: transactionID,
		}
	case store.ErrClosingTimeBeforeOpeningTime,
		store.ErrOpeningTimeAfterClosingTime:
		return &CustomError{
			Status:        http.StatusBadRequest,
			ErrorMessage:  err.Error(),
			TransactionID: transactionID,
		}
	default:
		return nil
	}

}
