package store

import (
	"context"
	"fmt"

	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

type UseCaseCreate struct {
	repository Repository
}

func NewUseCaseCreate(repository Repository) UseCaseCreate {
	return UseCaseCreate{
		repository: repository,
	}
}

func (c UseCaseCreate) Execute(ctx context.Context, params CreateParams) (id string, err error) {
	userID, ok := ctx.Value(string(logger.UserIDKey)).(string)
	if !ok {
		return "", fmt.Errorf("invalid userID: '%s'", userID)
	}
	id, err = c.repository.Create(ctx, params.Entity(userID))
	if err != nil {
		return "", fmt.Errorf("could not Create a store for this user: %w", err)
	}
	return id, nil
}
