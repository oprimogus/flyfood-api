package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/oprimogus/flyfood-api/internal/infra/database"
)

func main() {
	err := PopulateLocalDatabase()
	if err != nil {
		panic(err)
	}
}

func PopulateLocalDatabase() error {
    ctx := context.Background()
	db, err := database.GetPostgres(ctx)
	if err != nil {
		panic(err)
	}
	defer func(db *database.Postgres) {
		err := db.ClosePG()
		if err != nil {
			slog.Error("fail on close DB Connection")
		}
	}(db)

	err = db.Ping(ctx)
	if err != nil {
		panic("Fail connect in local database")
	}
	mocks := getMocks()
	for _, v := range mocks {
		err := executeSQLFile(ctx, db, v)
		if err != nil {
			return err
		}
	}
	return nil
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

func executeSQLFile(ctx context.Context, db *database.Postgres, mock string) error {
	query, err := os.ReadFile(mock)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("erro ao ler mock %v: %w", mock, err)
	}
	_, err = db.Exec(ctx, string(query))
	if err != nil {
		errJson, _ := json.Marshal(err)
		log.Println(string(errJson))
		return fmt.Errorf("erro ao executar o mock %v: %v", mock, err)
	}
	log.Printf("Mock %v adicionado com sucesso\n", mock)
	return nil
}
