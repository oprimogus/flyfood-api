package xerrors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/oprimogus/cardapiogo/internal/config"
)

var ErrInvalidCredentials = fmt.Errorf("invalid email or password")

func handleGocloakError(err error, traceID string) *CustomError {
	var gocloakApiError *gocloak.APIError
	if !errors.As(err, &gocloakApiError) {
		return nil
	}

	messages := strings.Split(gocloakApiError.Message, ":")

	if config.GetInstance().Api.Environment != string(config.Production) {
		if len(messages) == 1 {
			return &CustomError{
				Status:       gocloakApiError.Code,
				ErrorMessage: strings.TrimSpace(messages[len(messages)-1]),
				TraceID:      traceID,
				Debug:        err,
			}
		}

		if len(messages) == 2 {
			return &CustomError{
				Status:       gocloakApiError.Code,
				ErrorMessage: strings.TrimSpace(messages[len(messages)-1]),
				TraceID:      traceID,
				Debug:        err,
			}
		}

		if len(messages) == 3 {
			return &CustomError{
				Status:       gocloakApiError.Code,
				ErrorMessage: strings.TrimSpace(messages[0]),
				Details:      strings.TrimSpace(messages[len(messages)-2]),
				TraceID:      traceID,
				Debug:        gocloakApiError,
			}
		}
	}

	if gocloakApiError.Code == 401 && gocloakApiError.Type == "unknown" {
		return &CustomError{
			Status:       gocloakApiError.Code,
			ErrorMessage: ErrInvalidCredentials.Error(),
			TraceID:      traceID,
		}
	}

	return &CustomError{
		Status:       gocloakApiError.Code,
		ErrorMessage: "An authentication error occurred while processing your request.",
		TraceID:      traceID,
	}
}
