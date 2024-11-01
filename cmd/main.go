package main

import (
	"log/slog"
	"os"

	"github.com/oprimogus/cardapiogo/internal/api/router"
	"github.com/oprimogus/cardapiogo/internal/config"
	"github.com/oprimogus/cardapiogo/internal/database/postgres"
	"github.com/oprimogus/cardapiogo/internal/persistence"
	"github.com/oprimogus/cardapiogo/internal/services/adapter"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

//	@title			Cardapiogo API
//	@version		1.0
//	@description	Documentação da API de delivery Cardapiogo.
//	@contact.name	Gustavo Ferreira de Jesus
//	@contact.email	gustavo081900@gmail.com

//	@host		localhost:3000
//	@BasePath	/api
//	@accept		json
//	@produce	json

// @securityDefinitions.apikey BearerToken
// @in							header
// @name						Authorization
func main() {

	if config.GetInstance().Api.Environment == string(config.Production) {
		logger.InitLogger(os.Stdout, slog.LevelInfo)
	} else {
		logger.InitLogger(os.Stdout, slog.LevelDebug)
	}

	// database
	db := postgres.GetInstance()
	defer db.Close()

	// Service factory
	serviceFactory := adapter.NewServiceFactory()

	// Repository factory
	repositoryFactory := persistence.NewDataBaseRepositoryFactory(db, serviceFactory)

	// Web server
	router.Initialize(repositoryFactory, serviceFactory)
}
