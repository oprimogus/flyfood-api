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

func (s *Service) FindCustomer(ctx context.Context, id int) (*customer.Customer, error) {
	return s.r.FindByID(ctx, id)
}

func (s *Service) CreateCustomer(ctx context.Context, dto CreateProfileDTO) (*customer.Customer, error) {
	c, err := customer.NewCustomer(dto.ID, dto.Name, dto.LastName, dto.CPF, dto.Email, dto.Phone)
	if err != nil {
		return nil, err
	}

	return c, s.r.Save(ctx, c)
}

func (s *Service) UpdateCustomerProfile(ctx context.Context, id int, dto UpdateProfileDTO) error {
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

func (s *Service) AddAddress(ctx context.Context, id int, addr address.Address) error {
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

func (s *Service) RemoveAddress(ctx context.Context, id int, addr address.Address) error {
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
