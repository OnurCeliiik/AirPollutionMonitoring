package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/user/airpollution/internal/models"
)

const (
	RawAirDataTopic    = "raw-air-data"
	AnomalyAlertsTopic = "anomaly-alerts"
)

// Producer handles producing messages to Kafka
type Producer struct {
	writer *kafka.Writer
}

// Consumer handles consuming messages from Kafka
type Consumer struct {
	reader *kafka.Reader
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{
		writer: writer,
	}
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.FirstOffset,
		MaxWait:     time.Second,
	})

	return &Consumer{
		reader: reader,
	}
}

// ProduceAirQualityData produces an air quality data message with retries
func (p *Producer) ProduceAirQualityData(ctx context.Context, data *models.AirQualityData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling air quality data: %w", err)
	}

	// Retry logic
	maxRetries := 3
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		err = p.writer.WriteMessages(ctx, kafka.Message{
			Value: jsonData,
		})
		if err == nil {
			return nil
		}
		lastErr = err
		// Exponential backoff: 100ms, 200ms, 400ms
		backoff := time.Duration(100*(1<<i)) * time.Millisecond
		time.Sleep(backoff)
	}

	return fmt.Errorf("error writing message to Kafka after %d retries: %w", maxRetries, lastErr)
}

// ProduceAnomaly produces an anomaly message with retries
func (p *Producer) ProduceAnomaly(ctx context.Context, anomaly *models.Anomaly) error {
	jsonData, err := json.Marshal(anomaly)
	if err != nil {
		return fmt.Errorf("error marshaling anomaly: %w", err)
	}

	// Retry logic
	maxRetries := 3
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		err = p.writer.WriteMessages(ctx, kafka.Message{
			Value: jsonData,
		})
		if err == nil {
			return nil
		}
		lastErr = err
		// Exponential backoff: 100ms, 200ms, 400ms
		backoff := time.Duration(100*(1<<i)) * time.Millisecond
		time.Sleep(backoff)
	}

	return fmt.Errorf("error writing message to Kafka after %d retries: %w", maxRetries, lastErr)
}

// ConsumeAirQualityData consumes air quality data messages
func (c *Consumer) ConsumeAirQualityData(ctx context.Context) (*models.AirQualityData, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading message from Kafka: %w", err)
	}

	var data models.AirQualityData
	if err := json.Unmarshal(msg.Value, &data); err != nil {
		return nil, fmt.Errorf("error unmarshaling air quality data: %w", err)
	}

	return &data, nil
}

// ConsumeAnomaly consumes anomaly messages
func (c *Consumer) ConsumeAnomaly(ctx context.Context) (*models.Anomaly, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading message from Kafka: %w", err)
	}

	var anomaly models.Anomaly
	if err := json.Unmarshal(msg.Value, &anomaly); err != nil {
		return nil, fmt.Errorf("error unmarshaling anomaly: %w", err)
	}

	return &anomaly, nil
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.writer.Close()
}

// Close closes the consumer
func (c *Consumer) Close() error {
	return c.reader.Close()
}
