package xerrors

import "net/http"

type FieldError struct {
	Field   string      `json:"field"`
	Input   string      `json:"input"`
	Message string      `json:"message"`
	Debug   interface{} `json:"debug,omitempty"`
}

type InvalidField struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

const INVALID_INPUT_MESSAGE = "There are some problems with the data you submitted"

func InvalidInput(transactionID string, errs map[string]string) *CustomError {
	details := make([]InvalidField, 0, len(errs))
	for field, msg := range errs {
		details = append(details, InvalidField{
			Field: field,
			Error: msg,
		})
	}

	return &CustomError{
		Status:        http.StatusBadRequest,
		ErrorMessage:  INVALID_INPUT_MESSAGE,
		Details:       details,
		TransactionID: transactionID,
	}
}
