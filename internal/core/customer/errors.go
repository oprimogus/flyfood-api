package customer

import (
	"context"
	"errors"

	"github.com/oprimogus/flyfood-api/pkg/xerrors"
)

var (
	errMaxAddresses = errors.New("você pode cadastrar até 5 endereços")
	errCustomerNotFound = errors.New("cliente não encontrado")
)

func ErrMaxAddresses(ctx context.Context) error {
	return xerrors.NewWithContext(ctx, errMaxAddresses).WithStatusBadRequest()
}

func ErrCustomerNotFound(ctx context.Context) error {
	return xerrors.NewWithContext(ctx, errCustomerNotFound).WithStatusNotFound()
}

