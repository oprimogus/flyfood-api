package customer

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*Customer, error)
	Save(ctx context.Context, c *Customer) error
}
