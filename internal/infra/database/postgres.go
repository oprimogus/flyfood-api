package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/oprimogus/flyfood-api/internal/config"
	"github.com/oprimogus/flyfood-api/internal/infra/database/sqlc"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
)

var (
	pgdb  *Postgres
	pOnce sync.Once
)

type Postgres struct {
	*pgxpool.Pool
	sqlDb   *sql.DB
	Querier sqlc.Querier
}

func GetPostgres(ctx context.Context) (*Postgres, error) {
	var err error

	pOnce.Do(func() {
		cfg := config.Get().Postgres
		pgdb, err = newPostgres(ctx, cfg)
		if err != nil {
			return
		}
		if err = pgdb.Migrate(ctx); err != nil {
			err = fmt.Errorf("could not migrate database: %w", err)
		}
	})

	if err != nil {
		return nil, err
	}

	return pgdb, nil
}

func newPostgres(ctx context.Context, cfg *config.Postgres) (*Postgres, error) {
	connectionString := createStringConn(cfg)

	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database with pgxpool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("could not ping to database with pgxpool: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	querier := sqlc.New(pool)
	p := &Postgres{pool, db, querier}

	return p, nil
}

func createStringConn(conf *config.Postgres) string {
	return strings.Join([]string{
		"host=" + conf.Host,
		"port=" + conf.Port,
		"user=" + conf.UserName,
		"password=" + conf.Password,
		"dbname=" + conf.DatabaseName,
		"sslmode=disable",
		"search_path=public",
	}, " ")
}

//go:embed migrations/*.sql
var migrationsFS embed.FS

func (d *Postgres) Migrate(ctx context.Context) error {
    slog.Info("Migrating...")
	goose.SetBaseFS(migrationsFS)

	fsys, err := fs.Sub(migrationsFS, "migrations")

	postgresLock, err := lock.NewPostgresSessionLocker()
	if err != nil {
		return fmt.Errorf("could not create postgres session locker: %w", err)
	}

	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		d.sqlDb,
		fsys,
		goose.WithVerbose(false),
		goose.WithSessionLocker(postgresLock),
	)

	if err != nil {
		return fmt.Errorf("database: could not set goose provider: %w", err)
	}

	result, err := provider.Up(ctx)
	if err != nil {
		return err
	}
	for i := range result {
		slog.Info(fmt.Sprintf("database migration succeeded: %s", result[i]))
	}
	slog.Info("All migrations have been performed!")

	return nil
}

func (d *Postgres) ClosePG() error {
	if err := d.sqlDb.Close(); err != nil {
		return fmt.Errorf("database: could not close database: %w", err)
	}
	d.Close()
	return nil
}
