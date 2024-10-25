package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oprimogus/cardapiogo/internal/api/middleware"
	validatorutils "github.com/oprimogus/cardapiogo/internal/api/validator"
	"github.com/oprimogus/cardapiogo/internal/core/authentication"
	xerrors "github.com/oprimogus/cardapiogo/internal/errors"
)

// UserController struct
type AuthController struct {
	validator            *validatorutils.Validator
	authenticationModule authentication.AuthenticationModule
}

func NewAuthController(validator *validatorutils.Validator, authRepository authentication.Repository) *AuthController {
	return &AuthController{
		validator:            validator,
		authenticationModule: authentication.NewAuthenticationModule(authRepository),
	}
}

// SignIn godoc
//
//	@Summary		Sign-In with email and password
//	@Description	Authenticate a user using their email and password. If successful, a JWT is returned which can be used for subsequent authenticated requests.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		authentication.SignInParams	false	"Request body containing email and password."
//	@Success		200		{object}	authentication.JWT		"Returns the generated JWT upon successful authentication."
//	@Failure		400		{object}	xerrors.ErrorResponse	"Bad request due to validation errors or malformed input."
//	@Failure		500		{object}	xerrors.ErrorResponse	"Internal server error, something went wrong while processing the request."
//	@Failure		502		{object}	xerrors.ErrorResponse	"Bad gateway, likely due to an external service failure."
//	@Router			/v1/auth/sign-in [post]
func (c *AuthController) SignIn(ctx *gin.Context) {
	var params authentication.SignInParams
	transactionID := ctx.GetString(middleware.TransactionIDLabel)
	err := ctx.BindJSON(&params)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	errValidate := c.validator.Validate(params, transactionID)
	if errValidate != nil {
		xerror := xerrors.HandleError(errValidate, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	jwt, err := c.authenticationModule.SignIn.Execute(ctx, params.Email, params.Password)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	ctx.JSON(http.StatusOK, jwt)
}

// RefreshUserToken godoc
//
//	@Summary		Refresh JWT using a refresh token
//	@Description	Refresh an expired access token by providing a valid refresh token. Returns a new JWT if the refresh is successful.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		authentication.RefreshParams	true	"Request body containing the refresh token."
//	@Success		200		{object}	authentication.JWT		"Returns a new JWT upon successful token refresh."
//	@Failure		400		{object}	xerrors.ErrorResponse	"Bad request due to validation errors or malformed input."
//	@Failure		500		{object}	xerrors.ErrorResponse	"Internal server error, something went wrong while processing the request."
//	@Failure		502		{object}	xerrors.ErrorResponse	"Bad gateway, likely due to an external service failure."
//	@Router			/v1/auth/refresh [post]
func (c *AuthController) RefreshUserToken(ctx *gin.Context) {
	var params authentication.RefreshParams
	transactionID := ctx.GetString(middleware.TransactionIDLabel)
	err := ctx.BindJSON(&params)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	errValidate := c.validator.Validate(params, transactionID)
	if errValidate != nil {
		xerror := xerrors.HandleError(errValidate, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	jwt, err := c.authenticationModule.Refresh.Execute(ctx, params.RefreshToken)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	ctx.JSON(http.StatusOK, jwt)
}
