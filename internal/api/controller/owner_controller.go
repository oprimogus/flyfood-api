package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oprimogus/flyfood-api/internal/api/middleware"
	"github.com/oprimogus/flyfood-api/internal/config"
	"github.com/oprimogus/flyfood-api/internal/core/owner"
	"github.com/oprimogus/flyfood-api/internal/core/store"
	"github.com/oprimogus/flyfood-api/internal/infrastructure/database/persistence"
	"github.com/oprimogus/flyfood-api/internal/infrastructure/services/adapter"
	"github.com/oprimogus/flyfood-api/internal/infrastructure/services/zitadel"
	_ "github.com/oprimogus/flyfood-api/internal/xerrors"
	"github.com/oprimogus/flyfood-api/internal/xvalidator"
	"github.com/oprimogus/flyfood-api/pkg/converters"
)

type ownerController struct {
	command      owner.Command
	storeCommand store.Command
	storeQuery   store.Query
}

func newOwnerController(command owner.Command, storeCommand store.Command, storeQuery store.Query) ownerController {
	return ownerController{
		command:      command,
		storeCommand: storeCommand,
		storeQuery:   storeQuery,
	}
}

// createOwner godoc
//
//	@Summary		Owner: Create an owner
//	@Description	Register a customer as owner
//	@Tags			Owner V1
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string				true	"Bearer authentication token"
//	@Success		201				{object}	nil					"Store successfully created"
//	@Failure		400				{object}	xerrors.CustomError	"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError	"Unauthorized - authentication required"
//	@Failure		409				{object}	xerrors.CustomError	"Conflict - store may already exist"
//	@Failure		422				{object}	xerrors.CustomError	"Validation error - invalid store details"
//	@Failure		500				{object}	xerrors.CustomError	"Internal server error"
//	@Failure		502				{object}	xerrors.CustomError	"External service communication error"
//	@Router			/v1/owner [post]
func (c ownerController) createOwner(w http.ResponseWriter, r *http.Request) {
	zt := zitadel.GetInstance()
	authCtx := zt.GetContext(r.Context())

	err := c.command.NewOwner(r.Context(), authCtx.UserID())
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = zt.SetRole(r.Context(), authCtx.UserID(), zitadel.Owner)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// createStore godoc
//
//	@Summary		Store: Create a new store for an owner
//	@Description	Register a comprehensive store profile with all necessary details
//	@Tags			Store Management V1
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string					true	"Bearer authentication token"
//	@Param			request			body		store.CreateNewStoreDTO	true	"Detailed store creation information"
//	@Success		201				{object}	nil						"Store successfully created"
//	@Failure		400				{object}	xerrors.CustomError		"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError		"Unauthorized - authentication required"
//	@Failure		409				{object}	xerrors.CustomError		"Conflict - store may already exist"
//	@Failure		422				{object}	xerrors.CustomError		"Validation error - invalid store details"
//	@Failure		500				{object}	xerrors.CustomError		"Internal server error"
//	@Failure		502				{object}	xerrors.CustomError		"External service communication error"
//	@Router			/v1/owner/store [post]
func (c ownerController) createStore(w http.ResponseWriter, r *http.Request) {
	authCtx := zitadel.GetInstance().GetContext(r.Context())
	var params store.CreateNewStoreDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = xvalidator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeCommand.CreateNewStore(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// updateStore godoc
//
//	@Summary		Store: Update existing store profile
//	@Description	Modify comprehensive store details including contact, operational, and profile information
//	@Tags			Store Management V1
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string					true	"Bearer authentication token"
//	@Param			request			body		store.UpdateStoreDTO	true	"Updated store information"
//	@Success		200				{object}	nil						"Store successfully updated"
//	@Failure		400				{object}	xerrors.CustomError		"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError		"Unauthorized - authentication required"
//	@Failure		404				{object}	xerrors.CustomError		"Store not found"
//	@Failure		422				{object}	xerrors.CustomError		"Validation error - invalid store details"
//	@Failure		500				{object}	xerrors.CustomError		"Internal server error"
//	@Router			/v1/owner/store [put]
func (c ownerController) updateStore(w http.ResponseWriter, r *http.Request) {
	authCtx := zitadel.GetInstance().GetContext(r.Context())
	var params store.UpdateStoreDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = xvalidator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeCommand.Update(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// addBusinessHour godoc
//
//	@Summary		Store: Add business operating hours for store
//	@Description	Register specific business hours for a store's operational schedule
//	@Tags			Store Management V1
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string								true	"Bearer authentication token"
//	@Param			request			body		store.AddOrDeleteBusinessHourDTO	true	"Business hours details"
//	@Success		200				{object}	nil									"Business hours successfully added"
//	@Failure		400				{object}	xerrors.CustomError					"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError					"Unauthorized - authentication required"
//	@Failure		422				{object}	xerrors.CustomError					"Validation error - invalid business hours"
//	@Failure		500				{object}	xerrors.CustomError					"Internal server error"
//	@Router			/v1/owner/store/business-hour [post]
func (c ownerController) addBusinessHour(w http.ResponseWriter, r *http.Request) {
	authCtx := zitadel.GetInstance().GetContext(r.Context())
	var params store.AddOrDeleteBusinessHourDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = xvalidator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeCommand.AddStoreBusinessHour(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// removeBusinessHour godoc
//
//	@Summary		Store: Remove business operating hours for store
//	@Description	Delete specific business hours from store's operational schedule
//	@Tags			Store Management V1
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string								true	"Bearer authentication token"
//	@Param			request			body		store.AddOrDeleteBusinessHourDTO	true	"Business hours to remove"
//	@Success		200				{object}	nil									"Business hours successfully removed"
//	@Failure		400				{object}	xerrors.CustomError					"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError					"Unauthorized - authentication required"
//	@Failure		422				{object}	xerrors.CustomError					"Validation error - invalid business hours"
//	@Failure		500				{object}	xerrors.CustomError					"Internal server error"
//	@Router			/v1/owner/store/business-hour [delete]
func (c ownerController) removeBusinessHour(w http.ResponseWriter, r *http.Request) {
	authCtx := zitadel.GetInstance().GetContext(r.Context())
	var params store.AddOrDeleteBusinessHourDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = xvalidator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeCommand.RemoveStoreBusinessHour(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// addPaymentMethod godoc
//
//	@Summary		Store: Add payment method for store
//	@Description	Register a new payment method accepted by the store
//	@Tags			Store Management V1
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string								true	"Bearer authentication token"
//	@Param			request			body		store.AddOrDeletePaymentMethodDTO	true	"Payment method details"
//	@Success		200				{object}	nil									"Payment method successfully added"
//	@Failure		400				{object}	xerrors.CustomError					"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError					"Unauthorized - authentication required"
//	@Failure		422				{object}	xerrors.CustomError					"Validation error - invalid payment method"
//	@Failure		500				{object}	xerrors.CustomError					"Internal server error"
//	@Router			/v1/owner/store/payment-method [post]
func (c ownerController) addPaymentMethod(w http.ResponseWriter, r *http.Request) {
	authCtx := zitadel.GetInstance().GetContext(r.Context())
	var params store.AddOrDeletePaymentMethodDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = xvalidator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeCommand.AddStorePaymentMethod(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// removePaymentMethod godoc
//
//	@Summary		Store: Remove payment method for store
//	@Description	Delete a previously added payment method from store's accepted methods
//	@Tags			Store Management V1
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string								true	"Bearer authentication token"
//	@Param			request			body		store.AddOrDeletePaymentMethodDTO	true	"Payment method to remove"
//	@Success		200				{object}	nil									"Payment method successfully removed"
//	@Failure		400				{object}	xerrors.CustomError					"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError					"Unauthorized - authentication required"
//	@Failure		422				{object}	xerrors.CustomError					"Validation error - invalid payment method"
//	@Failure		500				{object}	xerrors.CustomError					"Internal server error"
//	@Router			/v1/owner/store/payment-method [delete]
func (c ownerController) removePaymentMethod(w http.ResponseWriter, r *http.Request) {
	authCtx := zitadel.GetInstance().GetContext(r.Context())
	var params store.AddOrDeletePaymentMethodDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = xvalidator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeCommand.RemoveStorePaymentMethod(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// openOrClose godoc
//
//	@Summary		Store: Change store's operational status for orders
//	@Description	Toggle store's availability to accept or stop accepting new orders
//	@Tags			Store Management V1
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string					true	"Bearer authentication token"
//	@Param			request			body		store.SetOpenStateDTO	true	"Store open/close status"
//	@Success		200				{object}	nil						"Store operational status updated successfully"
//	@Failure		400				{object}	xerrors.CustomError		"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError		"Unauthorized - authentication required"
//	@Failure		422				{object}	xerrors.CustomError		"Validation error - invalid status"
//	@Failure		500				{object}	xerrors.CustomError		"Internal server error"
//	@Router			/v1/owner/store/open [post]
func (c ownerController) openOrClose(w http.ResponseWriter, r *http.Request) {
	authCtx := zitadel.GetInstance().GetContext(r.Context())
	var params store.SetOpenStateDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = xvalidator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeCommand.OpenOrCloseStore(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// activatedOrDeactivate godoc
//
//	@Summary		Store: Change store's visibility and order acceptance status
//	@Description	Activate or deactivate store to control product visibility and order processing
//	@Tags			Store Management V1
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string				true	"Bearer authentication token"
//	@Param			request			body		store.SetActiveDTO	true	"Store activation status"
//	@Success		200				{object}	nil					"Store activation status updated successfully"
//	@Failure		400				{object}	xerrors.CustomError	"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError	"Unauthorized - authentication required"
//	@Failure		422				{object}	xerrors.CustomError	"Validation error - invalid status"
//	@Failure		500				{object}	xerrors.CustomError	"Internal server error"
//	@Router			/v1/owner/store/active [post]
func (c ownerController) activatedOrDeactivate(w http.ResponseWriter, r *http.Request) {
	authCtx := zitadel.GetInstance().GetContext(r.Context())
	var params store.SetActiveDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = xvalidator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeCommand.ActiveOrDeactivateStore(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// changeProfileImage godoc
//
//	@Summary		Store: Change store profile image
//	@Description	Change store profile image
//	@Tags			Store Management V1
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string				true	"Bearer authentication token"
//	@Param			id				path		string				true	"Store ID"
//	@Param			file			formData	file				true	"Store profile image"
//	@Success		200				{object}	nil					"Store activation status updated successfully"
//	@Failure		400				{object}	xerrors.CustomError	"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError	"Unauthorized - authentication required"
//	@Failure		422				{object}	xerrors.CustomError	"Validation error - invalid status"
//	@Failure		500				{object}	xerrors.CustomError	"Internal server error"
//	@Router			/v1/owner/store/{id}/profile-image [post]
func (c ownerController) changeProfileImage(w http.ResponseWriter, r *http.Request) {
	file, ext, err := GetFileFormData(w, r, int64(10), "image", []string{"image/jpeg", "image/png", "image/jpg"})
	if err != nil {
		HandleError(w, r, err)
		return
	}

	storeID := chi.URLParam(r, "id")

	zt := zitadel.GetInstance()
	authCtx := zt.Middleware.Context(r.Context())

	fileBytes, err := converters.FileToBytes(file)
	if err != nil {
		HandleError(w, r, err)
	}

	params := store.UploadStoreImageDTO{
		StoreID: storeID,
		Image:   fileBytes,
		Ext:     ext,
	}

	err = xvalidator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeCommand.ChangeStoreProfileImage(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// changeHeaderImage godoc
//
//	@Summary		Store: Change store header image
//	@Description	Change store header image
//	@Tags			Store Management V1
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string				true	"Bearer authentication token"
//	@Param			id				path		string				true	"Store ID"
//	@Param			file			formData	file				true	"Store header image"
//	@Success		200				{object}	nil					"Store activation status updated successfully"
//	@Failure		400				{object}	xerrors.CustomError	"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError	"Unauthorized - authentication required"
//	@Failure		422				{object}	xerrors.CustomError	"Validation error - invalid status"
//	@Failure		500				{object}	xerrors.CustomError	"Internal server error"
//	@Router			/v1/owner/store/{id}/header-image [post]
func (c ownerController) changeHeaderImage(w http.ResponseWriter, r *http.Request) {

	file, _, err := GetFileFormData(w, r, int64(10), "image", []string{"image/jpeg", "image/png"})
	if err != nil {
		HandleError(w, r, err)
		return
	}

	storeID := chi.URLParam(r, "id")

	zt := zitadel.GetInstance()
	authCtx := zt.Middleware.Context(r.Context())

	fileBytes, err := converters.FileToBytes(file)
	if err != nil {
		HandleError(w, r, err)
	}

	params := store.UploadStoreImageDTO{
		StoreID: storeID,
		Image:   fileBytes,
	}

	err = xvalidator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeCommand.ChangeStoreHeaderImage(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// getOwnerStoreByID godoc
//
//	@Summary		Store: Get a store by ID
//	@Description	Get a store by ID
//	@Tags			Store Management V1
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string					true	"Bearer authentication token"
//	@Param			id				path		string					true	"Store ID"
//	@Success		200				{object}	store.QueryOwnerStore	"QueryOwnerStore model"
//	@Failure		400				{object}	xerrors.CustomError		"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError		"Unauthorized - authentication required"
//	@Failure		422				{object}	xerrors.CustomError		"Validation error - invalid status"
//	@Failure		500				{object}	xerrors.CustomError		"Internal server error"
//	@Router			/v1/owner/store/{id} [get]
func (c ownerController) getOwnerStoreByID(w http.ResponseWriter, r *http.Request) {
	stID := chi.URLParam(r, "id")

	st, err := c.storeQuery.GetQueryOwnerStoreByID(r.Context(), stID)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	JSONResponse(w, http.StatusOK, st)
}

// getOwnerStores godoc
//
//	@Summary		Store: Get all stores of an owner
//	@Description	Get all stores of an owner
//	@Tags			Store Management V1
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string						true	"Bearer authentication token"
//	@Success		200				{object}	store.QueryOwnerStoreList	"store model for owner list"
//	@Failure		400				{object}	xerrors.CustomError			"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError			"Unauthorized - authentication required"
//	@Failure		422				{object}	xerrors.CustomError			"Validation error - invalid status"
//	@Failure		500				{object}	xerrors.CustomError			"Internal server error"
//	@Router			/v1/owner/store/all [get]
func (c ownerController) getOwnerStores(w http.ResponseWriter, r *http.Request) {
	zt := zitadel.GetInstance()
	authCtx := zt.Middleware.Context(r.Context())

	st, err := c.storeQuery.GetOwnerStores(r.Context(), authCtx.UserID())
	if err != nil {
		HandleError(w, r, err)
		return
	}
	JSONResponse(w, http.StatusOK, st)
}

func SetupOwnerRoutes(r *chi.Mux, repoFactory persistence.RepositoryFactory, services adapter.Factory) {
	basePath := config.GetInstance().Api.BasePath
	command := owner.NewCommand(repoFactory.NewOwnerRepository())
	stCommand := store.NewCommand(
		repoFactory.NewStoreRepository(),
		repoFactory.NewOwnerRepository(),
		services)
	stQuery := store.NewQueryService(repoFactory.NewStoreRepository(), repoFactory.NewSQLC())
	c := newOwnerController(command, stCommand, stQuery)

	r.Route(basePath+"/v1/owner", func(r chi.Router) {
		r.
			With(middleware.Authentication).
			Post("/", c.createOwner)

		r.Route("/store", func(r chi.Router) {
			r.
				With(middleware.Authorization(string(zitadel.Owner))).
				Post("/", c.createStore)
			r.
				With(middleware.Authorization(string(zitadel.Owner))).
				Put("/", c.updateStore)
			r.
				With(middleware.Authorization(string(zitadel.Owner))).
				Post("/business-hour", c.addBusinessHour)
			r.
				With(middleware.Authorization(string(zitadel.Owner))).
				Delete("/business-hour", c.removeBusinessHour)
			r.
				With(middleware.Authorization(string(zitadel.Owner))).
				Post("/payment-method", c.addPaymentMethod)
			r.
				With(middleware.Authorization(string(zitadel.Owner))).
				Delete("/payment-method", c.removePaymentMethod)
			r.
				With(middleware.Authorization(string(zitadel.Owner))).
				Post("/open", c.openOrClose)
			r.
				With(middleware.Authorization(string(zitadel.Owner))).
				Post("/active", c.activatedOrDeactivate)
			r.
				With(middleware.Authorization(string(zitadel.Owner))).
				Post("/{id}/profile-image", c.changeProfileImage)
			r.
				With(middleware.Authorization(string(zitadel.Owner))).
				Post("/{id}/header-image", c.changeHeaderImage)
			r.
				With(middleware.Authorization(string(zitadel.Owner))).
				Get("/all", c.getOwnerStores)
			r.
				With(middleware.Authorization(string(zitadel.Owner))).
				Get("/{id}", c.getOwnerStoreByID)

		})

	})
}
