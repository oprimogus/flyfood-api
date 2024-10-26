package router

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/oprimogus/cardapiogo/internal/api/middleware"
	"github.com/oprimogus/cardapiogo/internal/api/router/routes"
	validatorutils "github.com/oprimogus/cardapiogo/internal/api/validator"
	"github.com/oprimogus/cardapiogo/internal/config"
	"github.com/oprimogus/cardapiogo/internal/core"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

// Initialize API
func Initialize(factory core.RepositoryFactory) {
	validator, err := validatorutils.NewValidator("pt")
	if err != nil && validator == nil {
		panic(err)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logger.LoggerMiddleware())
	router.Use(middleware.CorsMiddleware())

	metrics := middleware.NewPrometheusMetrics()
	router.Use(middleware.PrometheusMiddleware(metrics))

	router.MaxMultipartMemory = 8 << 20

	routes.DefaultRoutes(router, metrics.Registry)
	routes.SwaggerRoutes(router)
	routes.AuthRoutes(router, validator, factory)
	routes.UserRoutes(router, validator, factory)
	routes.StoreRoutes(router, validator, factory)
	routes.ItemRoutes(router, validator, factory)
	config := config.GetInstance().Api
	port := config.Port()
	if port == "" {
		port = "3000"
	}

	ctx := context.Background()
	slog.InfoContext(ctx, fmt.Sprintf("Docs available in http://localhost:%s/api/v1/reference/index.html", port))
	slog.InfoContext(ctx, fmt.Sprintf("Docs available in http://localhost:%s/api/v2/reference", port))
	slog.InfoContext(ctx, fmt.Sprintf("Listening and serving in 0.0.0.0:%v", port))
	err = router.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
