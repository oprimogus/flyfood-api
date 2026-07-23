package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oprimogus/flyfood-api/internal/api/middleware"
	"github.com/oprimogus/flyfood-api/internal/config"
	"github.com/oprimogus/flyfood-api/internal/core/customer"
	"github.com/oprimogus/flyfood-api/internal/core/store"
	"github.com/oprimogus/flyfood-api/internal/infra/database"
	"github.com/oprimogus/flyfood-api/internal/infra/health"
	"github.com/oprimogus/flyfood-api/internal/infra/logger"
	"github.com/oprimogus/flyfood-api/internal/infra/swagger"
)

//	@title			FlyFood API
//	@version		1.0
//	@description	Documentação da API de delivery FlyFood.
//	@contact.name	Gustavo Ferreira de Jesus
//	@contact.email	gustavo081900@gmail.com

//	@servers.url			https://flyfood.com.br/flyfood-api
//	@servers.description	Production API

//	@servers.url			https://dev.flyfood.com.br/flyfood-api
//	@servers.description	Staging API

//	@servers.url			http://localhost:3000/flyfood-api
//	@servers.description	Dev API

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl								https://auth.flyfood.com.br/oauth/v2/token
// @authorizationurl						https://auth.flyfood.com.br/oauth/v2/authorize
// @in										header
// @scope.openid							OpenID Connect basic login
// @scope.email							Access to user's email
// @scope.profile							Access to user's profile

func main() {
	if err := run(); err != nil {
		slog.Error("fatal error", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.InitLogger(os.Stdout)
	slog.Info("starting FlyFood API")

	cfg := config.Get()

	db, err := database.GetPostgres(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := db.ClosePG(); cerr != nil {
			slog.Error("failed to close DB", "err", cerr)
		}
	}()

	r := chi.NewRouter()
	r.Use(logger.LoggingMiddleware)
	r.Use(middleware.Recovery)
	r.Use(middleware.Cors)
	r.Use(middleware.JSON)
	r.Use(middleware.Prometheus)

	// Resources / Modules
	health.SetupHealthRoutes(r, db)
	swagger.SetupSwaggerRoutes(r)
	
	customer.SetupRoutes(ctx, r, db)
	store.SetupRoutes(ctx, r, db)

	port := cfg.API.Port
	if port == "" {
		port = "3000"
	}

	slog.Info(fmt.Sprintf("Docs available in http://localhost:%s%s/docs", port, cfg.API.BasePath))
	slog.Info(fmt.Sprintf("Listening and serving in 0.0.0.0:%v", port))

	_ = chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		middleware.MakePath(route)
		return nil
	})

	server := &http.Server{
		Addr:    ":" + cfg.API.Port,
		Handler: r,
	}

	go func() {
		slog.Info("server listening", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			cancel()
		}
	}()

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	slog.Info("shutting down server")

	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()

	if err := server.Shutdown(ctxTimeout); err != nil {
		slog.Error("server shutdown failed", "err", err)
		return err
	}

	slog.Info("server exited cleanly")
	return nil
}
