package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	err := PopulateLocalDatabase()
	if err != nil {
		panic(err)
	}
}

func PopulateLocalDatabase() error {
	db, err := getSQLDBConnection(createStringConn())
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			slog.Error("fail on close DB Connection")
		}
	}(db)

	err = db.Ping()
	if err != nil {
		panic("Fail connect in local database")
	}
	mocks := getMocks()
	for _, v := range mocks {
		err := executeSQLFile(db, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func createStringConn() string {
	dbHost := "localhost"
	dbPort := "5435"
	dbUsername := "cardapiogo"
	dbPassword := "cardapiogo"
	dbName := "postgres"
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost,
		dbPort,
		dbUsername,
		dbPassword,
		dbName,
	)
}

func getSQLDBConnection(connStr string) (*sql.DB, error) {
	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("database: could not open sql connection: %w", err)
	}
	return sqlDB, nil
}

func getMocks() []string {
	files, err := os.ReadDir("internal/infrastructure/database/sql/mocks")
	if err != nil {
		panic(err)
	}
	filesPath := make([]string, len(files))
	for i, v := range files {
		filesPath[i] = strings.Replace(v.Name(), ".sql", "", -1)
	}
	log.Print(filesPath)
	return filesPath
}

func executeSQLFile(db *sql.DB, mock string) error {
	query, err := os.ReadFile(fmt.Sprintf("internal/infrastructure/database/sql/mocks/%v.sql", mock))
	if err != nil {
		log.Println(err)
		return fmt.Errorf("erro ao ler mock %v: %w", mock, err)
	}
	_, err = db.Exec(string(query))
	if err != nil {
		log.Println(err)
		return fmt.Errorf("erro ao executar o mock %v: %v", mock, err)
	}
	log.Printf("Mock %v adicionado com sucesso\n", mock)
	return nil
}
