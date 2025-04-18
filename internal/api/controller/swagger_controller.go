package controller

import (
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/oprimogus/flyfood-api/internal/config"
	"net/http"
)

func swaggerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./api/swaggerHandler.json")
}

func docsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: ".//swagger.json",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "FlyFood",
		},
		DarkMode: true,
	})
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(htmlContent))
}

func SetupSwaggerRoutes(r chi.Router) {
	api := config.GetInstance().Api

	r.Route(api.BasePath, func(r chi.Router) {
		r.Get("/swagger.json", swaggerHandler)
		r.Get("/docs", docsHandler)

	})
}
