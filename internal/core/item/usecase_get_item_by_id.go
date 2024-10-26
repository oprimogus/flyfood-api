package item

import (
	"context"
)

type useCaseGetByID struct {
	repository Repository
}

func newUseCaseGetByID(repository Repository) useCaseGetByID {
	return useCaseGetByID{repository: repository}
}

func (u useCaseGetByID) Execute(ctx context.Context, id int) (GetItemByIDOutput, error) {
	return u.repository.GetItemByID(ctx, id)
}
