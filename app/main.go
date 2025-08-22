package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/williamkoller/cloud-architecture-golang/internal/metrics"
	metrics_handler "github.com/williamkoller/cloud-architecture-golang/internal/metrics/handler"
	metrics_router "github.com/williamkoller/cloud-architecture-golang/internal/metrics/router"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/handler"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/repository"
	usr_router "github.com/williamkoller/cloud-architecture-golang/internal/usr/router"
)

var (
	router      *gin.Engine
	ginLambdaV2 *ginadapter.GinLambdaV2
)

func healthMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.String(http.StatusOK, "OK")
			c.Abort()
			return
		}
		c.Next()
	})
}

func init() {
	metrics.Init("cloud-arch-golang", "1.0.0")

	gin.SetMode(gin.ReleaseMode)
	router = gin.New()

	router.Use(healthMiddleware())
	router.Use(metrics.RecoveryMiddleware())
	// Temporariamente desabilitado: router.Use(rateLimitMiddleware())
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))
	router.Use(metrics.Middleware())

	api := router.Group("/api")

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	router.GET("/test/error500", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno simulado"})
	})

	router.GET("/test/panic", func(c *gin.Context) {
		panic("Panic simulado para teste de métricas")
	})

	mh := metrics_handler.NewMetricsHandler()
	metrics_router.RegisterMetricsRoute(router, mh)

	userRepo := repository.NewInMemoryUserRepository()
	userHandler := handler.NewUserHandler(userRepo)
	usr_router.RegisterUserRoutes(api, userHandler)

	ginLambdaV2 = ginadapter.NewV2(router)
}

func main() {
	if os.Getenv("LOCAL") == "true" {
		log.Println("Starting server locally on :8080")

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Servidor HTTP otimizado para alta performance
		srv := &http.Server{
			Addr:    ":8080",
			Handler: router,
			// Timeouts otimizados para alta carga
			ReadTimeout:       10 * time.Second, // Tempo para ler request
			WriteTimeout:      30 * time.Second, // Tempo para escrever response
			IdleTimeout:       60 * time.Second, // Tempo para keep-alive
			ReadHeaderTimeout: 5 * time.Second,  // Tempo para ler headers
			MaxHeaderBytes:    1 << 20,          // 1MB para headers
		}

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Failed to start server: %v", err)
			}
		}()

		log.Println("Server started successfully")
		log.Println("- Health check: http://localhost:8080/health")
		log.Println("- Metrics: http://localhost:8080/metrics")
		log.Println("- API: http://localhost:8080/api/v1")
		log.Println("Press Ctrl+C to shutdown...")

		<-sigChan
		log.Println("Shutting down server gracefully...")

		// Graceful shutdown com timeout generoso
		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}

		log.Println("Server stopped successfully")
		return
	}

	// Para execução em Lambda
	lambda.Start(ginLambdaV2.ProxyWithContext)
}
