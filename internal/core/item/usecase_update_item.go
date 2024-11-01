package item

import (
	"context"
	"fmt"

	"github.com/oprimogus/cardapiogo/internal/core/store"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

type UseCaseUpdate struct {
	repository      Repository
	storeRepository store.Repository
}

func NewUseCaseUpdate(repository Repository, storeRepository store.Repository) UseCaseUpdate {
	return UseCaseUpdate{repository: repository, storeRepository: storeRepository}
}

func (u UseCaseUpdate) Execute(ctx context.Context, params UpdateItemInput) error {
	userID, ok := ctx.Value(string(logger.UserIDKey)).(string)
	if !ok {
		return fmt.Errorf("invalid userID: '%s'", userID)
	}
	isOwner, err := u.storeRepository.IsOwner(ctx, params.StoreID, userID)
	if err != nil {
		return err
	}

	if !isOwner {
		return store.ErrNotOwner
	}

	return u.repository.UpdateItem(ctx, params)
}
