package user

import (
	"context"
)

type UseCaseFindByEmail struct {
	repository Repository
}

func NewUseCaseFindByEmail(repository Repository) UseCaseFindByEmail {
	return UseCaseFindByEmail{
		repository: repository,
	}
}

func (c UseCaseFindByEmail) Execute(ctx context.Context, email string) (User, error) {
	return c.repository.FindByEmail(ctx, email)
}
