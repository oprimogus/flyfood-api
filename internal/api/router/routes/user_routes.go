package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/oprimogus/cardapiogo/internal/api/controller"
	"github.com/oprimogus/cardapiogo/internal/api/middleware"
	validatorutils "github.com/oprimogus/cardapiogo/internal/api/validator"
	"github.com/oprimogus/cardapiogo/internal/config"
	"github.com/oprimogus/cardapiogo/internal/core/user"
	"github.com/oprimogus/cardapiogo/internal/services/adapter"
)

func UserRoutes(router *gin.Engine,
	validator *validatorutils.Validator,
	factory adapter.Factory) {
	repository := factory.NewIdentityService()
	userController := controller.NewUserController(validator, repository)

	basePath := config.GetInstance().Api.BasePath()

	v1 := router.Group(basePath + "/v1")
	{
		v1.POST("/auth/sign-up", userController.CreateUser)
		v1.PUT("/user",
			middleware.AuthenticationMiddleware(repository),
			middleware.AuthorizationMiddleware([]user.Role{}),
			userController.UpdateUser)
		v1.POST("/user/roles",
			middleware.AuthenticationMiddleware(repository),
			middleware.AuthorizationMiddleware([]user.Role{}),
			userController.AddRolesToUser)
	}
}
