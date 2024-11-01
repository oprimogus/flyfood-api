package routes

import (
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/oprimogus/cardapiogo/api"
	"github.com/oprimogus/cardapiogo/internal/config"
	xerrors "github.com/oprimogus/cardapiogo/internal/errors"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

func SwaggerRoutes(router *gin.Engine) {
	basePath := config.GetInstance().Api.BasePath()
	v1 := router.Group(basePath)

	v1.GET("/v1/reference/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1.GET("/v2/reference", func(c *gin.Context) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./api/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Cardapiogo",
			},
			DarkMode: true,
		})
		if err != nil {
			traceID := c.GetString(string(logger.TraceIDKey))
			c.JSON(500, xerrors.InternalServer(traceID, err.Error()))
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
	})
}
