package ownership

import (
	"context"

	"github.com/oprimogus/flyfood-api/internal/core/customer"
	"github.com/oprimogus/flyfood-api/internal/infra/database"
)

type Service interface {
    IsOwnerOf(ctx context.Context, customerExternalID, resourceID string) (bool, error)
    IsOwner(ctx context.Context, customerExternalID string) (Owner, error)
}

type service struct {
	r Repository
	customerSvc customer.Service
}

func NewService(db *database.Postgres) Service {
	return &service{
	    r: NewRepository(db),
		customerSvc: customer.NewService(db),
	}
}

func (s *service) IsOwnerOf(ctx context.Context, customerExternalID, resourceID string) (bool, error) {
	owner, err := s.r.FindByCustomerExternalID(ctx, customerExternalID)
	if err != nil {
		return false, err
	}

	return owner.ID.String() == resourceID, nil
}

func (s *service) IsOwner(ctx context.Context, customerExternalID string) (Owner, error) {
	owner, err := s.r.FindByCustomerExternalID(ctx, customerExternalID)
	if err != nil {
		return Owner{}, err
	}

	return *owner, nil
}
