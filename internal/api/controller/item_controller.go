package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oprimogus/cardapiogo/internal/api/middleware"
	validatorutils "github.com/oprimogus/cardapiogo/internal/api/validator"
	"github.com/oprimogus/cardapiogo/internal/core/item"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	xerrors "github.com/oprimogus/cardapiogo/internal/errors"
)

type ItemController struct {
	validator  *validatorutils.Validator
	itemModule item.ItemModule
}

func NewItemController(validator *validatorutils.Validator, repository item.Repository, storeRepository store.Repository) *ItemController {
	return &ItemController{
		validator:  validator,
		itemModule: item.NewItemModule(repository, storeRepository),
	}
}

// CreateItem godoc
//
//	@Summary		Owner can create store items.
//	@Description	Owner can create store items.
//	@Tags			Item
//	@Accept			json
//	@Produce		json
//
// @Security BearerToken
//
//	@Param			Params	body	item.CreateItemInput	true	"Params to create a store item"
//	@Success		201 {object}    item.CreatedItem
//	@Failure		400	{object}	xerrors.ErrorResponse
//	@Failure		401	{object}	xerrors.ErrorResponse
//	@Failure		403	{object}	xerrors.ErrorResponse
//	@Failure		409	{object}	xerrors.ErrorResponse
//	@Failure		500	{object}	xerrors.ErrorResponse
//	@Failure		502	{object}	xerrors.ErrorResponse
//	@Router			/v1/store/item [post]
func (c *ItemController) CreateItem(ctx *gin.Context) {
	transactionID := ctx.GetString(middleware.TransactionIDLabel)
	var params item.CreateItemInput
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
	itemID, err := c.itemModule.Create.Execute(ctx, params)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	response := item.CreatedItem{ID: itemID}
	ctx.JSON(http.StatusCreated, response)
}

// GetItemByID godoc
//
//	@Summary		Any user can view store items.
//	@Description	Any user can view store items.
//	@Tags			Item
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Item ID"
//	@Success		200	{object}	item.GetItemByIDOutput
//	@Failure		404	{object}	xerrors.ErrorResponse
//	@Failure		500	{object}	xerrors.ErrorResponse
//	@Failure		502	{object}	xerrors.ErrorResponse
//	@Router			/v1/store/item/{id} [get]
func (c *ItemController) GetItemByID(ctx *gin.Context) {
	transactionID := ctx.GetString(middleware.TransactionIDLabel)
	id := ctx.Param("id")
	idConverted, err := strconv.Atoi(id)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	itemInstance, err := c.itemModule.GetByID.Execute(ctx, idConverted)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	ctx.JSON(http.StatusOK, itemInstance)
}
