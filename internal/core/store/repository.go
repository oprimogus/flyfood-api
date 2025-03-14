package store

import (
	"context"
	"github.com/oprimogus/flyfood-api/internal/core/store/product"
)

type Repository interface {
	FindStoreByID(ctx context.Context, id string) (Store, error)
	FindStoreProductByID(ctx context.Context, id string) (product.Product, error)
	SaveStore(ctx context.Context, s Store) error
	SaveProduct(ctx context.Context, p product.Product) error
}
