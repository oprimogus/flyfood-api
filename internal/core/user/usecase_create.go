package user

import (
	"context"
	"fmt"
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
	users, err := c.repository.GetUsers(ctx, input)
	if err != nil {
		return fmt.Errorf("fail on find users with input parameters: %w", err)
	}
	for _, v := range *users {
		if v.Email == input.Email {
			return ErrExistUserWithEmail
		}
		if v.Profile.Phone == input.Profile.Phone {
			return ErrExistUserWithPhone
		}
	}
	return c.repository.Create(ctx, input.ToEntity())
}
