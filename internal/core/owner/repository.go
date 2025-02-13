package owner

import "context"

type Repository interface {
	SaveOwner(ctx context.Context, owner Owner) error
	FindOwnerByID(ctx context.Context, id string) (*Owner, error)
	IsOwner(ctx context.Context, customerID string) (bool, error)
	IsOwnerOf(ctx context.Context, customerID, storeID string) (bool, error)
}
