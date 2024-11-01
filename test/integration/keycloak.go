package integration

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/oprimogus/cardapiogo/internal/config"
	"github.com/oprimogus/cardapiogo/internal/utils"
	keycloak "github.com/stillya/testcontainers-keycloak"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func MakeKeycloak(ctx context.Context) (*Container, error) {
	_ = utils.SetWorkingDirToProjectRoot()
	keycloakContainer, err := keycloak.Run(ctx, "quay.io/keycloak/keycloak:24.0.4",
		testcontainers.WithWaitStrategy(wait.ForListeningPort("8080/tcp")),
		testcontainers.WithWaitStrategy(wait.ForLog("Running the server in development mode.")),
		keycloak.WithContextPath("/"),
		keycloak.WithRealmImportFile("test/integration/testdata/keycloak.json"),
		keycloak.WithAdminUsername("admin"),
		keycloak.WithAdminPassword("admin"),
	)
	if err != nil {
		return nil, fmt.Errorf("could not start Keycloak: %w", err)
	}
	hostPort, err := keycloakContainer.MappedPort(ctx, "8080")
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed to get mapped port: %s", err))
		return nil, err
	}
	configInstance := config.GetInstance().Keycloak
	slog.Info("", "value", configInstance)
	portFormatted := strings.Replace(string(hostPort), "/tcp", "", -1)
	configInstance.BaseURL = fmt.Sprintf("http://localhost:%s", portFormatted)
	configInstance.Realm = "cardapiogo"
	configInstance.ClientID = "cardapiogo"
	configInstance.ClientSecret = "**********"

	return &Container{name: "Keycloak", instance: keycloakContainer}, nil
}
