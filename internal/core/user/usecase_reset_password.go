package user

import (
	"context"
	"fmt"

	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

type UseCaseResetPassword struct {
	repository Repository
}

func newResetPassword(repository Repository) UseCaseResetPassword {
	return UseCaseResetPassword{
		repository: repository,
	}
}

func (r UseCaseResetPassword) Execute(ctx context.Context, id string) error {
	userID, ok := ctx.Value(string(logger.UserIDKey)).(string)
	if !ok {
		return fmt.Errorf("invalid userID: '%s'", userID)
	}
	return r.repository.ResetPasswordByEmail(ctx, id)
}
