package user

import (
	"context"
)

type UseCaseFindByID struct {
	repository Repository
}

func NewUseCaseFindByID(repository Repository) UseCaseFindByID {
	return UseCaseFindByID{
		repository: repository,
	}
}

func (c UseCaseFindByID) Execute(ctx context.Context, email string) (User, error) {
	return c.repository.FindByID(ctx, email)
}
