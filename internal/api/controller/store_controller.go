package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/oprimogus/cardapiogo/internal/api/middleware"
	validatorutils "github.com/oprimogus/cardapiogo/internal/api/validator"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	xerrors "github.com/oprimogus/cardapiogo/internal/errors"
)

type StoreController struct {
	validator   *validatorutils.Validator
	storeModule store.StoreModule
}

func NewStoreController(validator *validatorutils.Validator, repository store.Repository) *StoreController {
	return &StoreController{
		validator:   validator,
		storeModule: store.NewStoreModule(repository),
	}
}

// GetStoreByID godoc
//
//	@Summary		Any user can view a store.
//	@Description	Fetches details of a store based on its ID.
//	@Tags			Store
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Store ID"
//	@Success		200	{object}	store.GetStoreByIdOutput
//	@Failure		404	{object}	xerrors.ErrorResponse
//	@Failure		500	{object}	xerrors.ErrorResponse
//	@Failure		502	{object}	xerrors.ErrorResponse
//	@Router			/v1/store/{id} [get]
func (c *StoreController) GetStoreByID(ctx *gin.Context) {
	transactionID := ctx.GetString(middleware.TransactionIDLabel)
	id := ctx.Param("id")
	storeInstance, err := c.storeModule.GetByID.Execute(ctx, id)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	ctx.JSON(http.StatusOK, storeInstance)
}

// GetStoreByFilter godoc
//
//	@Summary		Any user can view filtered stores.
//	@Description	Fetches stores based on a set of filters such as range, score, name, city, or type.
//	@Tags			Store
//	@Accept			json
//	@Produce		json
//	@Param			range		query		int		false	"Specify max range"
//	@Param			score		query		int		false	"Specify in score"
//	@Param			name		query		string	false	"Specify name like"
//	@Param			city		query		string	false	"Specify city"
//	@Param			latitude	query		string	false	"latitude of address selected"
//	@Param			longitude	query		string	false	"longitude of address selected"
//	@Param			type		query		string	false	"Specify store type"
//	@Success		200			{object}	[]store.GetStoreByIdOutput
//	@Failure		404			{object}	xerrors.ErrorResponse
//	@Failure		500			{object}	xerrors.ErrorResponse
//	@Failure		502			{object}	xerrors.ErrorResponse
//	@Router			/v1/store [get]
func (c *StoreController) GetStoreByFilter(ctx *gin.Context) {
	transactionID := ctx.GetString(middleware.TransactionIDLabel)
	var params store.StoreFilter
	params.Name = ctx.Query("name")
	params.City = ctx.Query("city")
	params.Latitude = ctx.Query("latitude")
	params.Longitude = ctx.Query("longitude")
	params.Type = store.ShopType(ctx.Query("type"))

	queryRange := ctx.Query("range")
	if queryRange != "" {
		rangeValue, err := strconv.Atoi(queryRange)
		if err != nil {
			xerror := xerrors.HandleError(err, transactionID)
			ctx.JSON(xerror.Status, xerror)
			return
		}
		params.Range = rangeValue
	}

	queryScore := ctx.Query("score")
	if queryScore != "" {
		scoreValue, err := strconv.Atoi(queryScore)
		if err != nil {
			xerror := xerrors.HandleError(err, transactionID)
			ctx.JSON(xerror.Status, xerror)
			return
		}
		params.Score = scoreValue
	}

	errValidate := c.validator.Validate(params, transactionID)
	if errValidate != nil {
		xerror := xerrors.HandleError(errValidate, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	storeList, err := c.storeModule.GetByFilter.Execute(ctx, params)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	ctx.JSON(http.StatusOK, storeList)
}

// Create godoc
//
//	@Summary		Owner can create stores.
//	@Description	Allows an owner to create a new store.
//	@Tags			Store
//	@Accept			json
//	@Produce		json
//	@Param			Params	body	store.CreateParams	true	"Parameters for creating a store"
//	@Success		201 {object}    store.CreatedStore
//	@Failure		400	{object}	xerrors.ErrorResponse
//	@Failure		401	{object}	xerrors.ErrorResponse
//	@Failure		403	{object}	xerrors.ErrorResponse
//	@Failure		409	{object}	xerrors.ErrorResponse
//	@Failure		500	{object}	xerrors.ErrorResponse
//	@Failure		502	{object}	xerrors.ErrorResponse
//	@Router			/v1/store [post]
func (c *StoreController) Create(ctx *gin.Context) {
	transactionID := ctx.GetString(middleware.TransactionIDLabel)
	var params store.CreateParams
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
	storeID, err := c.storeModule.Create.Execute(ctx, params)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	response := store.CreatedStore{ID: storeID}
	ctx.JSON(http.StatusCreated, response)
}

// Update godoc
//
//	@Summary		Owner can update stores.
//	@Description	Allows an owner to update the details of their store.
//	@Tags			Store
//	@Accept			json
//	@Produce		json
//	@Param			Params	body	store.UpdateParams	true	"Parameters for updating a store"
//	@Success		200
//	@Failure		400	{object}	xerrors.ErrorResponse
//	@Failure		401	{object}	xerrors.ErrorResponse
//	@Failure		403	{object}	xerrors.ErrorResponse
//	@Failure		409	{object}	xerrors.ErrorResponse
//	@Failure		500	{object}	xerrors.ErrorResponse
//	@Failure		502	{object}	xerrors.ErrorResponse
//	@Router			/v1/store [put]
func (c *StoreController) Update(ctx *gin.Context) {
	transactionID := ctx.GetString(middleware.TransactionIDLabel)
	var params store.UpdateParams
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
	err = c.storeModule.Update.Execute(ctx, params)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	ctx.Status(http.StatusOK)
}

// AddBusinessHours godoc
//
//	@Summary		Owner can update business hours.
//	@Description	Allows an owner to update the business hours of their store.
//	@Tags			Store
//	@Accept			json
//	@Produce		json
//	@Param			Params	body	store.StoreBusinessHoursParams	true	"Parameters for updating business hours"
//	@Success		200
//	@Failure		400	{object}	xerrors.ErrorResponse
//	@Failure		401	{object}	xerrors.ErrorResponse
//	@Failure		403	{object}	xerrors.ErrorResponse
//	@Failure		409	{object}	xerrors.ErrorResponse
//	@Failure		500	{object}	xerrors.ErrorResponse
//	@Failure		502	{object}	xerrors.ErrorResponse
//	@Router			/v1/store/business-hours [put]
func (c *StoreController) AddBusinessHours(ctx *gin.Context) {
	transactionID := ctx.GetString(middleware.TransactionIDLabel)
	var params store.StoreBusinessHoursParams
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
	err = c.storeModule.AddBusinessHour.Execute(ctx, params)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	ctx.Status(http.StatusOK)
}

// DeleteBusinessHours godoc
//
//	@Summary		Owner can delete business hours.
//	@Description	Allows an owner to delete business hours of their store.
//	@Tags			Store
//	@Accept			json
//	@Produce		json
//	@Param			Params	body	store.StoreBusinessHoursParams	true	"Parameters for deleting business hours"
//	@Success		200
//	@Failure		400	{object}	xerrors.ErrorResponse
//	@Failure		401	{object}	xerrors.ErrorResponse
//	@Failure		403	{object}	xerrors.ErrorResponse
//	@Failure		409	{object}	xerrors.ErrorResponse
//	@Failure		500	{object}	xerrors.ErrorResponse
//	@Failure		502	{object}	xerrors.ErrorResponse
//	@Router			/v1/store/business-hours [delete]
func (c *StoreController) DeleteBusinessHours(ctx *gin.Context) {
	transactionID := ctx.GetString(middleware.TransactionIDLabel)
	var params store.StoreBusinessHoursParams
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
	err = c.storeModule.DeleteBusinessHour.Execute(ctx, params)
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	ctx.Status(http.StatusOK)
}

// SetProfileImage godoc
//
//	@Summary		Owner can update store profile image.
//	@Description	Allows an owner to update the profile image of their store.
//	@Tags			Store
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id	path	string true	"Store ID"
//	@Param			file	formData	file true	"jpeg/png image"
//	@Success		200 {object} setFileOutput
//	@Failure		400	{object}	xerrors.ErrorResponse
//	@Failure		401	{object}	xerrors.ErrorResponse
//	@Failure		403	{object}	xerrors.ErrorResponse
//	@Failure		409	{object}	xerrors.ErrorResponse
//	@Failure		500	{object}	xerrors.ErrorResponse
//	@Failure		502	{object}	xerrors.ErrorResponse
//	@Router			/v1/store/{id}/profile-image [post]
func (c *StoreController) SetProfileImage(ctx *gin.Context) {
	transactionID := ctx.GetString(middleware.TransactionIDLabel)
	storeID := ctx.Param("id")
	if storeID == "" {
		xerror := xerrors.New(http.StatusBadRequest, "Store ID is required", transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	image, err := ctx.FormFile("file")
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	openedFile, err := image.Open()
	if err != nil {
		message := fmt.Sprintf("file sent is invalid: %s", err)
		xerr := xerrors.BadRequest(message, transactionID)
		ctx.JSON(xerr.Status, xerr)
		return
	}
	defer openedFile.Close()

	buffer := make([]byte, 512)
	if _, err := openedFile.Read(buffer); err != nil {
		message := fmt.Sprintf("file sent is invalid: %s", err)
		xerr := xerrors.BadRequest(message, transactionID)
		ctx.JSON(xerr.Status, xerr)
		return
	}

	fileType := http.DetectContentType(buffer)

	if fileType != "image/jpeg" && fileType != "image/png" {
		xerr := xerrors.BadRequest("Unsupported file type", transactionID)
		ctx.JSON(xerr.Status, xerr)
		return
	}

	url, err := c.storeModule.SetProfileImage.Execute(ctx, storeID, image)

	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	ctx.JSON(http.StatusOK, setFileOutput{URL: url})
}

// SetHeaderImage godoc
//
//	@Summary		Owner can update store header image.
//	@Description	Allows an owner to update the header image of their store.
//	@Tags			Store
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id	path	string true	"Store ID"
//	@Param			file	formData	file true	"jpeg/png image"
//	@Success		200 {object} setFileOutput
//	@Failure		400	{object}	xerrors.ErrorResponse
//	@Failure		401	{object}	xerrors.ErrorResponse
//	@Failure		403	{object}	xerrors.ErrorResponse
//	@Failure		409	{object}	xerrors.ErrorResponse
//	@Failure		500	{object}	xerrors.ErrorResponse
//	@Failure		502	{object}	xerrors.ErrorResponse
//	@Router			/v1/store/{id}/header-image [post]
func (c *StoreController) SetHeaderImage(ctx *gin.Context) {
	transactionID := ctx.GetString(middleware.TransactionIDLabel)
	storeID := ctx.Param("id")
	if storeID == "" {
		xerror := xerrors.New(http.StatusBadRequest, "Store ID is required", transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	image, err := ctx.FormFile("file")
	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	openedFile, err := image.Open()
	if err != nil {
		message := fmt.Sprintf("file sent is invalid: %s", err)
		xerr := xerrors.BadRequest(message, transactionID)
		ctx.JSON(xerr.Status, xerr)
		return
	}
	defer openedFile.Close()

	buffer := make([]byte, 512)
	if _, err := openedFile.Read(buffer); err != nil {
		message := fmt.Sprintf("file sent is invalid: %s", err)
		xerr := xerrors.BadRequest(message, transactionID)
		ctx.JSON(xerr.Status, xerr)
		return
	}

	fileType := http.DetectContentType(buffer)

	if fileType != "image/jpeg" && fileType != "image/png" {
		xerr := xerrors.BadRequest("Unsupported file type", transactionID)
		ctx.JSON(xerr.Status, xerr)
		return
	}

	url, err := c.storeModule.SetHeaderImage.Execute(ctx, storeID, image)

	if err != nil {
		xerror := xerrors.HandleError(err, transactionID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	ctx.JSON(http.StatusOK, setFileOutput{URL: url})
}
