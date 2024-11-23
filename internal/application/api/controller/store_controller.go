package controller

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/oprimogus/cardapiogo/internal/application/store"
	"github.com/oprimogus/cardapiogo/internal/config"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/database/persistence"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/adapter"
	"github.com/oprimogus/cardapiogo/internal/xvalidator"
	"github.com/oprimogus/cardapiogo/pkg/converters"
	"net/http"
)

type storeController struct {
	validator    *xvalidator.Validator
	storeService store.Service
}

func newStoreController(validator *xvalidator.Validator, storeService store.Service) storeController {
	return storeController{validator: validator, storeService: storeService}
}

// createStore godoc
// @Summary     Create a new store for an owner
// @Description Register a comprehensive store profile with all necessary details
// @Tags        Owner
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer authentication token"
// @Param       request body     store.CreateNewStoreDTO true "Detailed store creation information"
// @Success     201 {object}     nil "Store successfully created"
// @Failure     400 {object}     xerrors.CustomError "Invalid request data or malformed JSON"
// @Failure     401 {object}     xerrors.CustomError "Unauthorized - authentication required"
// @Failure     409 {object}     xerrors.CustomError "Conflict - store may already exist"
// @Failure     422 {object}     xerrors.CustomError "Validation error - invalid store details"
// @Failure     500 {object}     xerrors.CustomError "Internal server error"
// @Failure     502 {object}     xerrors.CustomError "External service communication error"
// @Router      /v1/owner/store [post]
func (c storeController) createStore(w http.ResponseWriter, r *http.Request) {
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

	err = c.storeService.CreateNewStore(r.Context(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// updateStore godoc
// @Summary     Update existing store profile
// @Description Modify comprehensive store details including contact, operational, and profile information
// @Tags        Store Management
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer authentication token"
// @Param       request body     store.UpdateStoreDTO true "Updated store information"
// @Success     200 {object}     nil "Store successfully updated"
// @Failure     400 {object}     xerrors.CustomError "Invalid request data or malformed JSON"
// @Failure     401 {object}     xerrors.CustomError "Unauthorized - authentication required"
// @Failure     404 {object}     xerrors.CustomError "Store not found"
// @Failure     422 {object}     xerrors.CustomError "Validation error - invalid store details"
// @Failure     500 {object}     xerrors.CustomError "Internal server error"
// @Router      /v1/store [put]
func (c storeController) updateStore(w http.ResponseWriter, r *http.Request) {
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

	err = c.storeService.Update(r.Context(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// addBusinessHour godoc
// @Summary     Add business operating hours for store
// @Description Register specific business hours for a store's operational schedule
// @Tags        Store Management
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer authentication token"
// @Param       request body     store.AddOrDeleteBusinessHourDTO true "Business hours details"
// @Success     200 {object}     nil "Business hours successfully added"
// @Failure     400 {object}     xerrors.CustomError "Invalid request data or malformed JSON"
// @Failure     401 {object}     xerrors.CustomError "Unauthorized - authentication required"
// @Failure     422 {object}     xerrors.CustomError "Validation error - invalid business hours"
// @Failure     500 {object}     xerrors.CustomError "Internal server error"
// @Router      /v1/store/business-hour [post]
func (c storeController) addBusinessHour(w http.ResponseWriter, r *http.Request) {
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

	err = c.storeService.AddBusinessHour(r.Context(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// removeBusinessHour godoc
// @Summary     Remove business operating hours for store
// @Description Delete specific business hours from store's operational schedule
// @Tags        Store Management
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer authentication token"
// @Param       request body     store.AddOrDeleteBusinessHourDTO true "Business hours to remove"
// @Success     200 {object}     nil "Business hours successfully removed"
// @Failure     400 {object}     xerrors.CustomError "Invalid request data or malformed JSON"
// @Failure     401 {object}     xerrors.CustomError "Unauthorized - authentication required"
// @Failure     422 {object}     xerrors.CustomError "Validation error - invalid business hours"
// @Failure     500 {object}     xerrors.CustomError "Internal server error"
// @Router      /v1/store/business-hour [delete]
func (c storeController) removeBusinessHour(w http.ResponseWriter, r *http.Request) {
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

	err = c.storeService.RemoveBusinessHour(r.Context(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// addPaymentMethod godoc
// @Summary     Add payment method for store
// @Description Register a new payment method accepted by the store
// @Tags        Store Management
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer authentication token"
// @Param       request body     store.AddOrDeletePaymentMethodDTO true "Payment method details"
// @Success     200 {object}     nil "Payment method successfully added"
// @Failure     400 {object}     xerrors.CustomError "Invalid request data or malformed JSON"
// @Failure     401 {object}     xerrors.CustomError "Unauthorized - authentication required"
// @Failure     422 {object}     xerrors.CustomError "Validation error - invalid payment method"
// @Failure     500 {object}     xerrors.CustomError "Internal server error"
// @Router      /v1/store/payment-method [post]
func (c storeController) addPaymentMethod(w http.ResponseWriter, r *http.Request) {
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

	err = c.storeService.AddPaymentMethod(r.Context(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// removePaymentMethod godoc
// @Summary     Remove payment method for store
// @Description Delete a previously added payment method from store's accepted methods
// @Tags        Store Management
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer authentication token"
// @Param       request body     store.AddOrDeletePaymentMethodDTO true "Payment method to remove"
// @Success     200 {object}     nil "Payment method successfully removed"
// @Failure     400 {object}     xerrors.CustomError "Invalid request data or malformed JSON"
// @Failure     401 {object}     xerrors.CustomError "Unauthorized - authentication required"
// @Failure     422 {object}     xerrors.CustomError "Validation error - invalid payment method"
// @Failure     500 {object}     xerrors.CustomError "Internal server error"
// @Router      /v1/store/payment-method [delete]
func (c storeController) removePaymentMethod(w http.ResponseWriter, r *http.Request) {
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

	err = c.storeService.RemovePaymentMethod(r.Context(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// openOrClose godoc
// @Summary     Change store's operational status for orders
// @Description Toggle store's availability to accept or stop accepting new orders
// @Tags        Store Management
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer authentication token"
// @Param       request body     store.SetOpenStateDTO true "Store open/close status"
// @Success     200 {object}     nil "Store operational status updated successfully"
// @Failure     400 {object}     xerrors.CustomError "Invalid request data or malformed JSON"
// @Failure     401 {object}     xerrors.CustomError "Unauthorized - authentication required"
// @Failure     422 {object}     xerrors.CustomError "Validation error - invalid status"
// @Failure     500 {object}     xerrors.CustomError "Internal server error"
// @Router      /v1/store/open [post]
func (c storeController) openOrClose(w http.ResponseWriter, r *http.Request) {
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

	err = c.storeService.OpenOrCloseStore(r.Context(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// activatedOrDeactivate godoc
// @Summary     Change store's visibility and order acceptance status
// @Description Activate or deactivate store to control product visibility and order processing
// @Tags        Store Management
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer authentication token"
// @Param       request body     store.SetActiveDTO true "Store activation status"
// @Success     200 {object}    nil "Store activation status updated successfully"
// @Failure     400 {object}     xerrors.CustomError "Invalid request data or malformed JSON"
// @Failure     401 {object}     xerrors.CustomError "Unauthorized - authentication required"
// @Failure     422 {object}     xerrors.CustomError "Validation error - invalid status"
// @Failure     500 {object}     xerrors.CustomError "Internal server error"
// @Router      /v1/store/active [post]
func (c storeController) activatedOrDeactivate(w http.ResponseWriter, r *http.Request) {
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

	err = c.storeService.ActiveOrDeactivateStore(r.Context(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// changeProfileImage godoc
// @Summary     Change store profile image
// @Description Change store profile image
// @Tags        Store Management
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer authentication token"
// @Param       request body     store.SetActiveDTO true "Store activation status"
// @Success     200 {object}    nil "Store activation status updated successfully"
// @Failure     400 {object}     xerrors.CustomError "Invalid request data or malformed JSON"
// @Failure     401 {object}     xerrors.CustomError "Unauthorized - authentication required"
// @Failure     422 {object}     xerrors.CustomError "Validation error - invalid status"
// @Failure     500 {object}     xerrors.CustomError "Internal server error"
// @Router      /v1/store/{storeID}/active [post]
func (c storeController) changeProfileImage(w http.ResponseWriter, r *http.Request) {

	file, _, err := GetFileFormData(w, r, int64(10), "image", []mimeType{JPEG, PNG})
	if err != nil {
		HandleError(w, r, err)
		return
	}

	storeID := chi.URLParam(r, "storeID")
	ownerID := "" // TODO: Get from Zitadel context
	fileBytes, err := converters.FileToBytes(file)
	if err != nil {
		HandleError(w, r, err)
	}

	params := store.UploadStoreImageDTO{
		StoreID: storeID,
		OwnerID: ownerID,
		Image:   fileBytes,
	}

	err = c.validator.Validate(params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = c.storeService.ChangeProfileImage(r.Context(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func SetupStoreRoutes(r *chi.Mux, repoFactory persistence.RepositoryFactory, services adapter.Factory) {
	basePath := config.GetInstance().Api.BasePath

	validator := xvalidator.GetPtInstance()

	service := store.NewService(repoFactory.NewStoreRepository())
	c := newStoreController(validator, service)

	r.Route(basePath+"/v1/owner", func(r chi.Router) {
		r.Post("/store", c.createStore)
	})
	r.Route(basePath+"/v1/store", func(r chi.Router) {
		r.Put("/", c.updateStore)
		r.Post("/business-hour", c.addBusinessHour)
		r.Delete("/business-hour", c.removeBusinessHour)
		r.Post("/payment-method", c.addPaymentMethod)
		r.Delete("/payment-method", c.removePaymentMethod)
		r.Post("/open", c.openOrClose)
		r.Post("/active", c.activatedOrDeactivate)
	})
}
