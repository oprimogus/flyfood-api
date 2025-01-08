package main

import (
	"github.com/oprimogus/cardapiogo/internal/api"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/database/persistence"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/database/postgres"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/adapter"
	"log/slog"
	"os"

	"github.com/oprimogus/cardapiogo/internal/config"
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

	// Init database connection
	db := postgres.GetInstance()
	defer db.Close()

	// Init Service factory
	serviceFactory := adapter.NewServiceFactory()

	// Init repositories
	repoFactory := persistence.NewRepositoryFactory(db)

	// Web server
	api.InitRouter(db, repoFactory, serviceFactory)
}
