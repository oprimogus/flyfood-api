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
	validator              *validatorutils.Validator
	useCaseCreate          item.UseCaseCreate
	useCaseUpdate          item.UseCaseUpdate
	useCaseGetByID         item.UseCaseGetByID
	useCaseGetItemByFilter item.UseCaseGetByFilter
	useCaseDelete          item.UseCaseDelete
}

func NewItemController(validator *validatorutils.Validator, repository item.Repository, storeRepository store.Repository) *ItemController {
	return &ItemController{
		validator:              validator,
		useCaseCreate:          item.NewUseCaseCreate(repository, storeRepository),
		useCaseUpdate:          item.NewUseCaseUpdate(repository, storeRepository),
		useCaseGetByID:         item.NewUseCaseGetByID(repository),
		useCaseGetItemByFilter: item.NewUseCaseGetByFilter(repository),
		useCaseDelete:          item.NewUseCaseDelete(repository, storeRepository),
	}
}

// CreateItem godoc
//
//	@Summary		Create a new store item
//	@Description	Creates a new item in the store. Only authenticated store owners can perform this operation.
//	@Description	The item will be associated with the owner's store and available in the catalog.
//	@Tags			Item
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			Params	body	item.CreateItemInput	true	"Item creation parameters including name, description, price, and type"
//	@Success		201		{object}	item.CreatedItem	"Returns the created item ID"
//	@Failure		400		{object}	xerrors.CustomError	"Invalid input parameters"
//	@Failure		401		{object}	xerrors.CustomError	"Unauthorized - Bearer token missing or invalid"
//	@Failure		403		{object}	xerrors.CustomError	"Forbidden - User is not the store owner"
//	@Failure		409		{object}	xerrors.CustomError	"Conflict - Item already exists"
//	@Failure		500		{object}	xerrors.CustomError	"Internal server error"
//	@Failure		502		{object}	xerrors.CustomError	"Bad gateway"
//	@Router			/v1/store/item [post]
func (c *ItemController) CreateItem(ctx *gin.Context) {
	traceID := ctx.GetString(string(logger.TraceIDKey))
	var params item.CreateItemInput
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
	itemID, err := c.useCaseCreate.Execute(ctx, params)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	response := item.CreatedItem{ID: itemID}
	ctx.JSON(http.StatusCreated, response)
}

// GetItemByID godoc
//
//	@Summary		Retrieve a specific store item
//	@Description	Returns detailed information about a specific store item by its ID.
//	@Description	This endpoint is publicly accessible and does not require authentication.
//	@Tags			Item
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Item ID - Unique identifier of the item"
//	@Success		200	{object}	item.GetItemByIDOutput	"Item details including name, description, price, and availability"
//	@Failure		404	{object}	xerrors.CustomError	"Item not found"
//	@Failure		500	{object}	xerrors.CustomError	"Internal server error"
//	@Failure		502	{object}	xerrors.CustomError	"Bad gateway"
//	@Router			/v1/store/item/{id} [get]
func (c *ItemController) GetItemByID(ctx *gin.Context) {
	traceID := ctx.GetString(string(logger.TraceIDKey))
	id := ctx.Param("id")
	idConverted, err := strconv.Atoi(id)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	itemInstance, err := c.useCaseGetByID.Execute(ctx, idConverted)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	ctx.JSON(http.StatusOK, itemInstance)
}

// GetItemByFilter godoc
//
//	@Summary		Search and filter store items
//	@Description	Returns a list of store items based on the provided filter criteria.
//	@Description	Supports filtering by type, name, city, score, and maximum price.
//	@Tags			Item
//	@Accept			json
//	@Produce		json
//	@Param			type		query	string	false	"Filter by item type"
//	@Param			name		query	string	false	"Filter by item name (partial match)"
//	@Param			city		query	string	false	"Filter by store city"
//	@Param			score		query	integer	false	"Filter by minimum score"
//	@Param			maxPrice	query	integer	false	"Filter by maximum price"
//	@Success		200	{array}		item.GetItemByIDOutput	"List of items matching the filter criteria"
//	@Failure		400	{object}	xerrors.CustomError	"Invalid filter parameters"
//	@Failure		404	{object}	xerrors.CustomError	"No items found"
//	@Failure		500	{object}	xerrors.CustomError	"Internal server error"
//	@Failure		502	{object}	xerrors.CustomError	"Bad gateway"
//	@Router			/v1/store/item [get]
func (c *ItemController) GetItemByFilter(ctx *gin.Context) {
	traceID := ctx.GetString(string(logger.TraceIDKey))
	var params item.GetItemFilterInput

	itemType := ctx.Query("type")
	if itemType != "" {
		params.Type = item.Type(itemType)
	}

	name := ctx.Query("name")
	if name != "" {
		params.Name = name
	}

	city := ctx.Query("city")
	if city != "" {
		params.City = city
	}

	score := ctx.Query("score")
	if score != "" {
		scoreInt, err := strconv.Atoi(score)
		if err != nil {
			xerror := xerrors.HandleError(err, traceID)
			ctx.JSON(xerror.Status, xerror)
			return
		}
		params.Score = scoreInt
	}

	maxPrice := ctx.Query("maxPrice")
	if maxPrice != "" {
		maxPriceInt, err := strconv.Atoi(maxPrice)
		if err != nil {
			xerror := xerrors.HandleError(err, traceID)
			ctx.JSON(xerror.Status, xerror)
			return
		}
		params.MaxPrice = maxPriceInt
	}

	err := c.validator.Validate(traceID, params)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	itemsInstance, err := c.useCaseGetItemByFilter.Execute(ctx, params)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	ctx.JSON(http.StatusOK, itemsInstance)
}

// UpdateItem godoc
//
//	@Summary		Update an existing store item
//	@Description	Modifies the details of an existing store item. Only authenticated store owners can update their own items.
//	@Description	All fields in the update request will replace the existing values.
//	@Tags			Item
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			Params	body	item.UpdateItemInput	true	"Updated item details including name, description, price, and availability"
//	@Success		200		{object}	nil	"Item successfully updated"
//	@Failure		400		{object}	xerrors.CustomError	"Invalid input parameters"
//	@Failure		401		{object}	xerrors.CustomError	"Unauthorized - Bearer token missing or invalid"
//	@Failure		403		{object}	xerrors.CustomError	"Forbidden - User is not the store owner"
//	@Failure		404		{object}	xerrors.CustomError	"Item not found"
//	@Failure		409		{object}	xerrors.CustomError	"Conflict with existing data"
//	@Failure		500		{object}	xerrors.CustomError	"Internal server error"
//	@Failure		502		{object}	xerrors.CustomError	"Bad gateway"
//	@Router			/v1/store/item [put]
func (c *ItemController) UpdateItem(ctx *gin.Context) {
	traceID := ctx.GetString(string(logger.TraceIDKey))
	var params item.UpdateItemInput
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

// DeleteItem godoc
//
//	@Summary		Delete a store item
//	@Description	Removes an item from the store catalog. Only authenticated store owners can delete their own items.
//	@Description	This operation is irreversible. All associated data will be permanently removed.
//	@Tags			Item
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			id			path	string	true	"Item ID to be deleted"
//	@Success		200			{object}	nil	"Item successfully deleted"
//	@Failure		400			{object}	xerrors.CustomError	"Invalid item ID"
//	@Failure		401			{object}	xerrors.CustomError	"Unauthorized - Bearer token missing or invalid"
//	@Failure		403			{object}	xerrors.CustomError	"Forbidden - User is not the store owner"
//	@Failure		404			{object}	xerrors.CustomError	"Item not found"
//	@Failure		500			{object}	xerrors.CustomError	"Internal server error"
//	@Failure		502			{object}	xerrors.CustomError	"Bad gateway"
//	@Router			/v1/store/item/{id} [delete]
func (c *ItemController) DeleteItem(ctx *gin.Context) {
	traceID := ctx.GetString(string(logger.TraceIDKey))
	itemID := ctx.Param("id")
	itemIDConverted, err := strconv.Atoi(itemID)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	params := item.DeleteInput{
		ID: itemIDConverted,
	}

	err = c.validator.Validate(traceID, params)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	err = c.useCaseDelete.Execute(ctx, params.ID)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	ctx.Status(http.StatusOK)
}
