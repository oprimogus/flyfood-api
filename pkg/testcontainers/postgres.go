package testcontainers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	*Container
}

type PostgresConfig struct {
	DatabaseName string
	Username     string
	Password     string
}

func MakePostgres(ctx context.Context, cfg PostgresConfig) (*PostgresContainer, error) {
	postgresContainer, err := postgres.Run(ctx,
		"postgis/postgis:16-3.4",
		postgres.WithDatabase(cfg.DatabaseName),
		postgres.WithUsername(cfg.Username),
		postgres.WithPassword(cfg.Password),
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

	return &PostgresContainer{
		Container: &Container{
			instance: postgresContainer,
			name:     "postgres",
			Port:     port,
		},
	}, nil
}
