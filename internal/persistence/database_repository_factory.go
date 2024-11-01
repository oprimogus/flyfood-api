package persistence

import (
	"github.com/oprimogus/cardapiogo/internal/core"
	"github.com/oprimogus/cardapiogo/internal/core/authentication"
	"github.com/oprimogus/cardapiogo/internal/core/item"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	"github.com/oprimogus/cardapiogo/internal/core/user"
	"github.com/oprimogus/cardapiogo/internal/database/postgres"
	"github.com/oprimogus/cardapiogo/internal/database/sqlc"
	"github.com/oprimogus/cardapiogo/internal/services/adapter"
)

type DatabaseRepositoryFactory struct {
	db             *postgres.Database
	querier        *sqlc.Queries
	serviceFactory adapter.Factory
}

func NewDataBaseRepositoryFactory(db *postgres.Database, adapter adapter.Factory) core.RepositoryFactory {
	return &DatabaseRepositoryFactory{db: db, querier: sqlc.New(db.GetDB()), serviceFactory: adapter}
}

func (d *DatabaseRepositoryFactory) NewUserRepository() user.Repository {
	return NewUserRepository()
}

func (d *DatabaseRepositoryFactory) NewAuthenticationRepository() authentication.Repository {
	return NewAuthenticationRepository()
}

func (d *DatabaseRepositoryFactory) NewStoreRepository() store.Repository {
	return NewStoreRepository(d.db, d.querier, d.serviceFactory.NewStorageService())
}

func (d *DatabaseRepositoryFactory) NewItemRepository() item.Repository {
	return NewItemRepository(d.db, d.querier)
}
