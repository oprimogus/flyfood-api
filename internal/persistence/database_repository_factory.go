package persistence

import (
	"github.com/oprimogus/cardapiogo/internal/core"
	"github.com/oprimogus/cardapiogo/internal/core/authentication"
	"github.com/oprimogus/cardapiogo/internal/core/item"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	"github.com/oprimogus/cardapiogo/internal/core/user"
	"github.com/oprimogus/cardapiogo/internal/database/postgres"
	"github.com/oprimogus/cardapiogo/internal/database/sqlc"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

var (
	log = logger.NewLogger("RepositoryFactory")
)

type DatabaseRepositoryFactory struct {
	db      *postgres.PostgresDatabase
	querier *sqlc.Queries
}

func NewDataBaseRepositoryFactory(db *postgres.PostgresDatabase) core.RepositoryFactory {
	return &DatabaseRepositoryFactory{db: db, querier: sqlc.New(db.GetDB())}
}

func (d *DatabaseRepositoryFactory) NewUserRepository() user.Repository {
	return NewUserRepository()
}

func (d *DatabaseRepositoryFactory) NewAuthenticationRepository() authentication.Repository {
	return NewAuthenticationRepository()
}

func (d *DatabaseRepositoryFactory) NewStoreRepository() store.Repository {
	return NewStoreRepository(d.db, d.querier)
}

func (d *DatabaseRepositoryFactory) NewItemRepository() item.Repository {
	return NewItemRepository(d.db, d.querier)
}
