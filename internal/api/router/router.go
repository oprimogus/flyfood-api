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
	"github.com/oprimogus/cardapiogo/internal/services/adapter"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

// Initialize API
func Initialize(repositoryFactory core.RepositoryFactory, serviceFactory adapter.Factory) {
	validator, err := validatorutils.NewValidator("pt")
	if err != nil && validator == nil {
		panic(err)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logger.GinMiddleware())
	router.Use(middleware.CorsMiddleware())

	metrics := middleware.NewPrometheusMetrics()
	router.Use(middleware.PrometheusMiddleware(metrics))

	router.MaxMultipartMemory = 8 << 20

	routes.DefaultRoutes(router, metrics.Registry)
	routes.SwaggerRoutes(router)
	routes.AuthRoutes(router, validator, serviceFactory)
	routes.UserRoutes(router, validator, serviceFactory)
	routes.StoreRoutes(router, validator, repositoryFactory)
	routes.ItemRoutes(router, validator, repositoryFactory)
	configInstance := config.GetInstance().Api
	port := configInstance.Port()
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
