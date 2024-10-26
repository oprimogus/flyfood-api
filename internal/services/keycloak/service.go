package keycloak

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Nerzal/gocloak/v13"

	"github.com/oprimogus/cardapiogo/internal/config"
)

var (
	keycloakInstance *KeycloakService
)

type KeycloakService struct {
	Client       *gocloak.GoCloak
	Token        *gocloak.JWT
	Realm        string
	ClientID     string
	ClientSecret string
	mu           sync.Mutex
}

func GetInstance() (k *KeycloakService, err error) {
	if keycloakInstance == nil {
		keycloakInstance, err = newKeycloakService(context.TODO())
		if err != nil {
			return keycloakInstance, err
		}
	}
	return keycloakInstance, nil
}

func newKeycloakService(ctx context.Context) (k *KeycloakService, err error) {

	config := config.GetInstance().Keycloak

	client := gocloak.NewClient(config.BaseURL)
	token, err := client.LoginClient(ctx, config.ClientID, config.ClientSecret, config.Realm)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("fail on get keycloak client: %s", err))
		return nil, fmt.Errorf("fail on get keycloak client: %w", err)
	}
	service := &KeycloakService{
		Client:       client,
		Token:        token,
		Realm:        config.Realm,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
	}

	go service.startTokenRenewer(ctx)

	return service, nil
}

func shouldRefreshToken(expirationTime time.Time) bool {
	return time.Now().After(expirationTime)
}

func (k *KeycloakService) startTokenRenewer(ctx context.Context) {
	slog.InfoContext(ctx, "Starting renew service...", "service", "keycloak")

	refreshBefore := time.Second * 60 // 1 min

	accessTokenTicker, accessTokenExpireIn := k.startTicker(ctx, "accessToken", refreshBefore, k.Token.ExpiresIn)
	defer accessTokenTicker.Stop()

	refreshTokenTicker, refreshTokenExpireIn := k.startTicker(ctx, "refreshToken", refreshBefore, k.Token.RefreshExpiresIn)
	defer refreshTokenTicker.Stop()

	for {
		select {
		case <-accessTokenTicker.C:
			k.handleTokenRenewal(ctx, accessTokenExpireIn, false)
		case <-refreshTokenTicker.C:
			k.handleTokenRenewal(ctx, refreshTokenExpireIn, true)
		case <-ctx.Done():
			slog.InfoContext(ctx, "AccessToken renewal stopped", "service", "keycloak")
			return
		}
	}
}

func (k *KeycloakService) startTicker(ctx context.Context, name string, renewBefore time.Duration, tokenExpiresIn int) (ticker *time.Ticker, willExpireIn time.Time) {
	actualTime := time.Now()
	tokenExpireIn := time.Second * time.Duration(tokenExpiresIn)
	willExpireIn = actualTime.Add(time.Duration(tokenExpireIn) - renewBefore)
	slog.InfoContext(ctx, fmt.Sprintf("token %s will expire in: %s", name, willExpireIn), "service", "keycloak")

	return time.NewTicker(time.Duration(tokenExpireIn) - renewBefore), willExpireIn
}

func (k *KeycloakService) handleTokenRenewal(ctx context.Context, tokenExpireIn time.Time, isRefreshToken bool) {
	k.mu.Lock()
	defer k.mu.Unlock()

	if shouldRefreshToken(tokenExpireIn) {
		var token *gocloak.JWT
		var err error

		if isRefreshToken {
			token, err = k.Client.LoginClient(ctx, k.ClientID, k.ClientSecret, k.Realm)
		} else {
			token, err = k.Client.RefreshToken(ctx, k.Token.RefreshToken, k.ClientID, k.ClientSecret, k.Realm)
		}

		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("Failed to renew token: %v", err), "service", "keycloak")
			return
		}

		k.Token = token
		slog.InfoContext(ctx, fmt.Sprintf("Token refreshed successfully (isRefreshToken: %v)", isRefreshToken), "service", "keycloak")
	}
}
