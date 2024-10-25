package core

import (
	"github.com/oprimogus/cardapiogo/internal/core/authentication"
	"github.com/oprimogus/cardapiogo/internal/core/item"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	"github.com/oprimogus/cardapiogo/internal/core/user"
)

type RepositoryFactory interface {
	NewUserRepository() user.Repository
	NewAuthenticationRepository() authentication.Repository
	NewStoreRepository() store.Repository
	NewItemRepository() item.Repository
}
