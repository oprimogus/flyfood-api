package customer

import (
	"context"
	"github.com/oprimogus/flyfood-api/internal/core/address"
	"github.com/oprimogus/flyfood-api/internal/infrastructure/services/nominatim"
)

type Service struct {
	r Repository
}

func NewService(r Repository) Service {
	return Service{r}
}

func (s *Service) FindCustomer(ctx context.Context, id string) (*Customer, error) {
	return s.r.FindByID(ctx, id)
}

func (s *Service) CreateCustomer(ctx context.Context, dto CreateProfileDTO) (*Customer, error) {
	c, err := NewCustomer(dto.ID, dto.Name, dto.LastName, dto.CPF, dto.Email, dto.Phone)
	if err != nil {
		return nil, err
	}

	return c, s.r.Save(ctx, c)
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

	geoData, err := nominatim.Search(ctx, nominatim.Query{
		Street:     addr.AddressLine1,
		City:       addr.City,
		State:      addr.State,
		Country:    addr.Country,
		PostalCode: addr.PostalCode,
	})
	if err == nil {
		addr.Latitude = geoData[0].Latitude
		addr.Longitude = geoData[0].Longitude
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
