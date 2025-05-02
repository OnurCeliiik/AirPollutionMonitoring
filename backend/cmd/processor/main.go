package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/airpollution/internal/db"
	"github.com/user/airpollution/internal/services/anomaly"
	"github.com/user/airpollution/internal/services/kafka"
)

func main() {
	// Load configuration from environment variables
	kafkaBrokers := getEnv("KAFKA_BROKERS", "localhost:9092")
	dbConnStr := getEnv("DB_CONNECTION_STRING", "postgres://postgres:postgres@localhost:5432/timescaledb?sslmode=disable")

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to database
	database, err := db.New(dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create Kafka consumer for raw air data
	consumer := kafka.NewConsumer(
		[]string{kafkaBrokers},
		kafka.RawAirDataTopic,
		"processor-group",
	)
	defer consumer.Close()

	// Create Kafka producer for anomaly alerts
	producer := kafka.NewProducer(
		[]string{kafkaBrokers},
		kafka.AnomalyAlertsTopic,
	)
	defer producer.Close()

	// Create anomaly detector
	detector := anomaly.NewDetector()

	// Process messages in a goroutine
	go processMessages(ctx, consumer, producer, database, detector)

	// Wait for interrupt signal to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down processor service...")
	cancel()                    // Signal the processMessages goroutine to stop
	time.Sleep(1 * time.Second) // Give it a moment to clean up
	log.Println("Processor service exited")
}

// processMessages continuously processes messages from Kafka
func processMessages(ctx context.Context, consumer *kafka.Consumer, producer *kafka.Producer, database *db.DB, detector *anomaly.Detector) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Set a timeout for the consume operation
			msgCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			data, err := consumer.ConsumeAirQualityData(msgCtx)
			cancel()

			if err != nil {
				log.Printf("Error consuming message: %v", err)
				time.Sleep(1 * time.Second) // Wait a bit before retrying
				continue
			}

			log.Printf("Processing air quality data: %s at [%f,%f]: %f",
				data.Parameter, data.Latitude, data.Longitude, data.Value)

			// Insert into database
			if err := database.InsertAirQualityData(data); err != nil {
				log.Printf("Error inserting data into database: %v", err)
				continue
			}

			// Get recent data for anomaly detection
			recentData, err := database.GetRecentDataForParameter(
				data.Parameter,
				data.Latitude,
				data.Longitude,
				24, // Last 24 hours
			)
			if err != nil {
				log.Printf("Error fetching recent data: %v", err)
				continue
			}

			// Detect anomalies
			anomalyResult, err := detector.Detect(data, recentData)
			if err != nil {
				log.Printf("Error detecting anomalies: %v", err)
				continue
			}

			// If anomaly detected, publish to anomaly alerts topic
			if anomalyResult != nil {
				log.Printf("Anomaly detected: %s - %s - %f",
					anomalyResult.Type, anomalyResult.Parameter, anomalyResult.Value)

				// Insert anomaly into database
				if err := database.InsertAnomaly(anomalyResult); err != nil {
					log.Printf("Error inserting anomaly into database: %v", err)
				}

				// Publish to Kafka
				alertCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
				err := producer.ProduceAnomaly(alertCtx, anomalyResult)
				cancel()

				if err != nil {
					log.Printf("Error publishing anomaly alert: %v", err)
				}
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
