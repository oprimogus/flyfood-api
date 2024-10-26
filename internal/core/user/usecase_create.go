package user

import (
	"context"
	// "log/slog"
)

type useCaseCreate struct {
	repository Repository
}

func newUseCaseCreate(repository Repository) useCaseCreate {
	return useCaseCreate{
		repository: repository,
	}
}

func (c useCaseCreate) Execute(ctx context.Context, input CreateParams) error {
	// existUser, err := c.repository.FindByEmail(ctx, input.Email)
	// if err != nil {
	// 	return err
	// }
	// if existUser.Email == input.Email {
	// 	slog.InfoContext(ctx, "fail on create user")
	// 	return ErrExistUserWithEmail
	// }
	return c.repository.Create(ctx, input.ToEntity())
}
