package identity

import (
	"context"

	"github.com/oprimogus/cardapiogo/internal/core/authentication"
	"github.com/oprimogus/cardapiogo/internal/core/user"
)

type Service interface {
	Create(ctx context.Context, user user.User) error
	FindByID(ctx context.Context, id string) (user.User, error)
	FindByEmail(ctx context.Context, email string) (user.User, error)
	Update(ctx context.Context, user user.User) error
	Delete(ctx context.Context, id string) error
	ResetPasswordByEmail(ctx context.Context, id string) error
	AddRoles(ctx context.Context, userID string, roles []user.Role) error
	GetUsersWithUniqueParams(ctx context.Context, params user.GetUsersWithUniqueParams) (*[]user.User, error)
	SignIn(ctx context.Context, email, password string) (authentication.JWT, error)
	RefreshToken(ctx context.Context, refreshToken string) (authentication.JWT, error)
	IsValidToken(ctx context.Context, token string) (bool, error)
	DecodeAccessToken(ctx context.Context, accessToken string) (map[string]interface{}, error)
}
