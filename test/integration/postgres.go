package integration

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/oprimogus/cardapiogo/internal/config"
	postgresDB "github.com/oprimogus/cardapiogo/internal/database/postgres"
	"github.com/oprimogus/cardapiogo/internal/utils"
)

func MakePostgres(ctx context.Context) (*Container, error) {
	_ = utils.SetWorkingDirToProjectRoot()
	configInstance := config.GetInstance().Database
	configInstance.Host = "localhost"
	configInstance.User = "cardapiogo"
	configInstance.Name = "postgres"
	configInstance.Password = "cardapiogo"
	postgresContainer, err := postgres.Run(ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join("test", "integration", "testdata", "postgres-init.sh")),
		postgres.WithDatabase(configInstance.Name),
		postgres.WithUsername(configInstance.User),
		postgres.WithPassword(configInstance.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed to start container: %s", err))
		return nil, err
	}

	hostPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed to get mapped port: %s", err))
		return nil, err
	}
	configInstance.Port = strings.Replace(string(hostPort), "/tcp", "", -1)

	errOnMigration := postgresDB.GetInstance().Migrate()
	if errOnMigration != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed on do migrations: %s", errOnMigration))
		return nil, errOnMigration
	}

	return &Container{name: "postgres", instance: postgresContainer}, nil
}
