package integration

import (
	"context"
	"fmt"
	postgresDB "github.com/oprimogus/flyfood-api/internal/infrastructure/database/postgres"
	"log/slog"
	"strings"

	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/oprimogus/flyfood-api/internal/config"
)

func MakePostgres(ctx context.Context) (*Container, error) {
	configInstance := config.GetInstance().Database
	configInstance.Host = "localhost"
	configInstance.User = "flyfood-api"
	configInstance.Name = "postgres"
	configInstance.Password = "flyfood-api"
	postgresContainer, err := postgres.Run(ctx,
		"docker.io/postgres:16-alpine",
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
	port := strings.ReplaceAll(string(hostPort), "/tcp", "")

	errOnMigration := postgresDB.GetTestInstance(port).Migrate()
	if errOnMigration != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed on do migrations: %s", errOnMigration))
		return nil, errOnMigration
	}

	return &Container{name: "postgres", instance: postgresContainer, Port: port}, nil
}
