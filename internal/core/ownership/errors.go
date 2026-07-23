package ownership

import (
	"context"
	"errors"

	"github.com/oprimogus/flyfood-api/pkg/xerrors"
)

var (
	errNotOwnerOfResource = errors.New("você não é dono deste recurso")
	errNotAnOwner         = errors.New("você não é um dono de negócio")
)

func ErrNotOwnerOfResource(ctx context.Context) error {
	return xerrors.NewWithContext(ctx, errNotOwnerOfResource).WithStatusForbidden()
}

func ErrNotAnOwner(ctx context.Context) error {
	return xerrors.NewWithContext(ctx, errNotAnOwner).WithStatusForbidden()
}
