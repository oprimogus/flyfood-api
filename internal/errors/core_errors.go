package xerrors

import (
	"errors"
	"net/http"

	"github.com/oprimogus/cardapiogo/internal/core/store"
	"github.com/oprimogus/cardapiogo/internal/core/user"
)

func handleCoreError(err error, traceID string) *CustomError {
	type errorMapping struct {
		err    error
		status int
	}

	mappings := []errorMapping{
		{user.ErrExistUserWithDocument, http.StatusConflict},
		{user.ErrExistUserWithEmail, http.StatusConflict},
		{user.ErrExistUserWithPhone, http.StatusConflict},
		{store.ErrClosingTimeBeforeOpeningTime, http.StatusBadRequest},
		{store.ErrOpeningTimeAfterClosingTime, http.StatusBadRequest},
		{store.ErrNotOwner, http.StatusForbidden},
	}

	for _, mapping := range mappings {
		if errors.Is(err, mapping.err) {
			return &CustomError{
				Status:       mapping.status,
				ErrorMessage: err.Error(),
				TraceID:      traceID,
			}
		}
	}

	return nil
}
