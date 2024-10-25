package item

import (
	"context"
	"fmt"

	"github.com/oprimogus/cardapiogo/internal/core/store"
)

type useCaseCreate struct {
	repository      Repository
	storeRepository store.Repository
}

func newUseCaseCreate(repository Repository, storeRepository store.Repository) useCaseCreate {
	return useCaseCreate{repository: repository, storeRepository: storeRepository}
}

func (u useCaseCreate) Execute(ctx context.Context, params CreateItemInput) (id int, err error) {
	userID := ctx.Value("userID").(string)
	if userID == "" {
		return 0, fmt.Errorf("invalid userID: %s", userID)
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
