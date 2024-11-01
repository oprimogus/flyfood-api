package user

import (
	"context"
	"fmt"

	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

type UseCaseUpdate struct {
	repository Repository
}

func NewUseCaseUpdate(repository Repository) UseCaseUpdate {
	return UseCaseUpdate{
		repository: repository,
	}
}

func (c UseCaseUpdate) Execute(ctx context.Context, input UpdateProfileParams) error {
	userID, ok := ctx.Value(string(logger.UserIDKey)).(string)
	if !ok {
		return fmt.Errorf("invalid userID: '%s'", userID)
	}
	user := input.ToEntity()
	user.ID = userID
	return c.repository.Update(ctx, user)
}
