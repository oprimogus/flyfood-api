package main

import (
	"github.com/oprimogus/cardapiogo/internal/api/router"
	"github.com/oprimogus/cardapiogo/internal/database/postgres"
	"github.com/oprimogus/cardapiogo/internal/persistence"
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

	// database
	db := postgres.GetInstance()
	defer db.Close()

	// routes
	factoryRepository := persistence.NewDataBaseRepositoryFactory(db)
	router.Initialize(factoryRepository)
}
