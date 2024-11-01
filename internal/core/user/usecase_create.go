package user

import (
	"context"
	"fmt"
)

type UseCaseCreate struct {
	repository Repository
}

func NewUseCaseCreate(repository Repository) UseCaseCreate {
	return UseCaseCreate{
		repository: repository,
	}
}

func (c UseCaseCreate) Execute(ctx context.Context, input CreateParams) error {
	getUsersWithUniqueParams := GetUsersWithUniqueParams{
		Email: input.Email,
		Phone: input.Profile.Phone,
	}
	users, err := c.repository.GetUsersWithUniqueParams(ctx, getUsersWithUniqueParams)
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
