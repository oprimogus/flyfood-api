package authentication

import (
	"context"
)

type UseCaseSignIn struct {
	repository Repository
}

func NewUseCaseSignIn(repository Repository) UseCaseSignIn {
	return UseCaseSignIn{
		repository: repository,
	}
}

func (s UseCaseSignIn) Execute(ctx context.Context, email, password string) (JWT, error) {
	return s.repository.SignIn(ctx, email, password)
}
