package customer

import (
	"context"
	"errors"
	"net/http"

	"github.com/oprimogus/flyfood-api/internal/xerrors"
)

var (
	ErrMaxAddresses            = errors.New("você pode cadastrar até 5 endereços")
	ErrTryRemoveInvalidAddress = errors.New("o endereço não pertence a essa conta")
	ErrThereIsNoAddresses      = errors.New("você não possui este endereço cadastrado")
	ErrCustomerAlreadyExist    = errors.New("dados já cadastrados")
)

var errStatusMap = map[error]int{
	ErrMaxAddresses:            http.StatusUnprocessableEntity,
	ErrTryRemoveInvalidAddress: http.StatusUnprocessableEntity,
	ErrThereIsNoAddresses:      http.StatusUnprocessableEntity,
	ErrCustomerAlreadyExist:    http.StatusUnprocessableEntity,
}

func HandleError(ctx context.Context, err error) *xerrors.CustomError {
	for domainErr, status := range errStatusMap {
		if errors.Is(err, domainErr) {
			return xerrors.New(ctx, status, domainErr)
		}
	}

	return xerrors.New(
		ctx,
		http.StatusInternalServerError,
		err,
	)
}
