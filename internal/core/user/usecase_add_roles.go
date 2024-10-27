package user

import (
	"context"
	"fmt"

	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

type useCaseAddRoles struct {
	repository Repository
}

func newUseCaseAddRoles(repository Repository) useCaseAddRoles {
	return useCaseAddRoles{repository: repository}
}

func (a useCaseAddRoles) Execute(ctx context.Context, roles []Role) error {
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
