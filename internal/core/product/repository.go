package product

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*Product, error)
	Save(ctx context.Context, product *Product) error
}
