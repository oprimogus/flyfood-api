package controller

import (
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	validatorutils "github.com/oprimogus/cardapiogo/internal/api/validator"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	xerrors "github.com/oprimogus/cardapiogo/internal/errors"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

type StoreController struct {
	validator                 *validatorutils.Validator
	useCaseCreate             store.UseCaseCreate
	useCaseUpdate             store.UseCaseUpdate
	useCaseGetByID            store.UseCaseGetByID
	useCaseGetByFilter        store.UseCaseGetByFilter
	useCaseAddBusinessHour    store.UseCaseAddBusinessHour
	useCaseDeleteBusinessHour store.UseCaseDeleteBusinessHour
	useCaseSetProfileImage    store.UseCaseSetProfileImage
	useCaseSetHeaderImage     store.UseCaseSetHeaderImage
}

func NewStoreController(validator *validatorutils.Validator, repository store.Repository) *StoreController {
	return &StoreController{
		validator:                 validator,
		useCaseCreate:             store.NewUseCaseCreate(repository),
		useCaseUpdate:             store.NewUseCaseUpdate(repository),
		useCaseGetByID:            store.NewUseCaseGetByID(repository),
		useCaseGetByFilter:        store.NewUseCaseGetByFilter(repository),
		useCaseAddBusinessHour:    store.NewUseCaseAddBusinessHour(repository),
		useCaseDeleteBusinessHour: store.NewUseCaseDeleteBusinessHour(repository),
		useCaseSetProfileImage:    store.NewUseCaseSetProfileImage(repository),
		useCaseSetHeaderImage:     store.NewUseCaseSetHeaderImage(repository),
	}
}

// GetStoreByID godoc
//
//	@Summary		Retrieve a specific store
//	@Description	Returns detailed information about a store by its ID.
//	@Description	This endpoint is publicly accessible and provides store details including business hours,
//	@Description	contact information, location, and current status.
//	@Tags			Store
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Store ID - Unique identifier of the store"
//	@Success		200	{object}	store.GetStoreByIdOutput	"Store details including name, description, business hours, and contact info"
//	@Failure		404	{object}	xerrors.CustomError	"Store not found"
//	@Failure		500	{object}	xerrors.CustomError	"Internal server error"
//	@Failure		502	{object}	xerrors.CustomError	"Bad gateway"
//	@Router			/v1/store/{id} [get]
func (c *StoreController) GetStoreByID(ctx *gin.Context) {
	traceID := ctx.GetString(string(logger.TraceIDKey))
	id := ctx.Param("id")
	storeInstance, err := c.useCaseGetByID.Execute(ctx, id)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	ctx.JSON(http.StatusOK, storeInstance)
}

// GetStoreByFilter godoc
//
//	@Summary		Search and filter stores
//	@Description	Returns a list of stores based on various filter criteria including location, rating, and type.
//	@Description	Supports geographic search with radius filtering when latitude and longitude are provided.
//	@Tags			Store
//	@Accept			json
//	@Produce		json
//	@Param			range		query	int		false	"Maximum distance range in kilometers for geographical search"
//	@Param			score		query	int		false	"Minimum rating score (1-5)"
//	@Param			name		query	string	false	"Store name for partial matching"
//	@Param			city		query	string	false	"Filter by city name"
//	@Param			latitude	query	string	false	"Latitude for geographical search (requires longitude)"
//	@Param			longitude	query	string	false	"Longitude for geographical search (requires latitude)"
//	@Param			type		query	string	false	"Store type category (e.g., restaurant, cafe, bar)"
//	@Success		200	{array}		store.GetStoreByIdOutput	"List of stores matching the filter criteria"
//	@Failure		400	{object}	xerrors.CustomError	"Invalid filter parameters"
//	@Failure		404	{object}	xerrors.CustomError	"No stores found"
//	@Failure		500	{object}	xerrors.CustomError	"Internal server error"
//	@Failure		502	{object}	xerrors.CustomError	"Bad gateway"
//	@Router			/v1/store [get]
func (c *StoreController) GetStoreByFilter(ctx *gin.Context) {
	traceID := ctx.GetString(string(logger.TraceIDKey))
	var params store.GetStoresFilterInput
	params.Name = ctx.Query("name")
	params.City = ctx.Query("city")
	params.Latitude = ctx.Query("latitude")
	params.Longitude = ctx.Query("longitude")
	params.Type = store.ShopType(ctx.Query("type"))

	queryRange := ctx.Query("range")
	if queryRange != "" {
		rangeValue, err := strconv.Atoi(queryRange)
		if err != nil {
			xerror := xerrors.HandleError(err, traceID)
			ctx.JSON(xerror.Status, xerror)
			return
		}
		params.Range = rangeValue
	}

	queryScore := ctx.Query("score")
	if queryScore != "" {
		scoreValue, err := strconv.Atoi(queryScore)
		if err != nil {
			xerror := xerrors.HandleError(err, traceID)
			ctx.JSON(xerror.Status, xerror)
			return
		}
		params.Score = scoreValue
	}

	err := c.validator.Validate(traceID, params)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	storeList, errListStores := c.useCaseGetByFilter.Execute(ctx, params)
	if errListStores != nil {
		xerror := xerrors.HandleError(errListStores, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	ctx.JSON(http.StatusOK, storeList)
}

// Create godoc
//
//	@Summary		Create a new store
//	@Description	Allows authenticated users to register a new store in the system.
//	@Description	The store owner must provide essential information including name, location, and contact details.
//	@Tags			Store
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			Params	body	store.CreateParams	true	"Store creation parameters including name, address, contact info, and type"
//	@Success		201	{object}	store.CreatedStore	"Returns the created store ID"
//	@Failure		400	{object}	xerrors.CustomError	"Invalid input parameters"
//	@Failure		401	{object}	xerrors.CustomError	"Unauthorized - Bearer token missing or invalid"
//	@Failure		403	{object}	xerrors.CustomError	"Forbidden - User lacks permission"
//	@Failure		409	{object}	xerrors.CustomError	"Conflict - Store already exists"
//	@Failure		500	{object}	xerrors.CustomError	"Internal server error"
//	@Failure		502	{object}	xerrors.CustomError	"Bad gateway"
//	@Router			/v1/store [post]
func (c *StoreController) Create(ctx *gin.Context) {
	traceID := ctx.GetString(string(logger.TraceIDKey))
	var params store.CreateParams
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
	storeID, err := c.useCaseCreate.Execute(ctx, params)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	response := store.CreatedStore{ID: storeID}
	ctx.JSON(http.StatusCreated, response)
}

// Update godoc
//
//	@Summary		Update store information
//	@Description	Allows store owners to modify their store's details and settings.
//	@Description	All provided fields will be updated while maintaining existing data for omitted fields.
//	@Tags			Store
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			Params	body	store.UpdateParams	true	"Updated store information including modifiable fields"
//	@Success		200	{object}	nil	"Store successfully updated"
//	@Failure		400	{object}	xerrors.CustomError	"Invalid input parameters"
//	@Failure		401	{object}	xerrors.CustomError	"Unauthorized - Bearer token missing or invalid"
//	@Failure		403	{object}	xerrors.CustomError	"Forbidden - User is not the store owner"
//	@Failure		404	{object}	xerrors.CustomError	"Store not found"
//	@Failure		409	{object}	xerrors.CustomError	"Conflict with existing data"
//	@Failure		500	{object}	xerrors.CustomError	"Internal server error"
//	@Failure		502	{object}	xerrors.CustomError	"Bad gateway"
//	@Router			/v1/store [put]
func (c *StoreController) Update(ctx *gin.Context) {
	traceID := ctx.GetString(string(logger.TraceIDKey))
	var params store.UpdateParams
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

// AddBusinessHours godoc
//
//	@Summary		Add store business hours
//	@Description	Allows store owners to set or update their store's operating hours.
//	@Description	Supports setting different hours for each day of the week and special holiday schedules.
//	@Tags			Store
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			Params	body	store.StoreBusinessHoursParams	true	"Business hours configuration including days, opening and closing times"
//	@Success		200	{object}	nil	"Business hours successfully updated"
//	@Failure		400	{object}	xerrors.CustomError	"Invalid time format or parameters"
//	@Failure		401	{object}	xerrors.CustomError	"Unauthorized - Bearer token missing or invalid"
//	@Failure		403	{object}	xerrors.CustomError	"Forbidden - User is not the store owner"
//	@Failure		409	{object}	xerrors.CustomError	"Conflict with existing business hours"
//	@Failure		500	{object}	xerrors.CustomError	"Internal server error"
//	@Failure		502	{object}	xerrors.CustomError	"Bad gateway"
//	@Router			/v1/store/business-hours [put]
func (c *StoreController) AddBusinessHours(ctx *gin.Context) {
	traceID := ctx.GetString(string(logger.TraceIDKey))
	var params store.StoreBusinessHoursParams
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
	err = c.useCaseAddBusinessHour.Execute(ctx, params)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	ctx.Status(http.StatusOK)
}

// DeleteBusinessHours godoc
//
//	@Summary		Remove store business hours
//	@Description	Allows store owners to delete specific business hours or operating schedules.
//	@Description	Can be used to remove regular hours or special holiday schedules.
//	@Tags			Store
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			Params	body	store.StoreBusinessHoursParams	true	"Business hours to be removed"
//	@Success		200	{object}	nil	"Business hours successfully deleted"
//	@Failure		400	{object}	xerrors.CustomError	"Invalid parameters"
//	@Failure		401	{object}	xerrors.CustomError	"Unauthorized - Bearer token missing or invalid"
//	@Failure		403	{object}	xerrors.CustomError	"Forbidden - User is not the store owner"
//	@Failure		404	{object}	xerrors.CustomError	"Business hours not found"
//	@Failure		500	{object}	xerrors.CustomError	"Internal server error"
//	@Failure		502	{object}	xerrors.CustomError	"Bad gateway"
//	@Router			/v1/store/business-hours [delete]
func (c *StoreController) DeleteBusinessHours(ctx *gin.Context) {
	traceID := ctx.GetString(string(logger.TraceIDKey))
	var params store.StoreBusinessHoursParams
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
	err = c.useCaseDeleteBusinessHour.Execute(ctx, params)
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	ctx.Status(http.StatusOK)
}

// SetProfileImage godoc
//
//	@Summary		Upload store profile image
//	@Description	Allows store owners to upload or update their store's profile image.
//	@Description	Supports JPEG and PNG formats. Image will be processed and optimized for display.
//	@Tags			Store
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerToken
//	@Param			id		path		string	true	"Store ID"
//	@Param			file	formData	file	true	"Profile image file (JPEG/PNG only, max size: 5MB)"
//	@Success		200	{object}	setFileOutput	"Returns URL of the uploaded image"
//	@Failure		400	{object}	xerrors.CustomError	"Invalid file format or size"
//	@Failure		401	{object}	xerrors.CustomError	"Unauthorized - Bearer token missing or invalid"
//	@Failure		403	{object}	xerrors.CustomError	"Forbidden - User is not the store owner"
//	@Failure		413	{object}	xerrors.CustomError	"File size exceeds limit"
//	@Failure		500	{object}	xerrors.CustomError	"Internal server error"
//	@Failure		502	{object}	xerrors.CustomError	"Bad gateway"
//	@Router			/v1/store/{id}/profile-image [post]
func (c *StoreController) SetProfileImage(ctx *gin.Context) {
	traceID := ctx.GetString(string(logger.TraceIDKey))
	storeID := ctx.Param("id")
	if storeID == "" {
		xerror := xerrors.New(traceID, http.StatusBadRequest, "Store ID is required")
		ctx.JSON(xerror.Status, xerror)
		return
	}

	image, err := ctx.FormFile("file")
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	openedFile, err := image.Open()
	if err != nil {
		message := fmt.Sprintf("file sent is invalid: %s", err)
		xerror := xerrors.BadRequest(traceID, message)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	defer func(openedFile multipart.File) {
		err := openedFile.Close()
		if err != nil {
			slog.ErrorContext(ctx, "fail on close file", "err", err.Error())
		}
	}(openedFile)

	buffer := make([]byte, 512)
	if _, err := openedFile.Read(buffer); err != nil {
		message := fmt.Sprintf("file sent is invalid: %s", err)
		xerror := xerrors.BadRequest(traceID, message)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	fileType := http.DetectContentType(buffer)

	if fileType != "image/jpeg" && fileType != "image/png" {
		xerror := xerrors.BadRequest(traceID, "Unsupported file type")
		ctx.JSON(xerror.Status, xerror)
		return
	}

	url, err := c.useCaseSetProfileImage.Execute(ctx, storeID, image)

	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	ctx.JSON(http.StatusOK, setFileOutput{URL: url})
}

// SetHeaderImage godoc
//
//	@Summary		Upload store header image
//	@Description	Allows store owners to upload or update their store's header/banner image.
//	@Description	Supports JPEG and PNG formats. Image will be processed and optimized for display.
//	@Tags			Store
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerToken
//	@Param			id		path		string	true	"Store ID"
//	@Param			file	formData	file	true	"Header image file (JPEG/PNG only, max size: 5MB)"
//	@Success		200	{object}	setFileOutput	"Returns URL of the uploaded image"
//	@Failure		400	{object}	xerrors.CustomError	"Invalid file format or size"
//	@Failure		401	{object}	xerrors.CustomError	"Unauthorized - Bearer token missing or invalid"
//	@Failure		403	{object}	xerrors.CustomError	"Forbidden - User is not the store owner"
//	@Failure		413	{object}	xerrors.CustomError	"File size exceeds limit"
//	@Failure		500	{object}	xerrors.CustomError	"Internal server error"
//	@Failure		502	{object}	xerrors.CustomError	"Bad gateway"
//	@Router			/v1/store/{id}/header-image [post]
func (c *StoreController) SetHeaderImage(ctx *gin.Context) {
	traceID := ctx.GetString(string(logger.TraceIDKey))
	storeID := ctx.Param("id")
	if storeID == "" {
		xerror := xerrors.New(traceID, http.StatusBadRequest, "Store ID is required")
		ctx.JSON(xerror.Status, xerror)
		return
	}

	image, err := ctx.FormFile("file")
	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	openedFile, err := image.Open()
	if err != nil {
		message := fmt.Sprintf("file sent is invalid: %s", err)
		xerror := xerrors.BadRequest(traceID, message)
		ctx.JSON(xerror.Status, xerror)
		return
	}
	defer func(openedFile multipart.File) {
		err := openedFile.Close()
		if err != nil {
			slog.ErrorContext(ctx, "fail on close file", "err", err.Error())
		}
	}(openedFile)

	buffer := make([]byte, 512)
	if _, err := openedFile.Read(buffer); err != nil {
		message := fmt.Sprintf("file sent is invalid: %s", err)
		xerror := xerrors.BadRequest(traceID, message)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	fileType := http.DetectContentType(buffer)

	if fileType != "image/jpeg" && fileType != "image/png" {
		xerror := xerrors.BadRequest(traceID, "Unsupported file type")
		ctx.JSON(xerror.Status, xerror)
		return
	}

	url, err := c.useCaseSetHeaderImage.Execute(ctx, storeID, image)

	if err != nil {
		xerror := xerrors.HandleError(err, traceID)
		ctx.JSON(xerror.Status, xerror)
		return
	}

	ctx.JSON(http.StatusOK, setFileOutput{URL: url})
}
