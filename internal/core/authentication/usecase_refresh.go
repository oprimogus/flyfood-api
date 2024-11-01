package authentication

import (
	"context"
)

type UseCaseRefresh struct {
	repository Repository
}

func NewUseCaseRefresh(repository Repository) UseCaseRefresh {
	return UseCaseRefresh{
		repository: repository,
	}
}

func (s UseCaseRefresh) Execute(ctx context.Context, refreshToken string) (JWT, error) {
	return s.repository.RefreshToken(ctx, refreshToken)
}
