package item

import "github.com/oprimogus/cardapiogo/internal/core/store"

type ItemModule struct {
	Create      useCaseCreate
	GetByID     useCaseGetByID
	GetByFilter useCaseGetByFilter
}

func NewItemModule(repository Repository, storeRepository store.Repository) ItemModule {
	return ItemModule{
		Create:      newUseCaseCreate(repository, storeRepository),
		GetByID:     newUseCaseGetByID(repository),
		GetByFilter: newUseCaseGetByFilter(repository),
	}
}
