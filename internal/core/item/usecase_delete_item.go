package item

import (
	"context"
	"fmt"

	"github.com/oprimogus/cardapiogo/internal/core/store"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

type useCaseDelete struct {
	repository      Repository
	storeRepository store.Repository
}

func NewUseCaseDelete(repository Repository, storeRepository store.Repository) *useCaseDelete {
	return &useCaseDelete{repository: repository, storeRepository: storeRepository}
}

func (u *useCaseDelete) Execute(ctx context.Context, storeID string, itemID int) error {
	userID, ok := ctx.Value(string(logger.UserIDKey)).(string)
	if !ok {
		return fmt.Errorf("invalid userID: '%s'", userID)
	}
	isOwner, err := u.storeRepository.IsOwner(ctx, storeID, userID)
	if err != nil {
		return err
	}

	if !isOwner {
		return store.ErrNotOwner
	}

	return u.repository.DeleteItem(ctx, itemID)
}
