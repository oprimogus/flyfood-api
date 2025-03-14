package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
	"os"
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
	dbPort := "5432"
	dbUsername := "flyfood-api"
	dbPassword := "flyfood-api"
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
	files, err := os.ReadDir("test/data")
	if err != nil {
		panic(err)
	}
	filesPath := make([]string, len(files))
	for i, v := range files {
		filesPath[i] = fmt.Sprintf("test/data/%s", v.Name())
	}
	return filesPath
}

func executeSQLFile(db *sql.DB, mock string) error {
	query, err := os.ReadFile(mock)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("erro ao ler mock %v: %w", mock, err)
	}
	_, err = db.Exec(string(query))
	if err != nil {
		errJson, _ := json.Marshal(err)
		log.Println(string(errJson))
		return fmt.Errorf("erro ao executar o mock %v: %v", mock, err)
	}
	log.Printf("Mock %v adicionado com sucesso\n", mock)
	return nil
}
