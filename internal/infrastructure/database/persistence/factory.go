package persistence

import postgresDB "github.com/oprimogus/cardapiogo/internal/infrastructure/database/postgres"

type RepositoryFactory struct {
	db *postgresDB.Database
}

func NewRepositoryFactory(db *postgresDB.Database) RepositoryFactory {
	return RepositoryFactory{db}
}

func (f RepositoryFactory) NewCustomerRepository() CustomerRepository {
	return NewCustomerRepository(f.db)
}

func (f RepositoryFactory) NewStoreRepository() StoreRepository {
	return NewStoreRepository(f.db)
}
