package user

import (
	"context"
	"fmt"

	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

type useCaseResetPassword struct {
	repository Repository
}

func newResetPassword(repository Repository) useCaseResetPassword {
	return useCaseResetPassword{
		repository: repository,
	}
}

func (r useCaseResetPassword) Execute(ctx context.Context, id string) error {
	userID, ok := ctx.Value(string(logger.UserIDKey)).(string)
	if !ok {
		return fmt.Errorf("invalid userID: '%s'", userID)
	}
	return r.repository.ResetPasswordByEmail(ctx, id)
}
