package user

import (
	"context"
	"fmt"

	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

type useCaseDelete struct {
	repository Repository
}

func newUseCaseDelete(repository Repository) useCaseDelete {
	return useCaseDelete{repository: repository}
}

func (d useCaseDelete) Execute(ctx context.Context) error {
	userID, ok := ctx.Value(string(logger.UserIDKey)).(string)
	if !ok {
		return fmt.Errorf("invalid userID: '%s'", userID)
	}
	return d.repository.Delete(ctx, userID)
}
