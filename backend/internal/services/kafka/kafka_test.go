package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/user/airpollution/internal/models"
)

// MockProducer implements a mock Kafka producer for testing
type MockProducer struct {
	messages [][]byte
	topic    string
}

// NewMockProducer creates a new mock producer
func NewMockProducer(topic string) *MockProducer {
	return &MockProducer{
		messages: make([][]byte, 0),
		topic:    topic,
	}
}

func (m *MockProducer) ProduceAirQualityData(ctx context.Context, data *models.AirQualityData) error {
	m.messages = append(m.messages, []byte("air-data"))
	return nil
}

func (m *MockProducer) ProduceAnomaly(ctx context.Context, anomaly *models.Anomaly) error {
	m.messages = append(m.messages, []byte("anomaly"))
	return nil
}

func (m *MockProducer) Close() error {
	return nil
}

// MockConsumer implements a mock Kafka consumer for testing
type MockConsumer struct {
	airQualityData []*models.AirQualityData
	anomalies      []*models.Anomaly
	index          int
}

// NewMockConsumer creates a new mock consumer
func NewMockConsumer() *MockConsumer {
	return &MockConsumer{
		airQualityData: make([]*models.AirQualityData, 0),
		anomalies:      make([]*models.Anomaly, 0),
		index:          0,
	}
}

// AddAirQualityData adds air quality data to the mock consumer
func (m *MockConsumer) AddAirQualityData(data *models.AirQualityData) {
	m.airQualityData = append(m.airQualityData, data)
}

// AddAnomaly adds an anomaly to the mock consumer
func (m *MockConsumer) AddAnomaly(anomaly *models.Anomaly) {
	m.anomalies = append(m.anomalies, anomaly)
}

func (m *MockConsumer) ConsumeAirQualityData(ctx context.Context) (*models.AirQualityData, error) {
	if m.index >= len(m.airQualityData) {
		return nil, nil
	}
	data := m.airQualityData[m.index]
	m.index++
	return data, nil
}

func (m *MockConsumer) ConsumeAnomaly(ctx context.Context) (*models.Anomaly, error) {
	if m.index >= len(m.anomalies) {
		return nil, nil
	}
	anomaly := m.anomalies[m.index]
	m.index++
	return anomaly, nil
}

func (m *MockConsumer) Close() error {
	return nil
}

func TestMockProducer(t *testing.T) {
	mockProducer := NewMockProducer(RawAirDataTopic)

	// Test data
	testData := &models.AirQualityData{
		ID:        uuid.New(),
		Parameter: "PM2.5",
		Value:     25.0,
		Latitude:  41.015,
		Longitude: 28.979,
		Timestamp: time.Now(),
	}

	// Test produce
	ctx := context.Background()
	err := mockProducer.ProduceAirQualityData(ctx, testData)
	if err != nil {
		t.Fatalf("Failed to produce data: %v", err)
	}

	// Verify message was produced
	if len(mockProducer.messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(mockProducer.messages))
	}
}

func TestMockConsumer(t *testing.T) {
	mockConsumer := NewMockConsumer()

	// Test data
	testData := &models.AirQualityData{
		ID:        uuid.New(),
		Parameter: "PM2.5",
		Value:     25.0,
		Latitude:  41.015,
		Longitude: 28.979,
		Timestamp: time.Now(),
	}

	// Add test data to consumer
	mockConsumer.AddAirQualityData(testData)

	// Test consume
	ctx := context.Background()
	data, err := mockConsumer.ConsumeAirQualityData(ctx)
	if err != nil {
		t.Fatalf("Failed to consume data: %v", err)
	}

	// Verify data was consumed
	if data == nil {
		t.Errorf("Expected data, got nil")
		return
	}

	if data.ID != testData.ID {
		t.Errorf("Expected ID %s, got %s", testData.ID, data.ID)
	}
}
