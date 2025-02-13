package xerrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func handleJSONError(err error, traceID string) *CustomError {
	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return New(
			traceID,
			http.StatusBadRequest,
			fmt.Sprintf(
				"Invalid JSON: field %s is not valid for type %s",
				unmarshalTypeError.Field,
				unmarshalTypeError.Value,
			),
			unmarshalTypeError.Struct,
		)
	}

	var jsonSyntaxError *json.SyntaxError
	if errors.As(err, &jsonSyntaxError) {
		return New(
			traceID,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid JSON: %s", jsonSyntaxError),
			jsonSyntaxError.Offset,
		)
	}

	return nil
}
