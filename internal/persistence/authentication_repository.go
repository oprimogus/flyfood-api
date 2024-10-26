package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/oprimogus/cardapiogo/internal/core/authentication"
	"github.com/oprimogus/cardapiogo/internal/services/keycloak"
)

type AuthenticationRepository struct{}

func NewAuthenticationRepository() *AuthenticationRepository {
	return &AuthenticationRepository{}
}

func (a *AuthenticationRepository) SignIn(ctx context.Context, email, password string) (authentication.JWT, error) {
	k, err := keycloak.GetInstance()
	if err != nil {
		return authentication.JWT{}, err
	}
	jwtInstance, err := k.Client.Login(ctx, k.ClientID, k.ClientSecret, k.Realm, email, password)
	if err != nil {
		return authentication.JWT{}, err
	}
	return authentication.JWT{
		AccessToken:      jwtInstance.AccessToken,
		IDToken:          jwtInstance.IDToken,
		ExpiresIn:        jwtInstance.ExpiresIn,
		RefreshExpiresIn: jwtInstance.RefreshExpiresIn,
		RefreshToken:     jwtInstance.RefreshToken,
		TokenType:        jwtInstance.TokenType,
		NotBeforePolicy:  jwtInstance.NotBeforePolicy,
		SessionState:     jwtInstance.SessionState,
		Scope:            jwtInstance.Scope,
	}, nil
}

func (a *AuthenticationRepository) RefreshToken(ctx context.Context, refreshToken string) (authentication.JWT, error) {
	k, err := keycloak.GetInstance()
	if err != nil {
		return authentication.JWT{}, err
	}

	jwtInstance, err := k.Client.RefreshToken(ctx, refreshToken, k.ClientID, k.ClientSecret, k.Realm)
	if err != nil {
		return authentication.JWT{}, err
	}
	return authentication.JWT{
		AccessToken:      jwtInstance.AccessToken,
		IDToken:          jwtInstance.IDToken,
		ExpiresIn:        jwtInstance.ExpiresIn,
		RefreshExpiresIn: jwtInstance.RefreshExpiresIn,
		RefreshToken:     jwtInstance.RefreshToken,
		TokenType:        jwtInstance.TokenType,
		NotBeforePolicy:  jwtInstance.NotBeforePolicy,
		SessionState:     jwtInstance.SessionState,
		Scope:            jwtInstance.Scope,
	}, nil
}

func (a *AuthenticationRepository) IsValidToken(ctx context.Context, token string) (bool, error) {
	k, err := keycloak.GetInstance()
	if err != nil {
		return false, err
	}

	r, err := k.Client.RetrospectToken(ctx, token, k.ClientID, k.ClientSecret, k.Realm)
	if err != nil {
		return false, fmt.Errorf("ocurred an error while validate your access token: %w", err)
	}
	if a == nil {
		return false, errors.New("ocurred an error while validate your access token")
	}
	return *r.Active, nil
}

func (a *AuthenticationRepository) DecodeAccessToken(ctx context.Context, accessToken string) (map[string]interface{}, error) {
	k, err := keycloak.GetInstance()
	if err != nil {
		return map[string]interface{}{}, err
	}

	token, mapClaims, err := k.Client.DecodeAccessToken(ctx, accessToken, k.Realm)
	if err != nil {
		return nil, fmt.Errorf("unable to decode access token: %w", err)
	}
	if mapClaims == nil {
		return nil, fmt.Errorf("unable to decode access token and get claims")
	}
	if token == nil {
		return nil, fmt.Errorf("unable to decode access token and get metadata")
	}

	return *mapClaims, nil
}
