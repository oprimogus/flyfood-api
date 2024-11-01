package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	validatorutils "github.com/oprimogus/cardapiogo/internal/api/validator"
	"github.com/oprimogus/cardapiogo/internal/core/user"
	xerrors "github.com/oprimogus/cardapiogo/internal/errors"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

type UserController struct {
	validator       *validatorutils.Validator
	useCaseCreate   user.UseCaseCreate
	useCaseUpdate   user.UseCaseUpdate
	useCaseAddRoles user.UseCaseAddRoles
}

func NewUserController(validator *validatorutils.Validator, userRepository user.Repository) *UserController {
	return &UserController{
		validator:       validator,
		useCaseCreate:   user.NewUseCaseCreate(userRepository),
		useCaseUpdate:   user.NewUseCaseUpdate(userRepository),
		useCaseAddRoles: user.NewUseCaseAddRoles(userRepository),
	}
}

// CreateUser godoc
//
// @Summary     Create a new user account
// @Description Register a new user with provided credentials and profile information
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       request body     user.CreateParams true "User registration data"
// @Success     201    {object}  nil              "User successfully created"
// @Failure     400    {object}  xerrors.CustomError "Invalid request data"
// @Failure     409    {object}  xerrors.CustomError "User already exists"
// @Failure     422    {object}  xerrors.CustomError "Validation error"
// @Failure     500    {object}  xerrors.CustomError "Internal server error"
// @Failure     502    {object}  xerrors.CustomError "External service error"
// @Router      /v1/auth/sign-up [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	traceID := ctx.GetString(string(logger.TraceIDKey))
	var params user.CreateParams
	err := ctx.BindJSON(&params)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	err = c.validator.Validate(traceID, params)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	err = c.useCaseCreate.Execute(ctx.Request.Context(), params)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	ctx.Status(http.StatusCreated)
}

// UpdateUser godoc
//
// @Summary     Update user profile
// @Description Modify existing user profile information
// @Tags        User
// @Accept      json
// @Produce     json
// @Security    BearerToken
// @Param       request body     user.UpdateProfileParams true "Updated user profile data"
// @Success     200    {object}  nil                     "Profile successfully updated"
// @Failure     400    {object}  xerrors.CustomError     "Invalid request data"
// @Failure     401    {object}  xerrors.CustomError     "Unauthorized"
// @Failure     403    {object}  xerrors.CustomError     "Forbidden"
// @Failure     404    {object}  xerrors.CustomError     "User not found"
// @Failure     422    {object}  xerrors.CustomError     "Validation error"
// @Failure     500    {object}  xerrors.CustomError     "Internal server error"
// @Router      /v1/user [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	traceID := ctx.GetString(string(logger.TraceIDKey))
	var params user.UpdateProfileParams
	err := ctx.BindJSON(&params)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	err = c.validator.Validate(traceID, params)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	err = c.useCaseUpdate.Execute(ctx, params)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	ctx.Status(http.StatusOK)
}

// AddRolesToUser godoc
//
// @Summary     Assign new roles to user
// @Description Add one or more roles to an existing user's permissions
// @Tags        User
// @Accept      json
// @Produce     json
// @Security    BearerToken
// @Param       request body     user.AddRolesParams true "Roles to be added"
// @Success     200    {object}  nil                 "Roles successfully added"
// @Failure     400    {object}  xerrors.CustomError "Invalid request data"
// @Failure     401    {object}  xerrors.CustomError "Unauthorized"
// @Failure     403    {object}  xerrors.CustomError "Forbidden"
// @Failure     404    {object}  xerrors.CustomError "User not found"
// @Failure     422    {object}  xerrors.CustomError "Validation error"
// @Failure     409    {object}  xerrors.CustomError "Role already assigned"
// @Failure     500    {object}  xerrors.CustomError "Internal server error"
// @Router      /v1/user/roles [post]
func (c *UserController) AddRolesToUser(ctx *gin.Context) {
	traceID := ctx.GetString(string(logger.TraceIDKey))
	var params user.AddRolesParams
	err := ctx.BindJSON(&params)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	err = c.validator.Validate(traceID, params)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	err = c.useCaseAddRoles.Execute(ctx, params.Roles)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	ctx.Status(http.StatusOK)
}
