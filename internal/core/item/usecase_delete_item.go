package item

import (
	"context"
	"fmt"

	"github.com/oprimogus/cardapiogo/internal/core/store"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

type UseCaseDelete struct {
	repository      Repository
	storeRepository store.Repository
}

func NewUseCaseDelete(repository Repository, storeRepository store.Repository) UseCaseDelete {
	return UseCaseDelete{repository: repository, storeRepository: storeRepository}
}

func (u UseCaseDelete) Execute(ctx context.Context, id int) error {
	userID, ok := ctx.Value(string(logger.UserIDKey)).(string)
	if !ok {
		return fmt.Errorf("invalid userID: '%s'", userID)
	}

	item, err := u.repository.GetItemByID(ctx, id)
	if err != nil {
		return err
	}

	isOwner, err := u.storeRepository.IsOwner(ctx, item.StoreID, userID)
	if err != nil {
		return err
	}

	if !isOwner {
		return store.ErrNotOwner
	}

	return u.repository.DeleteItem(ctx, id)
}
