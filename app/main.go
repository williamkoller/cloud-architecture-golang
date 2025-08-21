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
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/handler"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/repository"
	usr_router "github.com/williamkoller/cloud-architecture-golang/internal/usr/router"
)

var (
	router      *gin.Engine
	ginLambdaV2 *ginadapter.GinLambdaV2
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	api := router.Group("/api")

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	userRepo := repository.NewInMemoryUserRepository()
	useHandler := handler.NewUserHandler(userRepo)
	usr_router.RegisterUserRoutes(api, useHandler)

	ginLambdaV2 = ginadapter.NewV2(router)
}

func main() {
	if os.Getenv("LOCAL") == "true" {
		log.Println("Starting server locally on :8080")

		// Configurar graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Criar servidor HTTP
		srv := &http.Server{
			Addr:    ":8080",
			Handler: router,
		}

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Failed to start server: %v", err)
			}
		}()

		log.Println("Server started. Press Ctrl+C to shutdown...")

		<-sigChan
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}

		log.Println("Server stopped")
		return
	}

	lambda.Start(ginLambdaV2.ProxyWithContext)
}
