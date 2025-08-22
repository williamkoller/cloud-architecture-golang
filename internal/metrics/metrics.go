package metrics

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	Registry = prometheus.NewRegistry()

	serviceLabel string
	versionLabel string

	httpInFlight           prometheus.Gauge
	httpRequestsTotal      *prometheus.CounterVec
	httpRequestsByClass    *prometheus.CounterVec
	httpReqDuration        *prometheus.HistogramVec
	httpReqDurationByClass *prometheus.HistogramVec
	httpRespSizeBytes      *prometheus.HistogramVec

	// Métricas de domínio simplificadas
	usersCreatedTotal *prometheus.CounterVec
	usersUpdatedTotal *prometheus.CounterVec
	usersDeletedTotal *prometheus.CounterVec

	appInfo             prometheus.Gauge
	panicRecoveredTotal *prometheus.CounterVec
)

func Init(service, version string) {
	serviceLabel = service
	versionLabel = version

	// Labels básicos para todas as métricas
	constLabels := prometheus.Labels{"service": serviceLabel, "version": versionLabel}

	httpInFlight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "http_in_flight_requests",
		Help:        "Current number of in-flight HTTP requests.",
		ConstLabels: constLabels,
	})

	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests by method/route/status.",
		},
		[]string{"method", "route", "status", "service", "version"},
	)

	httpRequestsByClass = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_class_total",
			Help: "HTTP requests by status class (2xx/4xx/5xx).",
		},
		[]string{"method", "route", "class", "service", "version"},
	)

	// Buckets otimizados para APIs REST
	httpReqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "route", "status", "service", "version"},
	)

	httpReqDurationByClass = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_class_duration_seconds",
			Help:    "HTTP request duration by status class.",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "route", "class", "service", "version"},
	)

	// Buckets otimizados para tamanhos de resposta típicos
	httpRespSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes.",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 500000},
		},
		[]string{"method", "route", "status", "service", "version"},
	)

	// Métricas de domínio simplificadas
	usersCreatedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "users_created_total",
			Help: "Total users created.",
		},
		[]string{"service", "version"},
	)

	usersUpdatedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "users_updated_total",
			Help: "Total users updated.",
		},
		[]string{"service", "version"},
	)

	usersDeletedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "users_deleted_total",
			Help: "Total users deleted.",
		},
		[]string{"service", "version"},
	)

	appInfo = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "app_info",
		Help:        "Application info (constant 1 with service/version labels).",
		ConstLabels: constLabels,
	})

	panicRecoveredTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "panic_recovered_total",
			Help: "Panics recovered by middleware.",
		},
		[]string{"service", "version"},
	)

	// Registrar métricas com coletor Go otimizado
	Registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),

		httpInFlight,
		httpRequestsTotal,
		httpRequestsByClass,
		httpReqDuration,
		httpReqDurationByClass,
		httpRespSizeBytes,

		usersCreatedTotal,
		usersUpdatedTotal,
		usersDeletedTotal,

		appInfo,
		panicRecoveredTotal,
	)

	appInfo.Set(1)
}

// Middleware otimizado para alta performance
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pular métricas para health checks para reduzir overhead
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		start := time.Now()
		method := c.Request.Method
		route := c.FullPath()
		if route == "" {
			route = "UNMATCHED"
		}

		httpInFlight.Inc()
		defer httpInFlight.Dec()

		c.Next()

		// Calcular métricas após processamento
		statusCode := c.Writer.Status()
		elapsed := time.Since(start).Seconds()

		// Labels pré-computados para melhor performance
		status := strconv.Itoa(statusCode)
		class := fmt.Sprintf("%dxx", statusCode/100)

		// Atualizar métricas de forma eficiente
		httpRequestsTotal.WithLabelValues(method, route, status, serviceLabel, versionLabel).Inc()
		httpRequestsByClass.WithLabelValues(method, route, class, serviceLabel, versionLabel).Inc()
		httpReqDuration.WithLabelValues(method, route, status, serviceLabel, versionLabel).Observe(elapsed)
		httpReqDurationByClass.WithLabelValues(method, route, class, serviceLabel, versionLabel).Observe(elapsed)

		// Métricas de tamanho apenas se significativas
		if sz := c.Writer.Size(); sz > 0 {
			httpRespSizeBytes.WithLabelValues(method, route, status, serviceLabel, versionLabel).Observe(float64(sz))
		}
	}
}

// RecoveryMiddleware otimizado para capturar panics
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		panicRecoveredTotal.WithLabelValues(serviceLabel, versionLabel).Inc()
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// Helpers de domínio otimizados (chamadas diretas sem overhead)
func UsersCreatedInc() {
	usersCreatedTotal.WithLabelValues(serviceLabel, versionLabel).Inc()
}

func UsersUpdatedInc() {
	usersUpdatedTotal.WithLabelValues(serviceLabel, versionLabel).Inc()
}

func UsersDeletedInc() {
	usersDeletedTotal.WithLabelValues(serviceLabel, versionLabel).Inc()
}

func PanicRecoveredInc() {
	panicRecoveredTotal.WithLabelValues(serviceLabel, versionLabel).Inc()
}
