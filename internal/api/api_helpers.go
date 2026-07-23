package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/oprimogus/flyfood-api/pkg/xerrors"
)

func HandleApiError(w http.ResponseWriter, r *http.Request, err error) {
	var xerr *xerrors.CustomError
	if errors.As(err, &xerr) {
		xerr = xerr.WithContext(r.Context())
		w.WriteHeader(xerr.Status)
		_ = json.NewEncoder(w).Encode(xerr)
		return
	}
	xerr = xerrors.New(err.Error()).
		WithContext(r.Context()).
		WithDetails(err)

	w.WriteHeader(xerr.Status)
	_ = json.NewEncoder(w).Encode(xerr)
}

func JSONResponse(w http.ResponseWriter, status int, response any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func IsValidBool(queryParam string) (bool, error) {
	return strconv.ParseBool(queryParam)
}

func IsValidInt(queryParam string) (int, error) {
	return strconv.Atoi(queryParam)
}