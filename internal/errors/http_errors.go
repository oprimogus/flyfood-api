package xerrors

import (
	"net/http"
)

func InternalServerError(transactionID string, msg string) *CustomError {
	if msg == "" {
		msg = "We encountered an error while processing your request."
	}
	return New(transactionID, http.StatusInternalServerError, msg)
}

func ConflictError(transactionID string, msg string) *CustomError {
	if msg == "" {
		msg = "We encountered an conflict error while processing your request."
	}
	return New(transactionID, http.StatusConflict, msg)
}

func NotFound(transactionID string, msg string) *CustomError {
	if msg == "" {
		msg = "The requested resource was not found."
	}
	return New(transactionID, http.StatusNotFound, msg)
}

func Unauthorized(transactionID string, msg string) *CustomError {
	if msg == "" {
		msg = "You are not authenticated to perform the requested action."
	}
	return New(transactionID, http.StatusUnauthorized, msg)
}

func Forbidden(transactionID string, msg string) *CustomError {
	if msg == "" {
		msg = "You are not authorized to perform the requested action."
	}
	return New(transactionID, http.StatusForbidden, msg)
}

func BadRequest(transactionID string, msg string) *CustomError {
	if msg == "" {
		msg = "Your request is in a bad format."
	}
	return New(transactionID, http.StatusBadRequest, msg)
}

type invalidField struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func InvalidInput(transactionID string, errs map[string]string) *CustomError {
	details := []invalidField{}
	for i, v := range errs {
		details = append(details, invalidField{
			Field: i,
			Error: v,
		})
	}
	return &CustomError{
		Status:        http.StatusBadRequest,
		ErrorMessage:  "There is some problem with the data you submitted.",
		Details:       details,
		TransactionID: transactionID,
	}
}
