package persistence

import (
	postgresDB "github.com/oprimogus/cardapiogo/internal/infrastructure/database/postgres"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/database/sqlc"
)

type RepositoryFactory struct {
	db *postgresDB.Database
	q  sqlc.Querier
}

func NewRepositoryFactory(db *postgresDB.Database) RepositoryFactory {
	return RepositoryFactory{db: db, q: sqlc.New(db.GetDB())}
}

func (f RepositoryFactory) NewSQLC() sqlc.Querier {
	return f.q
}

func (f RepositoryFactory) NewCustomerRepository() CustomerRepository {
	return NewCustomerRepository(f.db)
}

func (f RepositoryFactory) NewStoreRepository() StoreRepository {
	return NewStoreRepository(f.db)
}

func (f RepositoryFactory) NewProductRepository() ProductRepository {
	return NewProductRepository(f.db)
}

func (f RepositoryFactory) NewOwnerRepository() OwnerRepository {
	return NewOwnerRepository(f.db)
}
