package customer

import (
	"context"

	"github.com/google/uuid"
	"github.com/oprimogus/flyfood-api/internal/core/address"
	"github.com/oprimogus/flyfood-api/internal/infra/database"
	"github.com/oprimogus/flyfood-api/pkg/xerrors"
)

type Service struct {
	r Repository
}

func NewService(db *database.Postgres) Service {
	return Service{NewRepository(db)}
}

func (s *Service) FindCustomerByID(ctx context.Context, id uuid.UUID) (*Customer, error) {
	return s.r.FindByID(ctx, id)
}

func (s *Service) FindCustomerByExternalID(ctx context.Context, id string) (*Customer, error) {
	return s.r.FindByExternalID(ctx, id)
}

func (s *Service) CreateCustomer(ctx context.Context, dto CreateCustomerDTO) (*Customer, error) {
	c, err := NewCustomer(dto.ExternalID, dto.Name, dto.LastName, dto.Email, dto.CPF, dto.Phone)
	if err != nil {
		return nil, err
	}

	return c, s.r.Save(ctx, c)
}

func (s *Service) UpdateCustomerProfile(ctx context.Context, customerExternalId string, dto UpdateCustomerDTO) error {
	c, err := s.r.FindByExternalID(ctx, customerExternalId)
	if err != nil {
		return err
	}

	err = c.UpdateProfile(dto.Name, dto.LastName, dto.Phone)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, c)
}

func (s *Service) AddAddress(ctx context.Context, customerExternalId string, input address.CreateAddressDTO) error {
	c, err := s.r.FindByExternalID(ctx, customerExternalId)
	if err != nil {
		return err
	}
	
	addrs, err := s.r.FindAddressesByExternalCustomerID(ctx, customerExternalId)
	if err != nil {
		return err
	}

	if len(addrs) >= 5 {
		return xerrors.NewWithContext(ctx, errMaxAddresses)
	}

	addrID, err := uuid.NewV7()
	if err != nil {
		return err
	}
	
	addr := address.Address{
		ID:           addrID,
		Name:         input.Name,
		AddressLine1: input.AddressLine1,
		AddressLine2: input.AddressLine2,
		Neighborhood: input.Neighborhood,
		City:         input.City,
		State:        input.State,
		PostalCode:   input.PostalCode,
		Country:      input.Country,
	}

	return s.r.SaveAddress(ctx, c.ID, addr)
}

func (s *Service) RemoveAddress(ctx context.Context, customerExternalId string, addressId uuid.UUID) error {
	c, err := s.r.FindByExternalID(ctx, customerExternalId)
	if err != nil {
		return err
	}

	err = s.r.DeleteAddress(ctx, c.ID, addressId)
	if err != nil {
		return err
	}
	return nil
}
