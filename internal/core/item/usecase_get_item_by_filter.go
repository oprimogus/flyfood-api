package item

import (
	"context"
)

type UseCaseGetByFilter struct {
	repository Repository
}

func NewUseCaseGetByFilter(repository Repository) UseCaseGetByFilter {
	return UseCaseGetByFilter{repository: repository}
}

func (u UseCaseGetByFilter) Execute(ctx context.Context, filter GetItemFilterInput) (*[]GetItemFilterOutput, error) {
	return u.repository.GetItemsByFilter(ctx, filter)
}
