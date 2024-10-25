package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/oprimogus/cardapiogo/internal/api/controller"
	"github.com/oprimogus/cardapiogo/internal/api/middleware"
	validatorutils "github.com/oprimogus/cardapiogo/internal/api/validator"
	"github.com/oprimogus/cardapiogo/internal/config"
	"github.com/oprimogus/cardapiogo/internal/core"
	"github.com/oprimogus/cardapiogo/internal/core/user"
)

func ItemRoutes(router *gin.Engine,
	validator *validatorutils.Validator,
	factory core.RepositoryFactory) {
	itemRepository := factory.NewItemRepository()
	storeRepository := factory.NewStoreRepository()
	authenticationRepository := factory.NewAuthenticationRepository()
	itemController := controller.NewItemController(validator, itemRepository, storeRepository)

	basePath := config.GetInstance().Api.BasePath()

	v1 := router.Group(basePath + "/v1")
	{
		v1.POST("/store/item",
			middleware.AuthenticationMiddleware(authenticationRepository),
			middleware.AuthorizationMiddleware([]user.Role{user.RoleOwner}),
			itemController.CreateItem)
		v1.GET("/store/item/:id",
			middleware.AuthenticationMiddleware(authenticationRepository),
			itemController.GetItemByID)
	}
}
