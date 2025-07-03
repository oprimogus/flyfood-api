package customer

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oprimogus/flyfood-api/internal/core"
	"github.com/oprimogus/flyfood-api/internal/api/middleware"
	"github.com/oprimogus/flyfood-api/internal/config"
	"github.com/oprimogus/flyfood-api/internal/core/address"
	"github.com/oprimogus/flyfood-api/internal/infrastructure/services/zitadel"
	"github.com/oprimogus/flyfood-api/internal/xvalidator"
	postgresDB "github.com/oprimogus/flyfood-api/internal/infrastructure/database/postgres"
)

type Controller struct {
	service Service
}

func NewController(customerService Service) Controller {
	return Controller{service: customerService}
}

// getCustomer godoc
//
//	@Summary		Get a customer account
//	@Description	Register a comprehensive customer profile with full registration details
//	@Tags			Customer Profile Management V1
//	@Produce		json
//	@Success		200	{object}	customer.Customer	"Customer"
//	@Success		201	{object}	customer.Customer	"Customer"
//	@Failure		400	{object}	xerrors.CustomError	"Invalid request data or malformed JSON"
//	@Failure		409	{object}	xerrors.CustomError	"Conflict - customer may already exist"
//	@Failure		422	{object}	xerrors.CustomError	"Validation error - invalid customer details"
//	@Failure		500	{object}	xerrors.CustomError	"Internal server error"
//	@Failure		502	{object}	xerrors.CustomError	"External service communication error"
//	@Router			/v1/customer [get]
func (c Controller) getCustomer(w http.ResponseWriter, r *http.Request) {
	zt := zitadel.GetInstance()
	authCtx := zt.GetContext(r.Context())

	cs, err := c.service.FindCustomer(r.Context(), authCtx.UserID())
	if err != nil {
		params := CreateProfileDTO{
			ID:       authCtx.UserID(),
			Name:     authCtx.GivenName,
			Email:    authCtx.Email,
			LastName: authCtx.FamilyName,
			Phone:    authCtx.PhoneNumber,
		}

		err = xvalidator.Validate(params)
		if err != nil {
			core.HandleApiError(w, r, err)
			return
		}

		createdCustomer, err := c.service.CreateCustomer(r.Context(), params)
		if err != nil {
			core.HandleApiError(w, r, err)
			return
		}

		core.JSONResponse(w, http.StatusCreated, createdCustomer)
		return
	}

	core.JSONResponse(w, http.StatusOK, cs)

}

// updateProfile godoc
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
func (c Controller) updateProfile(w http.ResponseWriter, r *http.Request) {
	zt := zitadel.GetInstance()
	authCtx := zt.Middleware.Context(r.Context())

	var params UpdateProfileDTO
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}

	err = xvalidator.Validate(params)
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}

	err = c.service.UpdateCustomerProfile(r.Context(), authCtx.UserID(), params)
	if err != nil {
		core.HandleApiError(w, r, err)
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
func (c Controller) addAddress(w http.ResponseWriter, r *http.Request) {
	zt := zitadel.GetInstance()
	authCtx := zt.Middleware.Context(r.Context())

	var params address.Address
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}

	err = xvalidator.Validate(params)
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}

	err = c.service.AddAddress(r.Context(), authCtx.UserID(), params)
	if err != nil {
		core.HandleApiError(w, r, err)
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
func (c Controller) removeAddress(w http.ResponseWriter, r *http.Request) {
	zt := zitadel.GetInstance()
	authCtx := zt.Middleware.Context(r.Context())

	var params address.Address
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}

	err = xvalidator.Validate(params)
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}

	err = c.service.RemoveAddress(r.Context(), authCtx.UserID(), params)
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func SetupRoutes(r *chi.Mux, db *postgresDB.Database) {
	basePath := config.GetInstance().Api.BasePath
	service := NewService(NewCustomerRepository(db))
	c := NewController(service)

	r.Route(basePath+"/v1/customer", func(r chi.Router) {
		r.
			With(middleware.Authentication).
			Get("/", c.getCustomer)
		r.
			With(middleware.Authentication).
			Put("/", c.updateProfile)
		r.
			With(middleware.Authentication).
			Post("/address", c.addAddress)
		r.
			With(middleware.Authentication).
			Delete("/address", c.removeAddress)
	})
}
