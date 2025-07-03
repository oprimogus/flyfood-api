package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oprimogus/flyfood-api/internal/api/middleware"
	"github.com/oprimogus/flyfood-api/internal/config"
	"github.com/oprimogus/flyfood-api/internal/core"
	"github.com/oprimogus/flyfood-api/internal/core/store"
	"github.com/oprimogus/flyfood-api/internal/infrastructure/database/persistence"
	"github.com/oprimogus/flyfood-api/internal/infrastructure/services/adapter"
	_ "github.com/oprimogus/flyfood-api/internal/xerrors"
	"github.com/oprimogus/flyfood-api/internal/xvalidator"
)

type storeController struct {
	commandService store.Command
	queryService   store.Query
}

func newStoreController(command store.Command, query store.Query) storeController {
	return storeController{commandService: command, queryService: query}
}

// getQueryStoreByID godoc
//
//	@Summary		Get a store by ID
//	@Description	Get a store by ID
//	@Tags			Store V1
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string				true	"Bearer authentication token"
//	@Param			id				path		string				true	"Store ID"
//	@Success		200				{object}	store.QueryStore	"Query Store model"
//	@Failure		400				{object}	xerrors.CustomError	"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError	"Unauthorized - authentication required"
//	@Failure		422				{object}	xerrors.CustomError	"Validation error - invalid status"
//	@Failure		500				{object}	xerrors.CustomError	"Internal server error"
//	@Router			/v1/store/{id} [get]
func (c storeController) getQueryStoreByID(w http.ResponseWriter, r *http.Request) {
	stID := chi.URLParam(r, "id")

	st, err := c.queryService.GetQueryStoreByID(r.Context(), stID)
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}
	core.JSONResponse(w, http.StatusOK, st)
}

// getQueryStoreListByFilter godoc
//
//	@Summary		Get a stores by Filter
//	@Description	Get a stores by Filter
//	@Tags			Store V1
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string									true	"Bearer authentication token"
//	@Param			name			query		string									true	"Name"
//	@Param			isOpen			query		string									true	"IsOpen"
//	@Param			score			query		int										true	"Score"
//	@Param			type			query		string									true	"Store Type"
//	@Param			city			query		string									true	"City"
//	@Param			page			query		string									true	"Page"
//	@Param			maxItems		query		string									true	"Items per page"
//	@Success		200				{object}	core.Pagination[store.QueryStoreList]	"Query Store model"
//	@Failure		400				{object}	xerrors.CustomError						"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError						"Unauthorized - authentication required"
//	@Failure		422				{object}	xerrors.CustomError						"Validation error - invalid status"
//	@Failure		500				{object}	xerrors.CustomError						"Internal server error"
//	@Router			/v1/store [get]
func (c storeController) getQueryStoreListByFilter(w http.ResponseWriter, r *http.Request) {
	var params store.QueryStoresInput
	queryParams := r.URL.Query()

	name := queryParams.Get("name")
	if name != "" {
		params.Name = &name
	}

	city := queryParams.Get("city")
	if city != "" {
		params.City = &city
	}

	storeType := queryParams.Get("type")
	if storeType != "" {
		if store.IsValidType(storeType) {
			value := store.Type(storeType)
			params.Type = &value
		} else {
			xerror := xvalidator.NewFieldError("type", storeType)
			core.HandleApiError(w, r, xerror)
			return
		}
	}

	isOpenParam := queryParams.Get("isOpen")
	if isOpenParam != "" {
		isOpen, err := core.IsValidBool(isOpenParam)
		if err != nil {
			xerror := xvalidator.NewFieldError("isOpen", isOpenParam)
			core.HandleApiError(w, r, xerror)
			return
		} else {
			params.IsOpen = &isOpen
		}
	}

	scoreParam := queryParams.Get("score")
	if scoreParam != "" {
		score, err := core.IsValidInt(scoreParam)
		if err != nil {
			xerror := xvalidator.NewFieldError("score", scoreParam)
			core.HandleApiError(w, r, xerror)
			return
		} else {
			params.Score = &score
		}
	}

	pageParam := queryParams.Get("page")
	page, err := core.IsValidInt(pageParam)
	if err != nil {
		xerror := xvalidator.NewFieldError("page", pageParam)
		core.HandleApiError(w, r, xerror)
		return
	} else {
		params.Page = page
	}

	maxItemsParam := queryParams.Get("maxItems")
	maxItems, err := core.IsValidInt(maxItemsParam)
	if err != nil {
		xerror := xvalidator.NewFieldError("maxItems", maxItemsParam)
		core.HandleApiError(w, r, xerror)
		return
	} else {
		if maxItems > 50 {
			params.MaxItems = 50
		} else {
			params.MaxItems = maxItems
		}
	}

	st, err := c.queryService.GetStoreByFilter(r.Context(), params)
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}
	core.JSONResponse(w, http.StatusOK, st)
}

func SetupStoreRoutes(r *chi.Mux, repoFactory persistence.RepositoryFactory, services adapter.Factory) {
	basePath := config.GetInstance().Api.BasePath
	command := store.NewCommand(
		repoFactory.NewStoreRepository(),
		repoFactory.NewOwnerRepository(),
		services)
	query := store.NewQueryService(repoFactory.NewStoreRepository(), repoFactory.NewSQLC())
	c := newStoreController(command, query)

	r.Route(basePath+"/v1", func(r chi.Router) {
		r.
			With(middleware.Authentication).
			Get("/store/{id}", c.getQueryStoreByID)
		r.
			With(middleware.Authentication).
			Get("/store", c.getQueryStoreListByFilter)
	})
}
