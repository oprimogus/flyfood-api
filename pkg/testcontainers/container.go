package testcontainers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/testcontainers/testcontainers-go"
)

type Container struct {
    instance testcontainers.Container
    name     string
	Port     string
}

func (c *Container) Kill(ctx context.Context) {
	if err := c.instance.Terminate(ctx); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not stop %s: %s", c.name, err))
	}
}
