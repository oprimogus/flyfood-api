package ownership

import (
	"context"
	"errors"

	"github.com/oprimogus/flyfood-api/pkg/xerrors"
)

var (
    NotOwnerOfResourceErr = errors.New("Você não é dono deste recurso")
    NotOwnerErr = errors.New("Você não é um dono de negócio")
)

func ErrNotOwnerOfResource(ctx context.Context) error {
	return xerrors.NewWithContext(ctx, NotOwnerOfResourceErr).WithStatusForbidden()
}

func ErrNotAnOwner(ctx context.Context) error {
	return xerrors.NewWithContext(ctx, NotOwnerErr).WithStatusForbidden()
}
