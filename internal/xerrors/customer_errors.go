package xerrors

import (
	"errors"
	"github.com/oprimogus/flyfood-api/internal/core/customer"
	"net/http"
)

type errorMapping struct {
	err    error
	status int
}

func HandleCustomerError(err error) *CustomError {

	mappings := []errorMapping{
		{customer.ErrMaxAddresses, http.StatusUnprocessableEntity},
		{customer.ErrTryRemoveInvalidAddress, http.StatusUnprocessableEntity},
		{customer.ErrThereIsNoAddresses, http.StatusUnprocessableEntity},
	}

	for _, mapping := range mappings {
		if errors.Is(err, mapping.err) {
			return &CustomError{
				Status:       mapping.status,
				ErrorMessage: err.Error(),
				TraceID:      "",
			}
		}
	}

	return nil
}
