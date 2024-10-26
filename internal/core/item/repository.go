package item

import "context"

type Repository interface {
	CreateItem(ctx context.Context, params CreateItemInput, score int) (id int, err error)
	GetItemByID(ctx context.Context, id int) (GetItemByIDOutput, error)
	GetItemsByFilter(ctx context.Context, filter GetItemFilterInput) (*[]GetItemFilterOutput, error)
	UpdateItem(ctx context.Context, params UpdateItemInput) error
	DeleteItem(ctx context.Context, id int) error
}
