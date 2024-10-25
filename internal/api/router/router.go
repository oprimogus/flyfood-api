package router

import (
	"log"

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

	logger := logger.NewLogger("Gin Router")
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.LoggerMiddleware(logger))

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

	logger.Infof("Docs available in http://localhost:%s/api/v1/reference/index.html", port)
	logger.Infof("Docs available in http://localhost:%s/api/v2/reference", port)
	log.Printf("Listening and serving in HTTP :%v\n", port)
	err = router.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
