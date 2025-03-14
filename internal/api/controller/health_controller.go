package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/oprimogus/flyfood-api/internal/config"
	postgresDB "github.com/oprimogus/flyfood-api/internal/infrastructure/database/postgres"
	"github.com/oprimogus/flyfood-api/internal/infrastructure/observability"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

type healthController struct {
	db *postgresDB.Database
}

type HealthResponse struct {
	Status    string `json:"status" example:"UP"`
	Timestamp string `json:"timestamp" example:"2024-03-15T10:30:15Z"`
}

func newHealthHandler(db *postgresDB.Database) healthController {
	return healthController{db: db}
}

// livenessProbe godoc
//
//	@Summary		Liveness Prob
//	@Description	Liveness Prob
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	HealthResponse
//	@Failure		500	{object}	HealthResponse
//	@Router			/health/liveness [get]
func (h *healthController) livenessProbe(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "UP",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	JSONResponse(w, http.StatusOK, response)
}

// readinessProbe godoc
//
//	@Summary		Readiness Prob
//	@Description	Readiness Prob
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	HealthResponse
//	@Failure		500	{object}	HealthResponse
//	@Router			/health/readiness [get]
func (h *healthController) readinessProbe(w http.ResponseWriter, r *http.Request) {
	if err := h.db.GetDB().Ping(r.Context()); err != nil {
		response := HealthResponse{
			Status:    "DOWN",
			Timestamp: time.Now().Format(time.RFC3339),
		}
		JSONResponse(w, http.StatusServiceUnavailable, response)
		return
	}

	response := HealthResponse{
		Status:    "UP",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	JSONResponse(w, http.StatusOK, response)
}

// prometheusMetrics godoc
//
//	@Summary		Prometheus Metrics
//	@Description	Prometheus Metrics
//	@Tags			Health
//	@Produce		plain
//	@Success		200
//	@Failure		500	{object}	HealthResponse
//	@Router			/health/metrics [get]
func (h *healthController) prometheusMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := observability.GetPrometheusMetrics()
	prometheusHandler := promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{})
	prometheusHandler.ServeHTTP(w, r)
}

func SetupHealthRoutes(r chi.Router, db *postgresDB.Database) {
	basePath := config.GetInstance().Api.BasePath
	handler := newHealthHandler(db)

	r.Route(basePath+"/health", func(r chi.Router) {
		r.Get("/liveness", handler.livenessProbe)
		r.Get("/readiness", handler.readinessProbe)
		r.Get("/metrics", handler.prometheusMetrics)
	})
}
