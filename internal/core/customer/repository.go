package customer

import "context"

type Repository interface {
	FindByID(ctx context.Context, id int) (*Customer, error)
	Save(ctx context.Context, c *Customer) error
}
