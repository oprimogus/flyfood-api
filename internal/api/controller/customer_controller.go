package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/oprimogus/cardapiogo/internal/api/middleware"
	"github.com/oprimogus/cardapiogo/internal/config"
	"github.com/oprimogus/cardapiogo/internal/core/address"
	"github.com/oprimogus/cardapiogo/internal/core/customer"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/database/persistence"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/zitadel"
	"github.com/oprimogus/cardapiogo/internal/xvalidator"
	"net/http"
)

type customerController struct {
	validator       *xvalidator.Validator
	customerService customer.Service
}

func newCustomerController(validator *xvalidator.Validator, customerService customer.Service) customerController {
	return customerController{validator: validator, customerService: customerService}
}

// getCustomer godoc
//
//	@Summary		Get a customer account
//	@Description	Register a comprehensive customer profile with full registration details
//	@Tags			Customer Profile Management V1
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	customer.Customer	"Customer"
//	@Success		201	{object}	customer.Customer	"Customer"
//	@Failure		400	{object}	xerrors.CustomError	"Invalid request data or malformed JSON"
//	@Failure		409	{object}	xerrors.CustomError	"Conflict - customer may already exist"
//	@Failure		422	{object}	xerrors.CustomError	"Validation error - invalid customer details"
//	@Failure		500	{object}	xerrors.CustomError	"Internal server error"
//	@Failure		502	{object}	xerrors.CustomError	"External service communication error"
//	@Router			/v1/customer [get]
func (c customerController) getCustomer(w http.ResponseWriter, r *http.Request) {
	zt := zitadel.GetInstance()
	authCtx := zt.GetContext(r.Context())

	cs, err := c.customerService.FindCustomer(r.Context(), authCtx.UserID())
	if err != nil {
		params := customer.CreateProfileDTO{
			ID:       authCtx.UserID(),
			Name:     authCtx.GivenName,
			Email:    authCtx.Email,
			LastName: authCtx.FamilyName,
			Phone:    authCtx.PhoneNumber,
		}

		err = c.validator.Validate(params)
		if err != nil {
			HandleError(w, r, err)
			return
		}

		createdCustomer, err := c.customerService.CreateCustomer(r.Context(), params)
		if err != nil {
			HandleError(w, r, err)
			return
		}

		JSONResponse(w, http.StatusCreated, createdCustomer)
		return
	}

	JSONResponse(w, http.StatusOK, cs)

}

// updateCustomerProfile godoc
//
//	@Summary		Update existing customer profile
//	@Description	Modify comprehensive customer account details with full profile update
//	@Tags			Customer Profile Management V1
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string						true	"Bearer authentication token"
//	@Param			request			body		customer.UpdateProfileDTO	true	"Updated customer profile information"
//	@Success		200				{object}	nil							"Profile successfully updated"
//	@Failure		400				{object}	xerrors.CustomError			"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError			"Unauthorized - authentication required"
//	@Failure		403				{object}	xerrors.CustomError			"Forbidden - insufficient permissions"
//	@Failure		404				{object}	xerrors.CustomError			"Customer profile not found"
//	@Failure		422				{object}	xerrors.CustomError			"Validation error - invalid profile details"
//	@Failure		500				{object}	xerrors.CustomError			"Internal server error"
//	@Router			/v1/customer [put]
func (c customerController) updateCustomerProfile(w http.ResponseWriter, r *http.Request) {
	zt := zitadel.GetInstance()
	authCtx := zt.Middleware.Context(r.Context())

	var params customer.UpdateProfileDTO
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

	err = c.customerService.UpdateCustomerProfile(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// addAddress godoc
//
//	@Summary		Add new address to customer
//	@Description	Add new address to customer
//	@Tags			Customer Profile Management V1
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string				true	"Bearer authentication token"
//	@Param			request			body		address.Address		true	"Address Model"
//	@Success		200				{object}	nil					"Address added to user"
//	@Failure		400				{object}	xerrors.CustomError	"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError	"Unauthorized - authentication required"
//	@Failure		403				{object}	xerrors.CustomError	"Forbidden - insufficient permissions"
//	@Failure		404				{object}	xerrors.CustomError	"Customer profile not found"
//	@Failure		422				{object}	xerrors.CustomError	"Validation error - invalid profile details"
//	@Failure		500				{object}	xerrors.CustomError	"Internal server error"
//	@Router			/v1/customer/address [post]
func (c customerController) addAddress(w http.ResponseWriter, r *http.Request) {
	zt := zitadel.GetInstance()
	authCtx := zt.Middleware.Context(r.Context())

	var params address.Address
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

	err = c.customerService.AddAddress(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// removeAddress godoc
//
//	@Summary		Remove an address to customer
//	@Description	Remove an address to customer
//	@Tags			Customer Profile Management V1
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string				true	"Bearer authentication token"
//	@Param			request			body		address.Address		true	"Address Model"
//	@Success		200				{object}	nil					"Address removed from user"
//	@Failure		400				{object}	xerrors.CustomError	"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError	"Unauthorized - authentication required"
//	@Failure		403				{object}	xerrors.CustomError	"Forbidden - insufficient permissions"
//	@Failure		404				{object}	xerrors.CustomError	"Customer profile not found"
//	@Failure		422				{object}	xerrors.CustomError	"Validation error - invalid profile details"
//	@Failure		500				{object}	xerrors.CustomError	"Internal server error"
//	@Router			/v1/customer/address [delete]
func (c customerController) removeAddress(w http.ResponseWriter, r *http.Request) {
	zt := zitadel.GetInstance()
	authCtx := zt.Middleware.Context(r.Context())

	var params address.Address
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

	err = c.customerService.RemoveAddress(r.Context(), authCtx.UserID(), params)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func SetupCustomerRoutes(r *chi.Mux, repoFactory persistence.RepositoryFactory) {
	basePath := config.GetInstance().Api.BasePath

	validator := xvalidator.GetPtInstance()

	service := customer.NewService(repoFactory.NewCustomerRepository())
	c := newCustomerController(validator, service)

	r.Route(basePath+"/v1/customer", func(r chi.Router) {
		r.
			With(middleware.Authentication).
			Get("/", c.getCustomer)
		r.
			With(middleware.Authentication).
			Put("/", c.updateCustomerProfile)
		r.
			With(middleware.Authentication).
			Post("/address", c.addAddress)
		r.
			With(middleware.Authentication).
			Delete("/address", c.removeAddress)
	})
}
