package middleware

import (
	"fmt"
	"github.com/oprimogus/flyfood-api/internal/config"
	"github.com/oprimogus/flyfood-api/internal/infra/observability"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

func Prometheus(next http.Handler) http.Handler {
	metrics := observability.GetPrometheusMetrics()

	service := config.Get().API.Name

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := NewResponseRecorder(w)

		path := GetPath(r.URL.Path)
		method := r.Method

		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			metrics.ResponseTime.WithLabelValues(service, path, method, fmt.Sprintf("%d", rw.Status)).Observe(v)
		}))
		defer timer.ObserveDuration()

		metrics.ActiveConnections.Inc()
		defer metrics.ActiveConnections.Dec()

		next.ServeHTTP(rw, r)

		metrics.RequestCounter.WithLabelValues(service, path, method, fmt.Sprintf("%d", rw.Status)).Inc()
	})
}
