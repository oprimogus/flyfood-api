package customer

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oprimogus/flyfood-api/internal/core/address"
	"github.com/oprimogus/flyfood-api/internal/infra/database"
	"github.com/oprimogus/flyfood-api/internal/infra/database/sqlc"
	"github.com/oprimogus/flyfood-api/pkg/xerrors"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Customer, error)
	FindByExternalID(ctx context.Context, externalId string) (*Customer, error)
	Save(ctx context.Context, c *Customer) error

	FindAddressesByExternalCustomerID(ctx context.Context, externalId string) ([]address.Address, error)
	SaveAddress(ctx context.Context, customerId uuid.UUID, addr address.Address) error
	DeleteAddress(ctx context.Context, customerId uuid.UUID, addressId uuid.UUID) error
}

type repository struct {
	db *pgxpool.Pool
	q  sqlc.Querier
}

func NewRepository(db *database.Postgres) Repository {
	return &repository{db: db.Pool, q: sqlc.New(db.Pool)}
}

func (r repository) FindByID(ctx context.Context, id uuid.UUID) (*Customer, error) {
	c, err := r.q.FindCustomerByID(ctx, database.ToPgTypeUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCustomerNotFound(ctx)
		}
		return nil, err
	}

	cDomain := &Customer{
		ID:         database.ToUUID(c.ID),
		ExternalID: c.ExternalID,
		Name:       c.Name,
		LastName:   c.LastName,
		CPF:        c.Cpf.String,
		Email:      c.Email,
		Phone:      c.Phone,
	}
	return cDomain, nil
}

func (r repository) FindByExternalID(ctx context.Context, externalId string) (*Customer, error) {
	c, err := r.q.FindCustomerByExternalID(ctx, externalId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCustomerNotFound(ctx)
		}
		return nil, err
	}

	cDomain := &Customer{
		ID:         database.ToUUID(c.ID),
		ExternalID: c.ExternalID,
		Name:       c.Name,
		LastName:   c.LastName,
		CPF:        c.Cpf.String,
		Email:      c.Email,
		Phone:      c.Phone,
	}
	return cDomain, nil
}

func (r repository) Save(ctx context.Context, c *Customer) error {
	if err := r.q.SaveCustomer(ctx, sqlc.SaveCustomerParams{
		ID:         database.ToPgTypeUUID(c.ID),
		ExternalID: c.ExternalID,
		Name:       c.Name,
		LastName:   c.LastName,
		Cpf: pgtype.Text{
			String: c.CPF,
			Valid:  c.CPF != "",
		},
		Email: c.Email,
		Phone: c.Phone,
	}); err != nil {
		return err
	}
	return nil
}

func (r repository) FindAddressesByExternalCustomerID(ctx context.Context, externalId string) ([]address.Address, error) {
	addressesSqlc, err := r.q.FindAddressesByExternalCustomerID(ctx, externalId)
	if err != nil {
		return nil, err
	}

	addrs := make([]address.Address, len(addressesSqlc))
	for i, a := range addressesSqlc {
		addrs[i] = address.Address{
			ID:           database.ToUUID(a.ID),
			Name:         a.Name,
			Default:      a.IsDefault.Bool,
			AddressLine1: a.AddressLine1,
			AddressLine2: a.AddressLine2,
			Neighborhood: a.Neighborhood,
			City:         a.City,
			State:        a.State,
			PostalCode:   a.PostalCode,
			Country:      a.Country,
			Latitude:     a.Latitude,
			Longitude:    a.Longitude,
		}
	}
	return addrs, nil
}

func (r repository) SaveAddress(ctx context.Context, customerId uuid.UUID, addr address.Address) error {
    transaction, err := r.db.Begin(ctx)
    if err != nil {
        return xerrors.NewWithContext(ctx, err)
    }
    defer transaction.Rollback(ctx)
    
    tq := sqlc.New(transaction)

	if err := tq.SaveAddress(ctx, sqlc.SaveAddressParams{
		ID:           database.ToPgTypeUUID(addr.ID),
		Name:         addr.Name,
		AddressLine1: addr.AddressLine1,
		AddressLine2: addr.AddressLine2,
		Neighborhood: addr.Neighborhood,
		City:         addr.City,
		State:        addr.State,
		PostalCode:   addr.PostalCode,
		Country:      addr.Country,
		Latitude:     addr.Latitude,
		Longitude:    addr.Longitude,
	}); err != nil {
		return err
	}

	if err = tq.SaveCustomerAddress(ctx, sqlc.SaveCustomerAddressParams{
		CustomerID: database.ToPgTypeUUID(customerId),
		AddressID:  database.ToPgTypeUUID(addr.ID),
		IsDefault:  addr.Default,
	}); err != nil {
		return err
	}

	return transaction.Commit(ctx)
}

func (r repository) DeleteAddress(ctx context.Context, customerId uuid.UUID, addressId uuid.UUID) error {
    transaction, err := r.db.Begin(ctx)
    if err != nil {
        return xerrors.NewWithContext(ctx, err)
    }
    defer transaction.Rollback(ctx)
    
    tq := sqlc.New(transaction)

    _, err = tq.DeleteCustomerAddress(ctx, sqlc.DeleteCustomerAddressParams{
        CustomerID: database.ToPgTypeUUID(customerId),
        AddressID:  database.ToPgTypeUUID(addressId),
    })
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return xerrors.NewWithContext(ctx, fmt.Errorf("Address of ID %s does not exist", addressId.String()))
        }
        return err
    }

    _, err = tq.DeleteAddress(ctx, database.ToPgTypeUUID(addressId))
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return xerrors.NewWithContext(ctx, fmt.Errorf("Address of ID %s does not exist", addressId.String()))
        }
        return err
    }

    return transaction.Commit(ctx)
}