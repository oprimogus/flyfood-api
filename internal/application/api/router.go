package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/oprimogus/cardapiogo/internal/application/api/controller"
	"github.com/oprimogus/cardapiogo/internal/application/api/middleware"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/database/persistence"
	postgresDB "github.com/oprimogus/cardapiogo/internal/infrastructure/database/postgres"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/adapter"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oprimogus/cardapiogo/internal/config"
)

func InitRouter(db *postgresDB.Database, repoFactory persistence.RepositoryFactory,
	serviceFactory adapter.Factory) {

	r := chi.NewRouter()
	r.Use(middleware.Recovery)
	r.Use(middleware.JSON)
	r.Use(middleware.Logging)
	r.Use(middleware.Cors)
	r.Use(middleware.Prometheus)

	controller.SetupHealthRoutes(r, db)
	controller.SetupSwaggerRoutes(r)
	controller.SetupCustomerRoutes(r, repoFactory, serviceFactory)
	controller.SetupStoreRoutes(r, repoFactory, serviceFactory)

	configInstance := config.GetInstance().Api
	port := configInstance.Port
	if port == "" {
		port = "3000"
	}

	ctx := context.Background()
	slog.InfoContext(ctx, fmt.Sprintf("Docs available in http://localhost:%s/api/docs", port))
	slog.InfoContext(ctx, fmt.Sprintf("Listening and serving in 0.0.0.0:%v", port))

	if configInstance.Environment == string(config.Production) {
		srv := &http.Server{
			Addr:    ":" + port,
			Handler: r,
		}

		go func() {
			// service connections
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
		// catching ctx.Done(). timeout of 5 seconds.
		select {
		case <-ctx.Done():
			slog.Info("timeout of 5 seconds.")
		}
		slog.Info("Server exiting")

	} else {
		//err := router.Run(":" + port)
		//if err != nil {
		//	panic(err)
		//}
		log.Fatal(http.ListenAndServe(":"+port, r))
	}
}
