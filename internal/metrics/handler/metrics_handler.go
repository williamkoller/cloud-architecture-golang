package metrics_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/williamkoller/cloud-architecture-golang/internal/metrics"
)

type MetricsHandler struct {
	httpHandler http.Handler
}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{
		httpHandler: promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{}),
	}
}

func (h *MetricsHandler) Serve(c *gin.Context) {
	h.httpHandler.ServeHTTP(c.Writer, c.Request)
}
