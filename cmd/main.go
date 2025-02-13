package main

import (
	"context"
	"errors"
	"github.com/oprimogus/cardapiogo/internal/api"
	"github.com/oprimogus/cardapiogo/internal/config"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/database/persistence"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/database/postgres"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/adapter"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//	@title			FlyFood API
//	@version		1.0
//	@description	Documentação da API de delivery FlyFood.
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
	if err := run(); err != nil {
		log.Fatal("deu ruim")
	}
}

func run() (err error) {
	configInstance := config.GetInstance()

	if configInstance.Api.Environment == string(config.Production) {
		logger.InitLogger(os.Stdout, slog.LevelInfo)
	} else {
		logger.InitLogger(os.Stdout, slog.LevelDebug)
	}

	ctx := context.Background()

	// Init database connection
	db := postgres.GetInstance()
	defer db.Close()

	// Init Service factory
	serviceFactory := adapter.NewServiceFactory()

	// Init repositories
	repoFactory := persistence.NewRepositoryFactory(db)

	// Web server
	handler := api.InitRouter(db, repoFactory, serviceFactory)

	srv := &http.Server{
		Addr:    ":" + configInstance.Api.Port,
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		slog.Info("timeout of 5 seconds.")
	}

	slog.Info("Server exiting")
	err = srv.Shutdown(context.Background())

	return err
}
