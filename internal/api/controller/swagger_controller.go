package controller

import (
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/oprimogus/flyfood-api/api"
)

func swaggerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./api/swaggerHandler.json")
}

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
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(htmlContent))
}

func SetupSwaggerRoutes(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Get("/swagger.json", swaggerHandler)
		r.Get("/docs", docsHandler)
	})
}
