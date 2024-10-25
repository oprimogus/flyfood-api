package xerrors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/oprimogus/cardapiogo/internal/config"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	"github.com/oprimogus/cardapiogo/internal/core/user"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

var (
	log         = logger.NewLogger("Error Handling")
	environment = config.GetInstance().Api.Environment
)

type ErrorResponse struct {
	Status        int         `json:"-"`
	ErrorMessage  string      `json:"error"`
	Details       interface{} `json:"details"`
	TransactionID string      `json:"transactionID"`
	Debug         interface{} `json:"debug,omitempty"`
}

func New(status int, message string, transactionID string, details ...interface{}) *ErrorResponse {
	return &ErrorResponse{
		Status:        status,
		ErrorMessage:  message,
		Details:       details,
		TransactionID: transactionID,
	}
}

func (e *ErrorResponse) Error() string {
	return e.ErrorMessage
}

func (e *ErrorResponse) StatusCode() int {
	return e.Status
}

func HandleError(err error, transactionID string) *ErrorResponse {
	log.Debug(err.Error())
	if errResp, ok := err.(*ErrorResponse); ok {
		errResp.TransactionID = transactionID
		return errResp
	}
	if errResp, ok := err.(*json.UnmarshalTypeError); ok {
		return &ErrorResponse{
			Status:        http.StatusBadRequest,
			ErrorMessage:  fmt.Sprintf("Invalid JSON: field %s is not valid for type %s", errResp.Field, errResp.Value),
			Details:       errResp.Struct,
			TransactionID: transactionID,
		}
	}
	if errResp, ok := err.(*gocloak.APIError); ok {
		messages := strings.Split(errResp.Message, ":")
		if environment != string(config.Production) {
			return &ErrorResponse{
				Status:        errResp.Code,
				ErrorMessage:  strings.TrimSpace(messages[len(messages)-1]),
				Details:       strings.TrimSpace(messages[len(messages)-2]),
				TransactionID: transactionID,
			}
		}
		return &ErrorResponse{
			Status:        errResp.Code,
			ErrorMessage:  "Occurred an error when you request was processed.",
			TransactionID: transactionID,
		}
	}

	if errResp, ok := err.(*pgconn.PgError); ok {
		return handleDatabaseErrors(errResp, transactionID)
	}

	return handleCoreError(err, transactionID)

}

func handleCoreError(err error, transactionID string) *ErrorResponse {
	switch err {
	case pgx.ErrNoRows:
		return New(http.StatusNotFound, NOT_FOUND_RECORD, transactionID)
	case pgx.ErrTooManyRows:
		return New(http.StatusInternalServerError, TOO_MANY_VALUES, transactionID)
	case user.ErrExistUserWithDocument,
		user.ErrExistUserWithEmail,
		user.ErrExistUserWithPhone:
		return &ErrorResponse{
			Status:        http.StatusConflict,
			ErrorMessage:  err.Error(),
			TransactionID: transactionID,
		}
	case store.ErrClosingTimeBeforeOpeningTime,
		store.ErrOpeningTimeAfterClosingTime:
		return &ErrorResponse{
			Status:        http.StatusBadRequest,
			ErrorMessage:  err.Error(),
			TransactionID: transactionID,
		}
	default:
		return &ErrorResponse{
			Status:        http.StatusInternalServerError,
			ErrorMessage:  err.Error(),
			TransactionID: transactionID,
		}
	}

}
