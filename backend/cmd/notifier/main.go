package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/airpollution/internal/db"
	"github.com/user/airpollution/internal/services/kafka"
	"github.com/user/airpollution/internal/services/websocket"
)

func main() {
	// Load configuration from environment variables
	kafkaBrokers := getEnv("KAFKA_BROKERS", "localhost:9092")
	dbConnStr := getEnv("DB_CONNECTION_STRING", "postgres://postgres:postgres@localhost:5432/timescaledb?sslmode=disable")
	port := getEnv("PORT", "8081")

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to database
	database, err := db.New(dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create Kafka consumer for anomaly alerts
	consumer := kafka.NewConsumer(
		[]string{kafkaBrokers},
		kafka.AnomalyAlertsTopic,
		"notifier-group",
	)
	defer consumer.Close()

	// Create WebSocket hub
	hub := websocket.NewHub()
	go hub.Run(ctx)

	// Create Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

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
			"service": "notifier",
		})
	})

	// WebSocket endpoint
	router.GET("/ws/alerts", func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	})

	// Get recent anomalies endpoint
	router.GET("/api/anomalies", func(c *gin.Context) {
		// Default to last 24 hours
		hours := 24
		if hoursParam := c.Query("hours"); hoursParam != "" {
			if _, err := time.ParseDuration(hoursParam + "h"); err == nil {
				hours = 24 // Keep default if invalid
			}
		}

		anomalies, err := database.GetRecentAnomalies(hours)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch anomalies: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, anomalies)
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		log.Printf("Notifier Service starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Process messages in a goroutine
	go processAnomalyAlerts(ctx, consumer, hub)

	// Wait for interrupt signal to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down notifier service...")
	cancel() // Signal goroutines to stop

	// Shutdown HTTP server
	httpCtx, httpCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer httpCancel()
	if err := srv.Shutdown(httpCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	time.Sleep(1 * time.Second) // Give other goroutines a moment to clean up
	log.Println("Notifier service exited")
}

// processAnomalyAlerts continuously processes anomaly alerts from Kafka
func processAnomalyAlerts(ctx context.Context, consumer *kafka.Consumer, hub *websocket.Hub) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Set a timeout for the consume operation
			msgCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			anomaly, err := consumer.ConsumeAnomaly(msgCtx)
			cancel()

			if err != nil {
				log.Printf("Error consuming anomaly message: %v", err)
				time.Sleep(1 * time.Second) // Wait a bit before retrying
				continue
			}

			log.Printf("Received anomaly alert: %s - %s - %f",
				anomaly.Type, anomaly.Parameter, anomaly.Value)

			// Convert to alert format for WebSocket
			alert := anomaly.ToAnomalyAlert()

			// Broadcast to all WebSocket clients
			if err := hub.BroadcastAnomaly(alert); err != nil {
				log.Printf("Error broadcasting anomaly: %v", err)
			}
		}
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
