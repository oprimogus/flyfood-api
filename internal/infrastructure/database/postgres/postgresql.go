package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/oprimogus/flyfood-api/internal/config"
)

var (
	instance *Database
)

type Database struct {
	pool  *pgxpool.Pool
	sqlDB *sql.DB
}

func GetInstance() *Database {
	conf := config.GetInstance().Database
	if instance == nil {
		instance = createInstance(*conf)
	}
	if err := instance.GetDB().Ping(context.Background()); err != nil {
		panic(fmt.Sprintf("could not connect to database: %v", err))
	}
	return instance
}

func GetTestInstance(port string) *Database {
	conf := config.GetInstance().Database
	dbConf := config.DBConf{
		Host:     conf.Host,
		Port:     port,
		User:     conf.User,
		Password: conf.Password,
		Name:     conf.Name,
	}
	instance = createInstance(dbConf)
	if err := instance.GetDB().Ping(context.Background()); err != nil {
		panic(fmt.Sprintf("could not connect to database: %v", err))
	}
	return instance
}

func createInstance(conf config.DBConf) *Database {
	database := &Database{}
	strConnection := database.createStringConn(conf)

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

func (d Database) createStringConn(conf config.DBConf) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=public",
		conf.Host,
		conf.Port,
		conf.User,
		conf.Password,
		conf.Name,
	)
}

func (d Database) getPgxConnection(connStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("database: could not open pgx connection: %w", err)
	}
	return pool, nil
}

func (d Database) getSQLDBConnection(connStr string) (*sql.DB, error) {
	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("database: could not open sql connection: %w", err)
	}
	return sqlDB, nil
}

func (d Database) GetDB() *pgxpool.Pool {
	return d.pool
}

func (d Database) Close() {
	d.pool.Close()
}

func (d Database) Migrate() error {
	sourceURL := "file://internal/infrastructure/database/migrations"
	dbName := os.Getenv("DB_NAME")
	slog.Info("starting migration execution")
	driver, err := postgres.WithInstance(d.sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("database: could not create migration driver: %w", err)
	}
	migrator, err := migrate.NewWithDatabaseInstance(sourceURL, dbName, driver)
	if err != nil {
		return fmt.Errorf("database: Could not create migrator: %w", err)
	}
	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("database: Could not apply migrations: %w", err)
	}
	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info("No migrations to run.")
	} else {
		slog.Info("Migrations applied successfully.")
	}
	return nil
}
