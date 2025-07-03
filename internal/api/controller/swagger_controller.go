package controller

import (
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/oprimogus/flyfood-api/api"
	"github.com/oprimogus/flyfood-api/internal/core"
	_ "github.com/oprimogus/flyfood-api/api"
	"github.com/oprimogus/flyfood-api/internal/config"
	httpSwagger "github.com/swaggo/http-swagger"
)

func docsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecContent: api.SwaggerInfo.ReadDoc(),
		CustomOptions: scalar.CustomOptions{
			PageTitle: "FlyFood",
		},
		DarkMode: true,
	})
	if err != nil {
		core.HandleApiError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(htmlContent))
}

func SetupSwaggerRoutes(r chi.Router) {
	basePath := config.GetInstance().Api.BasePath
	if basePath == "" {
		basePath = "/"
	}
	r.Route(basePath, func(r chi.Router) {
		r.Get("/docs", docsHandler)
		r.Get("/swagger/*", httpSwagger.WrapHandler)
	})
}
