package customer

import (
	"context"
	"github.com/oprimogus/cardapiogo/internal/core/address"
	"github.com/oprimogus/cardapiogo/internal/core/customer"
)

type Service struct {
	r customer.Repository
}

func NewService(r customer.Repository) Service {
	return Service{r}
}

func (s *Service) FindCustomer(ctx context.Context, id string) (*customer.Customer, error) {
	return s.r.FindByID(ctx, id)
}

func (s *Service) CreateCustomer(ctx context.Context, dto CreateProfileDTO) error {
	c, err := customer.NewCustomer(dto.Name, dto.LastName, dto.CPF, dto.Email, dto.Phone)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, c)
}

func (s *Service) UpdateCustomerProfile(ctx context.Context, id string, dto UpdateProfileDTO) error {
	c, err := s.r.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = c.UpdateProfile(dto.Name, dto.LastName, dto.Phone)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, c)
}

func (s *Service) AddAddress(ctx context.Context, id string, addr address.Address) error {
	c, err := s.r.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = c.SaveNewAddress(addr)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, c)
}

func (s *Service) RemoveAddress(ctx context.Context, id string, addr address.Address) error {
	c, err := s.r.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = c.RemoveAddress(addr)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, c)
}
