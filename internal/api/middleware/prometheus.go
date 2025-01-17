package middleware

import (
	"fmt"
	"github.com/oprimogus/cardapiogo/internal/config"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/observability"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

func Prometheus(next http.Handler) http.Handler {
	metrics := observability.GetPrometheusMetrics()

	service := config.GetInstance().Api.ServiceName

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
