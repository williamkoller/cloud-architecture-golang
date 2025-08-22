package metrics_router

import (
	"github.com/gin-gonic/gin"
	mhandler "github.com/williamkoller/cloud-architecture-golang/internal/metrics/handler"
)

func RegisterMetricsRoute(router *gin.Engine, h *mhandler.MetricsHandler) {
	router.GET("/metrics", h.Serve)
}
