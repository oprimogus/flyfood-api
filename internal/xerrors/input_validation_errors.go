package xerrors

import (
	"errors"
	"github.com/oprimogus/cardapiogo/internal/xvalidator"
	"net/http"
)

func HandleValidationError(err error) *CustomError {
	var validationError *xvalidator.ValidationError
	if errors.As(err, &validationError) {
		return &CustomError{
			Status:       http.StatusBadRequest,
			ErrorMessage: validationError.Message,
			Details:      validationError.Fields,
		}
	}
	return nil
}
