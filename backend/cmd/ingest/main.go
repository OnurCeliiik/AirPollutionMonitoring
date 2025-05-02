package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/user/airpollution/internal/api"
	"github.com/user/airpollution/internal/services/kafka"
	"github.com/user/airpollution/internal/services/logger"
)

// @title Air Quality Monitoring API
// @version 1.0
// @description API for air quality monitoring platform
// @host localhost:8080
// @BasePath /
func main() {
	// Configure logger based on environment
	logLevel := getEnv("LOG_LEVEL", "INFO")
	logger.SetDefaultLogLevel(logLevel)

	// Use production mode in non-local environments
	env := getEnv("ENVIRONMENT", "development")
	if env != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Load configuration from environment variables
	kafkaBrokers := getEnv("KAFKA_BROKERS", "localhost:9092")
	port := getEnv("PORT", "8080")

	// CORS allowed origins
	allowedOrigins := getEnv("ALLOWED_ORIGINS", "*")

	logger.Info("Starting Ingest Service with Kafka brokers: %s", kafkaBrokers)

	// Create Kafka producer
	producer := kafka.NewProducer([]string{kafkaBrokers}, kafka.RawAirDataTopic)
	defer producer.Close()

	// Create Gin router
	router := gin.Default()

	// Add CORS middleware with improved security
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "up",
			"service": "ingest",
		})
	})

	// Setup API routes
	ingestHandler := api.NewIngestHandler(producer)
	ingestHandler.RegisterRoutes(router)

	// Setup Swagger
	if env == "development" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		logger.Info("Swagger UI available at http://localhost:%s/swagger/index.html", port)
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		logger.Info("Ingest Service starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited gracefully")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
