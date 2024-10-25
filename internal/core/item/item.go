package item

import "github.com/oprimogus/cardapiogo/internal/core/store"

type ItemModule struct {
	Create  useCaseCreate
	GetByID useCaseGetByID
}

func NewItemModule(repository Repository, storeRepository store.Repository) ItemModule {
	return ItemModule{
		Create:  newUseCaseCreate(repository, storeRepository),
		GetByID: *NewUseCaseGetByID(repository),
	}
}
