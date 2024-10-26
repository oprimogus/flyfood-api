package store

import (
	"context"
	"fmt"

	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

type useCaseCreate struct {
	repository Repository
}

func newUseCaseCreate(repository Repository) useCaseCreate {
	return useCaseCreate{
		repository: repository,
	}
}

func (c useCaseCreate) Execute(ctx context.Context, params CreateParams) (id string, err error) {
	userID, ok := ctx.Value(string(logger.UserIDKey)).(string)
	if !ok {
		userID = ""
	}
	if userID == "" {
		return "", fmt.Errorf("invalid userID: '%s'", userID)
	}
	id, err = c.repository.Create(ctx, params.Entity(userID))
	if err != nil {
		return "", fmt.Errorf("could not Create a store for this user: %w", err)
	}
	return id, nil
}
