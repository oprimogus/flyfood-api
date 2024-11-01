package user

import (
	"context"
	"fmt"

	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

type UseCaseAddRoles struct {
	repository Repository
}

func NewUseCaseAddRoles(repository Repository) UseCaseAddRoles {
	return UseCaseAddRoles{repository: repository}
}

func (a UseCaseAddRoles) Execute(ctx context.Context, roles []Role) error {
	userID, ok := ctx.Value(string(logger.UserIDKey)).(string)
	if !ok {
		return fmt.Errorf("invalid userID: '%s'", userID)
	}
	err := a.repository.AddRoles(ctx, userID, roles)
	if err != nil {
		return fmt.Errorf("fail on add roles to user: %w", err)
	}
	return nil
}
