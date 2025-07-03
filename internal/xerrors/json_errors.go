package xerrors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func handleJSONError(ctx context.Context, err error) *CustomError {
	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return New(
			ctx,
			http.StatusBadRequest,
			fmt.Errorf(
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
			ctx,
			http.StatusBadRequest,
			fmt.Errorf("Invalid JSON: %s", jsonSyntaxError),
			jsonSyntaxError.Offset,
		)
	}

	return nil
}
