package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	validatorutils "github.com/oprimogus/cardapiogo/internal/api/validator"
	"github.com/oprimogus/cardapiogo/internal/core/item"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	xerrors "github.com/oprimogus/cardapiogo/internal/errors"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
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
//	@Failure		400	{object}	xerrors.CustomError
//	@Failure		401	{object}	xerrors.CustomError
//	@Failure		403	{object}	xerrors.CustomError
//	@Failure		409	{object}	xerrors.CustomError
//	@Failure		500	{object}	xerrors.CustomError
//	@Failure		502	{object}	xerrors.CustomError
//	@Router			/v1/store/item [post]
func (c *ItemController) CreateItem(ctx *gin.Context) {
	transactionID := ctx.GetString(string(logger.TransactionIDKey))
	var params item.CreateItemInput
	err := ctx.BindJSON(&params)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	err = c.validator.Validate(transactionID, params)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
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
//	@Summary		Any user can view store item
//	@Description	Any user can view store item
//	@Tags			Item
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Item ID"
//	@Success		200	{object}	item.GetItemByIDOutput
//	@Failure		404	{object}	xerrors.CustomError
//	@Failure		500	{object}	xerrors.CustomError
//	@Failure		502	{object}	xerrors.CustomError
//	@Router			/v1/store/item/{id} [get]
func (c *ItemController) GetItemByID(ctx *gin.Context) {
	transactionID := ctx.GetString(string(logger.TransactionIDKey))
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

// GetItemByFilter godoc
//
//	@Summary		Any user can view store items.
//	@Description	Any user can view store items.
//	@Tags			Item
//	@Accept			json
//	@Produce		json
//	@Param			request	body		item.GetItemFilterInput	true	"Item Filter"
//	@Success		200	{object}	[]item.GetItemByIDOutput
//	@Failure		404	{object}	xerrors.CustomError
//	@Failure		500	{object}	xerrors.CustomError
//	@Failure		502	{object}	xerrors.CustomError
//	@Router			/v1/store/item [get]
func (c *ItemController) GetItemByFilter(ctx *gin.Context) {
	transactionID := ctx.GetString(string(logger.TransactionIDKey))
	var params item.GetItemFilterInput
	err := ctx.BindJSON(&params)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	err = c.validator.Validate(transactionID, params)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	itemsInstance, err := c.itemModule.GetByFilter.Execute(ctx, params)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	ctx.JSON(http.StatusOK, itemsInstance)
}
