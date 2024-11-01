package item

import (
	"context"
	"fmt"

	"github.com/oprimogus/cardapiogo/internal/core/store"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

type UseCaseCreate struct {
	repository      Repository
	storeRepository store.Repository
}

func NewUseCaseCreate(repository Repository, storeRepository store.Repository) UseCaseCreate {
	return UseCaseCreate{repository: repository, storeRepository: storeRepository}
}

func (u UseCaseCreate) Execute(ctx context.Context, params CreateItemInput) (id int, err error) {
	userID, ok := ctx.Value(string(logger.UserIDKey)).(string)
	if !ok {
		return 0, fmt.Errorf("invalid userID: '%s'", userID)
	}
	isOwner, err := u.storeRepository.IsOwner(ctx, params.StoreID, userID)
	if err != nil {
		return 0, err
	}

	if !isOwner {
		return 0, store.ErrNotOwner
	}

	return u.repository.CreateItem(ctx, params, defaultScore)
}
