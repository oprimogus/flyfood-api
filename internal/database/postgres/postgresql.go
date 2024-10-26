package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/oprimogus/cardapiogo/internal/config"
)

var (
	instance *PostgresDatabase
)

type PostgresDatabase struct {
	pool  *pgxpool.Pool
	sqlDB *sql.DB
}

func GetInstance() *PostgresDatabase {
	if instance == nil {
		instance = createInstance()
	}
	return instance
}

func createInstance() *PostgresDatabase {
	database := &PostgresDatabase{}
	strConnection := database.createStringConn()

	var err error
	database.pool, err = database.getPgxConnection(strConnection)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	database.sqlDB, err = database.getSQLDBConnection(strConnection)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
	return database
}

func (d PostgresDatabase) createStringConn() string {
	config := config.GetInstance()
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=public",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
	)
}

func (d PostgresDatabase) getPgxConnection(connStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("database: could not open pgx connection: %w", err)
	}
	return pool, nil
}

func (d PostgresDatabase) getSQLDBConnection(connStr string) (*sql.DB, error) {
	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("database: could not open sql connection: %w", err)
	}
	return sqlDB, nil
}

func (d PostgresDatabase) GetDB() *pgxpool.Pool {
	return d.pool
}

func (d PostgresDatabase) Close() {
	d.pool.Close()
	if d.sqlDB != nil {
		d.sqlDB.Close()
	}
}

func (d PostgresDatabase) Migrate() error {
	sourceURL := "file://internal/database/migrations"
	dbName := os.Getenv("DB_NAME")
	slog.Info("starting migration execution")
	driver, err := postgres.WithInstance(d.sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("database: could not create migration driver: %w", err)
	}
	migrator, err := migrate.NewWithDatabaseInstance(sourceURL, dbName, driver)
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("database: Could not create migrator: %w", err)
	}
	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("database: Could not apply migrations: %w", err)
	}
	if err == migrate.ErrNoChange {
		slog.Info("No migrations to run.")
	} else {
		slog.Info("Migrations applied successfully.")
	}
	return nil
}
