package xerrors

import (
	"context"
	"errors"
	"net/http"
)

func InternalServer(ctx context.Context, msg string) *CustomError {
	if msg == "" {
		msg = "We encountered an error while processing your request."
	}
	return New(ctx, http.StatusInternalServerError, errors.New(msg))
}

func Conflict(ctx context.Context, msg string) *CustomError {
	if msg == "" {
		msg = "We encountered a conflict error while processing your request."
	}
	return New(ctx, http.StatusConflict, errors.New(msg))
}

func NotFound(ctx context.Context, msg string) *CustomError {
	if msg == "" {
		msg = "The requested resource was not found."
	}
	return New(ctx, http.StatusNotFound, errors.New(msg))
}

func Unauthorized(ctx context.Context, msg string) *CustomError {
	if msg == "" {
		msg = "You are not authenticated to perform the requested action."
	}
	return New(ctx, http.StatusUnauthorized, errors.New(msg))
}

func Forbidden(ctx context.Context, msg string) *CustomError {
	if msg == "" {
		msg = "You are not authorized to perform the requested action."
	}
	return New(ctx, http.StatusForbidden, errors.New(msg))
}

func BadRequest(ctx context.Context, msg string) *CustomError {
	if msg == "" {
		msg = "Your request is in a bad format."
	}
	return New(ctx, http.StatusBadRequest, errors.New(msg))
}
