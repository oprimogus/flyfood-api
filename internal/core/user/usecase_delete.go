package user

import (
	"context"
	"fmt"

	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

type UseCaseDelete struct {
	repository Repository
}

func NewUseCaseDelete(repository Repository) UseCaseDelete {
	return UseCaseDelete{repository: repository}
}

func (d UseCaseDelete) Execute(ctx context.Context) error {
	userID, ok := ctx.Value(string(logger.UserIDKey)).(string)
	if !ok {
		return fmt.Errorf("invalid userID: '%s'", userID)
	}
	return d.repository.Delete(ctx, userID)
}
