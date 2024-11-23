package xerrors

import "net/http"

func InternalServer(traceID string, msg string) *CustomError {
	if msg == "" {
		msg = "We encountered an error while processing your request."
	}
	return New(traceID, http.StatusInternalServerError, msg)
}

func Conflict(traceID string, msg string) *CustomError {
	if msg == "" {
		msg = "We encountered a conflict error while processing your request."
	}
	return New(traceID, http.StatusConflict, msg)
}

func NotFound(traceID string, msg string) *CustomError {
	if msg == "" {
		msg = "The requested resource was not found."
	}
	return New(traceID, http.StatusNotFound, msg)
}

func Unauthorized(traceID string, msg string) *CustomError {
	if msg == "" {
		msg = "You are not authenticated to perform the requested action."
	}
	return New(traceID, http.StatusUnauthorized, msg)
}

func Forbidden(traceID string, msg string) *CustomError {
	if msg == "" {
		msg = "You are not authorized to perform the requested action."
	}
	return New(traceID, http.StatusForbidden, msg)
}

func BadRequest(traceID string, msg string) *CustomError {
	if msg == "" {
		msg = "Your request is in a bad format."
	}
	return New(traceID, http.StatusBadRequest, msg)
}
