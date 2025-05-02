package db

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/user/airpollution/internal/models"
)

// TestDB is a simple test to verify DB operations
// Note: These tests require a real TimescaleDB instance to run
// In a production environment, you would use a test database or mock the DB
func TestDB_Operations(t *testing.T) {
	t.Skip("Skipping DB tests that require a real database. Remove this line to run tests with a real DB.")

	// Connect to test database
	connStr := "postgres://postgres:postgres@localhost:5432/timescaledb?sslmode=disable"
	db, err := New(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test data
	testData := &models.AirQualityData{
		ID:        uuid.New(),
		Parameter: "PM2.5",
		Value:     25.0,
		Latitude:  41.015,
		Longitude: 28.979,
		Timestamp: time.Now(),
	}

	// Test insert
	err = db.InsertAirQualityData(testData)
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// Test query
	results, err := db.GetRecentDataForParameter(testData.Parameter, testData.Latitude, testData.Longitude, 24)
	if err != nil {
		t.Fatalf("Failed to query data: %v", err)
	}

	// Verify results
	found := false
	for _, data := range results {
		if data.ID == testData.ID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Inserted data not found in query results")
	}
}

// MockDB implements a mock DB for testing
type MockDB struct {
	data      []models.AirQualityData
	anomalies []models.Anomaly
}

// NewMockDB creates a new mock DB
func NewMockDB() *MockDB {
	return &MockDB{
		data:      make([]models.AirQualityData, 0),
		anomalies: make([]models.Anomaly, 0),
	}
}

func (m *MockDB) InsertAirQualityData(data *models.AirQualityData) error {
	m.data = append(m.data, *data)
	return nil
}

func (m *MockDB) InsertAnomaly(anomaly *models.Anomaly) error {
	m.anomalies = append(m.anomalies, *anomaly)
	return nil
}

func (m *MockDB) GetRecentDataForParameter(parameter string, latitude, longitude float64, hours int) ([]models.AirQualityData, error) {
	results := make([]models.AirQualityData, 0)
	for _, data := range m.data {
		if data.Parameter == parameter {
			results = append(results, data)
		}
	}
	return results, nil
}

func (m *MockDB) GetRecentAnomalies(hours int) ([]models.Anomaly, error) {
	return m.anomalies, nil
}

func TestMockDB_Operations(t *testing.T) {
	mockDB := NewMockDB()

	// Test data
	testData := &models.AirQualityData{
		ID:        uuid.New(),
		Parameter: "PM2.5",
		Value:     25.0,
		Latitude:  41.015,
		Longitude: 28.979,
		Timestamp: time.Now(),
	}

	// Test insert
	err := mockDB.InsertAirQualityData(testData)
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// Test query
	results, err := mockDB.GetRecentDataForParameter(testData.Parameter, testData.Latitude, testData.Longitude, 24)
	if err != nil {
		t.Fatalf("Failed to query data: %v", err)
	}

	// Verify results
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
}
