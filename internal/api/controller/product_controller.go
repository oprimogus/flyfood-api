package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oprimogus/flyfood-api/internal/api/middleware"
	"github.com/oprimogus/flyfood-api/internal/core"
	"github.com/oprimogus/flyfood-api/internal/config"
	"github.com/oprimogus/flyfood-api/internal/core/store"
	"github.com/oprimogus/flyfood-api/internal/core/store/product"
	"github.com/oprimogus/flyfood-api/internal/infrastructure/database/persistence"
	"github.com/oprimogus/flyfood-api/internal/infrastructure/services/adapter"
	"github.com/oprimogus/flyfood-api/internal/infrastructure/services/zitadel"
	_ "github.com/oprimogus/flyfood-api/internal/xerrors"
	"github.com/oprimogus/flyfood-api/internal/xvalidator"
	"github.com/oprimogus/flyfood-api/pkg/converters"
)

type productController struct {
	command   store.Command
}

func newProductController(command store.Command) productController {
	return productController{command: command}
}

// addNewProduct godoc
//
//	@Summary		Product: Add a new product into store
//	@Description	Add a new product into store
//	@Tags			Store Management V1
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string						true	"Bearer authentication token"
//	@Param			request			body		product.CreateProductDTO	true	"Detailed product creation information"
//	@Success		201				{object}	nil							"Product successfully created"
//	@Failure		400				{object}	xerrors.CustomError			"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError			"Unauthorized - authentication required"
//	@Failure		409				{object}	xerrors.CustomError			"Conflict - store may already exist"
//	@Failure		422				{object}	xerrors.CustomError			"Validation error - invalid store details"
//	@Failure		500				{object}	xerrors.CustomError			"Internal server error"
//	@Failure		502				{object}	xerrors.CustomError			"External service communication error"
//	@Router			/v1/owner/store/product [post]
func (c productController) addNewProduct(w http.ResponseWriter, r *http.Request) {
	zt := zitadel.GetInstance()
	authCtx := zt.GetContext(r.Context())

	var params product.CreateProductDTO
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

	err = c.command.NewStoreProduct(r.Context(), authCtx.UserID(), params)
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// updateProduct godoc
//
//	@Summary		Product: Update a product of store
//	@Description	Update a product of store
//	@Tags			Store Management V1
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string						true	"Bearer authentication token"
//	@Param			request			body		product.UpdateProductDTO	true	"Detailed product information"
//	@Success		200				{object}	nil							"Stock successfully updated"
//	@Failure		400				{object}	xerrors.CustomError			"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError			"Unauthorized - authentication required"
//	@Failure		409				{object}	xerrors.CustomError			"Conflict - store may already exist"
//	@Failure		422				{object}	xerrors.CustomError			"Validation error - invalid store details"
//	@Failure		500				{object}	xerrors.CustomError			"Internal server error"
//	@Failure		502				{object}	xerrors.CustomError			"External service communication error"
//	@Router			/v1/owner/store/product [put]
func (c productController) updateProduct(w http.ResponseWriter, r *http.Request) {
	zt := zitadel.GetInstance()
	authCtx := zt.GetContext(r.Context())

	var params product.UpdateProductDTO
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

	err = c.command.UpdateProduct(r.Context(), authCtx.UserID(), params)
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// increaseProductStock godoc
//
//	@Summary		Product: Increase stock of a product
//	@Description	Increase stock of a product
//	@Tags			Store Management V1
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string							true	"Bearer authentication token"
//	@Param			request			body		product.ChangeStockProductDTO	true	"Detailed product information"
//	@Success		200				{object}	nil								"Stock successfully updated"
//	@Failure		400				{object}	xerrors.CustomError				"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError				"Unauthorized - authentication required"
//	@Failure		409				{object}	xerrors.CustomError				"Conflict - store may already exist"
//	@Failure		422				{object}	xerrors.CustomError				"Validation error - invalid store details"
//	@Failure		500				{object}	xerrors.CustomError				"Internal server error"
//	@Failure		502				{object}	xerrors.CustomError				"External service communication error"
//	@Router			/v1/owner/store/product/increase-stock [post]
func (c productController) increaseProductStock(w http.ResponseWriter, r *http.Request) {
	zt := zitadel.GetInstance()
	authCtx := zt.GetContext(r.Context())

	var params product.ChangeStockProductDTO
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

	err = c.command.IncreaseStock(r.Context(), authCtx.UserID(), params)
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// decreaseProductStock godoc
//
//	@Summary		Product: Decrease stock of a product
//	@Description	Decrease stock of a product
//	@Tags			Store Management V1
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string							true	"Bearer authentication token"
//	@Param			request			body		product.ChangeStockProductDTO	true	"Detailed product information"
//	@Success		200				{object}	nil								"Stock successfully updated"
//	@Failure		400				{object}	xerrors.CustomError				"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError				"Unauthorized - authentication required"
//	@Failure		409				{object}	xerrors.CustomError				"Conflict - store may already exist"
//	@Failure		422				{object}	xerrors.CustomError				"Validation error - invalid store details"
//	@Failure		500				{object}	xerrors.CustomError				"Internal server error"
//	@Failure		502				{object}	xerrors.CustomError				"External service communication error"
//	@Router			/v1/owner/store/product/decrease-stock [post]
func (c productController) decreaseProductStock(w http.ResponseWriter, r *http.Request) {
	zt := zitadel.GetInstance()
	authCtx := zt.GetContext(r.Context())

	var params product.ChangeStockProductDTO
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

	err = c.command.DecreaseStock(r.Context(), authCtx.UserID(), params)
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// changeProductImage godoc
//
//	@Summary		Product: Change product image
//	@Description	Change product image
//	@Tags			Store Management V1
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string				true	"Bearer authentication token"
//	@Param			store_id		path		string				true	"Store ID"
//	@Param			product_id		path		string				true	"Product ID"
//	@Param			file			formData	file				true	"Product Image"
//	@Success		200				{object}	nil					"Product image updated with success"
//	@Failure		400				{object}	xerrors.CustomError	"Invalid request data or malformed JSON"
//	@Failure		401				{object}	xerrors.CustomError	"Unauthorized - authentication required"
//	@Failure		422				{object}	xerrors.CustomError	"Validation error - invalid status"
//	@Failure		500				{object}	xerrors.CustomError	"Internal server error"
//	@Router			/v1/store/{store_id}/product/{product_id}/header-image [post]
func (c productController) changeProductImage(w http.ResponseWriter, r *http.Request) {
	file, _, err := core.GetFileFormData(w, r, int64(10), "image", []string{"image/jpeg", "image/png"})
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}

	storeID := chi.URLParam(r, "store_id")
	productID := chi.URLParam(r, "product_id")

	zt := zitadel.GetInstance()
	authCtx := zt.Middleware.Context(r.Context())

	fileBytes, err := converters.FileToBytes(file)
	if err != nil {
		core.HandleApiError(w, r, err)
	}

	params := product.UploadProductImageDTO{
		StoreID:   storeID,
		ProductID: productID,
		Image:     fileBytes,
	}

	err = xvalidator.Validate(params)
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}

	err = c.command.ChangeProductImage(r.Context(), authCtx.UserID(), params)
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func SetupProductRoutes(r *chi.Mux, repoFactory persistence.RepositoryFactory, services adapter.Factory) {
	basePath := config.GetInstance().Api.BasePath
	command := store.NewCommand(
		repoFactory.NewStoreRepository(),
		repoFactory.NewOwnerRepository(),
		services)

	c := newProductController(command)

	r.
		With(middleware.Authorization(string(zitadel.Owner))).
		Post(basePath+"/v1/owner/store/{store_id}/product/{product_id}/image", c.changeProductImage)

	r.Route(basePath+"/v1/owner/store/product", func(r chi.Router) {
		r.
			With(middleware.Authorization(string(zitadel.Owner))).
			Post("/", c.addNewProduct)
		r.
			With(middleware.Authorization(string(zitadel.Owner))).
			Put("/", c.updateProduct)
		r.
			With(middleware.Authorization(string(zitadel.Owner))).
			Post("/increase-stock", c.increaseProductStock)
		r.
			With(middleware.Authorization(string(zitadel.Owner))).
			Post("/decrease-stock", c.decreaseProductStock)
	})

}
