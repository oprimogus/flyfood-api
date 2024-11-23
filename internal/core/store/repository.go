package store

import "context"

type Repository interface {
	BecomeOwner(ctx context.Context, customerID string) error
	FindStoreByID(ctx context.Context, id string) (*Store, error)
	FindOwnerByID(ctx context.Context, id string) (*Owner, error)
	IsOwner(ctx context.Context, customerID string) (bool, error)
	Save(ctx context.Context, s *Store) error
}
