package item

import (
	"context"
)

type useCaseGetByFilter struct {
	repository Repository
}

func newUseCaseGetByFilter(repository Repository) useCaseGetByFilter {
	return useCaseGetByFilter{repository: repository}
}

func (u useCaseGetByFilter) Execute(ctx context.Context, filter GetItemFilterInput) (*[]GetItemFilterOutput, error) {
	return u.repository.GetItemsByFilter(ctx, filter)
}
