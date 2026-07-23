package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"sync"
)

var (
	once            sync.Once
	metricsInstance *PrometheusMetrics
)

type PrometheusMetrics struct {
	Registry          *prometheus.Registry
	RequestCounter    *prometheus.CounterVec
	ResponseTime      *prometheus.HistogramVec
	ErrorCounter      *prometheus.CounterVec
	ActiveConnections prometheus.Gauge
}

func newPrometheusMetrics() *PrometheusMetrics {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector())

	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "metrics_requests_total",
		Help: "Total de requisições recebidas.",
	}, []string{"service", "path", "method", "status"})

	responseTime := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "metrics_response_time_seconds",
		Help:    "Tempo de resposta da API.",
		Buckets: prometheus.DefBuckets,
	}, []string{"service", "path", "method", "status"})

	errorCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "metrics_errors_total",
		Help: "Total de erros da API por endpoint, método e código de status.",
	}, []string{"service", "path", "method", "status"})

	activeConnections := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "metrics_active_connections",
		Help: "Número atual de conexões ativas.",
	})

	reg.MustRegister(requestCounter, responseTime, errorCounter, activeConnections)

	metricsInstance = &PrometheusMetrics{
		Registry:          reg,
		RequestCounter:    requestCounter,
		ResponseTime:      responseTime,
		ErrorCounter:      errorCounter,
		ActiveConnections: activeConnections,
	}

	return metricsInstance
}

func GetPrometheusMetrics() *PrometheusMetrics {
	once.Do(func() {
		metricsInstance = newPrometheusMetrics()
	})
	return metricsInstance
}
