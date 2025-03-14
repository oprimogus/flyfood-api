package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/oprimogus/flyfood-api/internal/api/controller"
	"github.com/oprimogus/flyfood-api/internal/api/middleware"
	"github.com/oprimogus/flyfood-api/internal/infrastructure/database/persistence"
	postgresDB "github.com/oprimogus/flyfood-api/internal/infrastructure/database/postgres"
	"github.com/oprimogus/flyfood-api/internal/infrastructure/services/adapter"
	"log/slog"
	"net/http"

	"github.com/oprimogus/flyfood-api/internal/config"
)

func InitRouter(db *postgresDB.Database, repoFactory persistence.RepositoryFactory,
	serviceFactory adapter.Factory) http.Handler {

	r := chi.NewRouter()
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)
	r.Use(middleware.Cors)
	r.Use(middleware.JSON)
	r.Use(middleware.Prometheus)

	controller.SetupHealthRoutes(r, db)
	controller.SetupSwaggerRoutes(r)
	controller.SetupCustomerRoutes(r, repoFactory)
	controller.SetupOwnerRoutes(r, repoFactory, serviceFactory)
	controller.SetupStoreRoutes(r, repoFactory, serviceFactory)
	controller.SetupProductRoutes(r, repoFactory, serviceFactory)

	configInstance := config.GetInstance().Api
	port := configInstance.Port
	if port == "" {
		port = "3000"
	}

	slog.Info(fmt.Sprintf("Docs available in http://localhost:%s/api/docs", port))
	slog.Info(fmt.Sprintf("Listening and serving in 0.0.0.0:%v", port))

	_ = chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		middleware.MakePath(route)
		return nil
	})

	return r
}
