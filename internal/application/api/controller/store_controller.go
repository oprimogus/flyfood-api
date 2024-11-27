package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/oprimogus/cardapiogo/internal/application/api/middleware"
	"github.com/oprimogus/cardapiogo/internal/application/store"
	"github.com/oprimogus/cardapiogo/internal/config"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/database/persistence"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/adapter"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/zitadel"
	"github.com/oprimogus/cardapiogo/internal/xvalidator"
	"github.com/oprimogus/cardapiogo/pkg/converters"
	"net/http"
	"strconv"
)

type storeController struct {
	validator    *xvalidator.Validator
	storeService store.Service
}

func newStoreController(validator *xvalidator.Validator, storeService store.Service) storeController {
	return storeController{validator: validator, storeService: storeService}
}

// createOwner godoc
//
//	@Summary		Create an owner
//	@Description	Register a customer as owner
//	@Tags			Owner
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string					true	"Bearer authentication token"
//	@Success		201				{object}	nil						"Store successfully created"
//	@Failure		400				{object}	xerrors.CustomError		"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError		"Unauthorized - authentication required"
//	@Failure		409				{object}	xerrors.CustomError		"Conflict - store may already exist"
//	@Failure		422				{object}	xerrors.CustomError		"Validation error - invalid store details"
//	@Failure		500				{object}	xerrors.CustomError		"Internal server error"
//	@Failure		502				{object}	xerrors.CustomError		"External service communication error"
//	@Router			/v1/owner [post]
func (c storeController) createOwner(w http.ResponseWriter, r *http.Request) {
	zt := zitadel.GetInstance()
	authCtx := zt.GetContext(r.Context())
	userID, err := strconv.Atoi(authCtx.UserID())
	if err != nil {
		HandleError(w, r, err)
		return
	}
	err = c.storeService.NewOwner(r.Context(), userID)
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
//	@Summary		Create a new store for an owner
//	@Description	Register a comprehensive store profile with all necessary details
//	@Tags			Owner
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
func (c storeController) createStore(w http.ResponseWriter, r *http.Request) {
	authCtx := zitadel.GetInstance().GetContext(r.Context())
	var params store.CreateNewStoreDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.validator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeService.CreateNewStore(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// updateStore godoc
//
//	@Summary		Update existing store profile
//	@Description	Modify comprehensive store details including contact, operational, and profile information
//	@Tags			Store Management
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
//	@Router			/v1/store [put]
func (c storeController) updateStore(w http.ResponseWriter, r *http.Request) {
	authCtx := zitadel.GetInstance().GetContext(r.Context())
	var params store.UpdateStoreDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.validator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeService.Update(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// addBusinessHour godoc
//
//	@Summary		Add business operating hours for store
//	@Description	Register specific business hours for a store's operational schedule
//	@Tags			Store Management
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
//	@Router			/v1/store/business-hour [post]
func (c storeController) addBusinessHour(w http.ResponseWriter, r *http.Request) {
	authCtx := zitadel.GetInstance().GetContext(r.Context())
	var params store.AddOrDeleteBusinessHourDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.validator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeService.AddBusinessHour(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// removeBusinessHour godoc
//
//	@Summary		Remove business operating hours for store
//	@Description	Delete specific business hours from store's operational schedule
//	@Tags			Store Management
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
//	@Router			/v1/store/business-hour [delete]
func (c storeController) removeBusinessHour(w http.ResponseWriter, r *http.Request) {
	authCtx := zitadel.GetInstance().GetContext(r.Context())
	var params store.AddOrDeleteBusinessHourDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.validator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeService.RemoveBusinessHour(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// addPaymentMethod godoc
//
//	@Summary		Add payment method for store
//	@Description	Register a new payment method accepted by the store
//	@Tags			Store Management
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
//	@Router			/v1/store/payment-method [post]
func (c storeController) addPaymentMethod(w http.ResponseWriter, r *http.Request) {
	authCtx := zitadel.GetInstance().GetContext(r.Context())
	var params store.AddOrDeletePaymentMethodDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.validator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeService.AddPaymentMethod(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// removePaymentMethod godoc
//
//	@Summary		Remove payment method for store
//	@Description	Delete a previously added payment method from store's accepted methods
//	@Tags			Store Management
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
//	@Router			/v1/store/payment-method [delete]
func (c storeController) removePaymentMethod(w http.ResponseWriter, r *http.Request) {
	authCtx := zitadel.GetInstance().GetContext(r.Context())
	var params store.AddOrDeletePaymentMethodDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.validator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeService.RemovePaymentMethod(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// openOrClose godoc
//
//	@Summary		Change store's operational status for orders
//	@Description	Toggle store's availability to accept or stop accepting new orders
//	@Tags			Store Management
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
//	@Router			/v1/store/open [post]
func (c storeController) openOrClose(w http.ResponseWriter, r *http.Request) {
	authCtx := zitadel.GetInstance().GetContext(r.Context())
	var params store.SetOpenStateDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.validator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeService.OpenOrCloseStore(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// activatedOrDeactivate godoc
//
//	@Summary		Change store's visibility and order acceptance status
//	@Description	Activate or deactivate store to control product visibility and order processing
//	@Tags			Store Management
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
//	@Router			/v1/store/active [post]
func (c storeController) activatedOrDeactivate(w http.ResponseWriter, r *http.Request) {
	authCtx := zitadel.GetInstance().GetContext(r.Context())
	var params store.SetActiveDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.validator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeService.ActiveOrDeactivateStore(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// changeProfileImage godoc
//
//	@Summary		Change store profile image
//	@Description	Change store profile image
//	@Tags			Store Management
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
func (c storeController) changeProfileImage(w http.ResponseWriter, r *http.Request) {
	file, _, err := GetFileFormData(w, r, int64(10), "image", []string{"image/jpeg", "image/png", "image/jpg"})
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

	err = c.validator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeService.ChangeProfileImage(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// changeHeaderImage godoc
//
//	@Summary		Change store profile image
//	@Description	Change store profile image
//	@Tags			Store Management
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
//	@Router			/v1/store/{id}/header-image [post]
func (c storeController) changeHeaderImage(w http.ResponseWriter, r *http.Request) {

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

	err = c.validator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeService.ChangeHeaderImage(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func SetupStoreRoutes(r *chi.Mux, repoFactory persistence.RepositoryFactory, services adapter.Factory) {
	basePath := config.GetInstance().Api.BasePath

	validator := xvalidator.GetPtInstance()

	service := store.NewService(
		repoFactory.NewStoreRepository(),
		repoFactory.NewCustomerRepository(),
		services)
	c := newStoreController(validator, service)

	r.Route(basePath+"/v1/owner", func(r chi.Router) {
		r.
			With(middleware.Authentication).
			Post("/", c.createOwner)

		r.Route("/store", func(r chi.Router) {
			r.
				With(middleware.Authorization("owner")).
				Post("/", c.createStore)
			r.
				With(middleware.Authorization("owner")).
				Put("/", c.updateStore)
			r.
				With(middleware.Authorization("owner")).
				Post("/business-hour", c.addBusinessHour)
			r.
				With(middleware.Authorization("owner")).
				Delete("/business-hour", c.removeBusinessHour)
			r.
				With(middleware.Authorization("owner")).
				Post("/payment-method", c.addPaymentMethod)
			r.
				With(middleware.Authorization("owner")).
				Delete("/payment-method", c.removePaymentMethod)
			r.
				With(middleware.Authorization("owner")).
				Post("/open", c.openOrClose)
			r.
				With(middleware.Authorization("owner")).
				Post("/active", c.activatedOrDeactivate)
			r.
				With(middleware.Authorization("owner")).
				Post("/{id}/profile-image", c.changeProfileImage)
			r.
				With(middleware.Authorization("owner")).
				Post("/{id}/header-image", c.changeHeaderImage)
		})
	})

}
