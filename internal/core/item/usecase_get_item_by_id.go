package item

import (
	"context"
)

type UseCaseGetByID struct {
	repository Repository
}

func NewUseCaseGetByID(repository Repository) UseCaseGetByID {
	return UseCaseGetByID{repository: repository}
}

func (u UseCaseGetByID) Execute(ctx context.Context, id int) (Item, error) {
	return u.repository.GetItemByID(ctx, id)
}
