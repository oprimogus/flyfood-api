package item

import (
	"context"
)

type useCaseGetByFilter struct {
	repository Repository
}

func NewUseCaseGetByFilter(repository Repository) *useCaseGetByFilter {
	return &useCaseGetByFilter{repository: repository}
}

func (u *useCaseGetByFilter) Execute(ctx context.Context, filter GetItemFilter) (*[]GetItemFilter, error) {
	return u.repository.GetItemsByFilter(ctx, filter)
}
