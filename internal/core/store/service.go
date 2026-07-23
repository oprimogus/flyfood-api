package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/oprimogus/flyfood-api/internal/core/ownership"
	"github.com/oprimogus/flyfood-api/internal/infra/database"
	"github.com/oprimogus/flyfood-api/pkg/xerrors"
)

type Service interface {
    CreateNewStore(ctx context.Context, customerExternalID string, input CreateStoreDTO) error
    UpdateStore(ctx context.Context, customerExternalID string, input UpdateStoreDTO) error
    FindStoreByID(ctx context.Context, id string) (*Store, error)
}

type service struct {
	r Repository
	ownerShipService ownership.Service
}

func NewService(db *database.Postgres) Service {
	return &service{
		r:                NewRepository(db),
		ownerShipService: ownership.NewService(db),
	}
}

func (s *service) FindStoreByID(ctx context.Context, id string) (*Store, error) {
	return s.r.FindByID(ctx, id)
}

func (s *service) CreateNewStore(ctx context.Context, customerExternalID string, input CreateStoreDTO) error {
    ow, err := s.ownerShipService.IsOwner(ctx, customerExternalID)
    if err != nil {
        return err
    }
    if ow.ID == uuid.Nil {
        return ownership.ErrNotAnOwner(ctx)
    }
    
	st, err := NewStore(
		ow.ID,
		input.CNPJ,
		input.Name,
		input.Description,
		input.Phone,
		input.Address,
		input.Type)
	if err != nil {
		return xerrors.NewWithContext(ctx, err)
	}

	err = s.r.Save(ctx, st)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateStore(ctx context.Context, customerExternalID string, input UpdateStoreDTO) error {
    isOwnerOf, err := s.ownerShipService.IsOwnerOf(ctx, customerExternalID, input.ID)
    if err != nil {
        return err
    }
    if !isOwnerOf {
        return ownership.ErrNotOwnerOfResource(ctx)
    }
    
	st, err := s.r.FindByID(ctx, input.ID)
	if err != nil {
		return err
	}

	err = st.UpdateStore(
		input.Name,
		input.Description,
		input.Phone,
		input.Address,
		input.Type,
	)

	err = s.r.Save(ctx, st)
	if err != nil {
		return err
	}
	return nil
}
