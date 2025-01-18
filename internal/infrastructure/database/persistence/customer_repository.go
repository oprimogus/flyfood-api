package persistence

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oprimogus/cardapiogo/internal/core/address"
	"github.com/oprimogus/cardapiogo/internal/core/customer"
	postgresDB "github.com/oprimogus/cardapiogo/internal/infrastructure/database/postgres"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/database/sqlc"
	"log/slog"
	"sync"
)

type CustomerRepository struct {
	db *pgxpool.Pool
	q  sqlc.Querier
}

func NewCustomerRepository(db *postgresDB.Database) CustomerRepository {
	return CustomerRepository{db: db.GetDB(), q: sqlc.New(db.GetDB())}
}

func (r CustomerRepository) FindByID(ctx context.Context, id string) (*customer.Customer, error) {
	var customerRepo sqlc.FindCustomerByIDRow
	var customerErr error
	var addressesRepo []sqlc.FindAddressesByCustomerIDRow
	var addressesErr error
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		customerRepo, customerErr = r.q.FindCustomerByID(ctx, id)
	}()

	go func() {
		defer wg.Done()
		addressesRepo, addressesErr = r.q.FindAddressesByCustomerID(ctx, id)
	}()

	wg.Wait()

	if customerErr != nil {
		return nil, customerErr
	}
	if addressesErr != nil {
		return nil, addressesErr
	}

	addresses := make([]address.Address, len(addressesRepo))
	for i, addr := range addressesRepo {
		addresses[i] = address.Address{
			Name:         addr.Name,
			AddressLine1: addr.AddressLine1,
			AddressLine2: addr.AddressLine2,
			Neighborhood: addr.Neighborhood,
			City:         addr.City,
			State:        addr.State,
			PostalCode:   addr.PostalCode,
			Country:      addr.Country,
		}
	}

	customerDomain := &customer.Customer{
		ID:       id,
		Name:     customerRepo.Name,
		LastName: customerRepo.LastName,
		CPF:      customerRepo.Cpf,
		Email:    customerRepo.Email,
		Phone:    customerRepo.Phone,
	}

	customerDomain.Addresses = addresses

	return customerDomain, nil
}

func (r CustomerRepository) Save(ctx context.Context, c *customer.Customer) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	qtx := sqlc.New(tx)

	if err := qtx.SaveCustomer(ctx, sqlc.SaveCustomerParams{
		ID:       c.ID,
		Name:     c.Name,
		LastName: c.LastName,
		Cpf:      c.CPF,
		Email:    c.Email,
		Phone:    c.Phone,
	}); err != nil {
		return err
	}

	if err := r.syncAddresses(ctx, qtx, c.ID, c.Addresses); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r CustomerRepository) syncAddresses(ctx context.Context,
	qtx *sqlc.Queries, customerID string, addresses []address.Address) error {

	addressesRepo, err := qtx.FindAddressesByCustomerID(ctx, customerID)
	if err != nil {
		return err
	}

	addressesToDelete := make([]sqlc.DeleteCustomerAddressesParams, len(addressesRepo))
	for i, v := range addressesRepo {
		addressesToDelete[i] = sqlc.DeleteCustomerAddressesParams{
			CustomerID:   customerID,
			Name:         v.Name,
			AddressLine1: v.AddressLine1,
			AddressLine2: v.AddressLine2,
			Neighborhood: v.Neighborhood,
			City:         v.City,
			State:        v.State,
			PostalCode:   v.PostalCode,
			Country:      v.Country,
		}
	}
	resultDelete := qtx.DeleteCustomerAddresses(ctx, addressesToDelete)
	resultDelete.Exec(func(i int, err error) {
		if err != nil {
			slog.Error(fmt.Sprintf("failed to delete address in index %d", i), "error", err.Error())
		}
	})
	if err := resultDelete.Close(); err != nil {
		return err
	}

	addressesToAdd := make([]sqlc.AddNewCustomerAddressesParams, len(addresses))
	for i, v := range addresses {
		addressesToAdd[i] = sqlc.AddNewCustomerAddressesParams{
			CustomerID:   customerID,
			Name:         v.Name,
			AddressLine1: v.AddressLine1,
			AddressLine2: v.AddressLine2,
			Neighborhood: v.Neighborhood,
			City:         v.City,
			State:        v.State,
			PostalCode:   v.PostalCode,
			Country:      v.Country,
		}
	}

	resultAdd := qtx.AddNewCustomerAddresses(ctx, addressesToAdd)
	resultAdd.Exec(func(i int, err error) {
		if err != nil {
			slog.Error(fmt.Sprintf("failed to add address in index %d", i), "error", err.Error())
		}
	})
	if err := resultAdd.Close(); err != nil {
		return err
	}
	return nil
}
