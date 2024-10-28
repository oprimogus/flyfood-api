package xerrors

import (
	"errors"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/oprimogus/cardapiogo/internal/config"
)

func handleGocloakError(err error, transactionID string) *CustomError {
	var gocloakApiError *gocloak.APIError
	if !errors.As(err, &gocloakApiError) {
		return nil
	}

	messages := strings.Split(gocloakApiError.Message, ":")

	if config.GetInstance().Api.Environment != string(config.Production) {
		if len(messages) == 1 {
			return &CustomError{
				Status:        gocloakApiError.Code,
				ErrorMessage:  strings.TrimSpace(messages[len(messages)-1]),
				TransactionID: transactionID,
				Debug: err,
			}
		}

		if len(messages) == 2 {
			return &CustomError{
				Status:        gocloakApiError.Code,
				ErrorMessage:  strings.TrimSpace(messages[len(messages)-1]),
				TransactionID: transactionID,
				Debug: err,
			}
		}

		if len(messages) == 3 {
			return &CustomError{
				Status:        gocloakApiError.Code,
				ErrorMessage:  strings.TrimSpace(messages[len(messages)-1]),
				Details:       strings.TrimSpace(messages[len(messages)-2]),
				TransactionID: transactionID,
				Debug:         gocloakApiError,
			}
		}
	}

	return &CustomError{
		Status:        gocloakApiError.Code,
		ErrorMessage:  "An authentication error occurred while processing your request.",
		TransactionID: transactionID,
	}
}
