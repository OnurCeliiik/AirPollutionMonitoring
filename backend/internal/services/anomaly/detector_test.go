package anomaly

import (
	"testing"
	"time"

	"github.com/user/airpollution/internal/models"
)

func TestCheckThresholdExceeded(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		name      string
		parameter string
		value     float64
		expected  bool
	}{
		{"PM2.5 Below Threshold", "PM2.5", 10.0, false},
		{"PM2.5 Above Threshold", "PM2.5", 20.0, true},
		{"PM10 Below Threshold", "PM10", 40.0, false},
		{"PM10 Above Threshold", "PM10", 50.0, true},
		{"NO2 Below Threshold", "NO2", 20.0, false},
		{"NO2 Above Threshold", "NO2", 30.0, true},
		{"O3 Below Threshold", "O3", 90.0, false},
		{"O3 Above Threshold", "O3", 110.0, true},
		{"Unknown Parameter", "CO", 100.0, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data := &models.AirQualityData{
				Parameter: tc.parameter,
				Value:     tc.value,
				Latitude:  41.015,
				Longitude: 28.979,
				Timestamp: time.Now(),
			}

			anomaly := detector.checkThresholdExceeded(data)

			if tc.expected && anomaly == nil {
				t.Errorf("Expected anomaly but got nil")
			}

			if !tc.expected && anomaly != nil {
				t.Errorf("Expected no anomaly but got %v", anomaly)
			}

			if anomaly != nil && anomaly.Type != string(models.ThresholdExceeded) {
				t.Errorf("Expected anomaly type %s but got %s", models.ThresholdExceeded, anomaly.Type)
			}
		})
	}
}

func TestCheckStatisticalOutlier(t *testing.T) {
	detector := NewDetector()

	// Create historical data
	historicalData := make([]models.AirQualityData, 0)
	now := time.Now()

	// Add 10 data points around value 20.0
	for i := 0; i < 10; i++ {
		historicalData = append(historicalData, models.AirQualityData{
			Parameter: "PM2.5",
			Value:     20.0 + float64(i%3), // Values will be 20.0, 21.0, 22.0
			Latitude:  41.015,
			Longitude: 28.979,
			Timestamp: now.Add(-time.Duration(i) * time.Hour),
		})
	}

	// Test normal value (not an outlier)
	normal := &models.AirQualityData{
		Parameter: "PM2.5",
		Value:     23.0, // Close to historical values
		Latitude:  41.015,
		Longitude: 28.979,
		Timestamp: now,
	}

	// Test outlier value
	outlier := &models.AirQualityData{
		Parameter: "PM2.5",
		Value:     50.0, // Far from historical values
		Latitude:  41.015,
		Longitude: 28.979,
		Timestamp: now,
	}

	// Test normal value
	anomaly := detector.checkStatisticalOutlier(normal, historicalData)
	if anomaly != nil {
		t.Errorf("Expected no anomaly for normal value but got %v", anomaly)
	}

	// Test outlier
	anomaly = detector.checkStatisticalOutlier(outlier, historicalData)
	if anomaly == nil {
		t.Errorf("Expected anomaly for outlier value but got nil")
	}

	if anomaly != nil && anomaly.Type != string(models.StatisticalOutlier) {
		t.Errorf("Expected anomaly type %s but got %s", models.StatisticalOutlier, anomaly.Type)
	}
}
