package models

import (
	"time"

	"github.com/google/uuid"
)

// AirQualityData represents a data point for air quality
type AirQualityData struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
	Parameter string    `json:"parameter" db:"parameter"`
	Value     float64   `json:"value" db:"value"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

// Anomaly represents an anomaly in air quality data
type Anomaly struct {
	ID                      uuid.UUID `json:"id" db:"id"`
	Type                    string    `json:"type" db:"type"`
	Parameter               string    `json:"parameter" db:"parameter"`
	Value                   float64   `json:"value" db:"value"`
	Latitude                float64   `json:"latitude" db:"latitude"`
	Longitude               float64   `json:"longitude" db:"longitude"`
	DetectedAt              time.Time `json:"detected_at" db:"detected_at"`
	AirQualityDataID        uuid.UUID `json:"air_quality_data_id,omitempty" db:"air_quality_data_id"`
	AirQualityDataTimestamp time.Time `json:"air_quality_data_timestamp,omitempty" db:"air_quality_data_timestamp"`
}

// AnomalyType represents the type of anomaly detected
type AnomalyType string

const (
	ThresholdExceeded       AnomalyType = "ThresholdExceeded"
	StatisticalOutlier      AnomalyType = "StatisticalOutlier"
	SpikeDetected           AnomalyType = "SpikeDetected"
	GeographicInconsistency AnomalyType = "GeographicInconsistency"
)

// AnomalyAlert represents the message sent via WebSocket to clients
type AnomalyAlert struct {
	Parameter string    `json:"parameter"`
	Value     float64   `json:"value"`
	Type      string    `json:"type"`
	Location  []float64 `json:"location"` // [latitude, longitude]
	Timestamp time.Time `json:"timestamp"`
}

// NewAirQualityData creates a new air quality data point
func NewAirQualityData(latitude, longitude float64, parameter string, value float64, timestamp time.Time) *AirQualityData {
	return &AirQualityData{
		ID:        uuid.New(),
		Latitude:  latitude,
		Longitude: longitude,
		Parameter: parameter,
		Value:     value,
		Timestamp: timestamp,
	}
}

// NewAnomaly creates a new anomaly
func NewAnomaly(anomalyType string, parameter string, value, latitude, longitude float64) *Anomaly {
	return &Anomaly{
		ID:         uuid.New(),
		Type:       anomalyType,
		Parameter:  parameter,
		Value:      value,
		Latitude:   latitude,
		Longitude:  longitude,
		DetectedAt: time.Now(),
		// AirQualityDataID and AirQualityDataTimestamp are optional and can be set later
	}
}

// NewAnomalyFromData creates a new anomaly from air quality data
func NewAnomalyFromData(anomalyType string, data *AirQualityData) *Anomaly {
	return &Anomaly{
		ID:                      uuid.New(),
		Type:                    anomalyType,
		Parameter:               data.Parameter,
		Value:                   data.Value,
		Latitude:                data.Latitude,
		Longitude:               data.Longitude,
		DetectedAt:              time.Now(),
		AirQualityDataID:        data.ID,
		AirQualityDataTimestamp: data.Timestamp,
	}
}

// ToAnomalyAlert converts an Anomaly to an AnomalyAlert
func (a *Anomaly) ToAnomalyAlert() *AnomalyAlert {
	return &AnomalyAlert{
		Parameter: a.Parameter,
		Value:     a.Value,
		Type:      a.Type,
		Location:  []float64{a.Latitude, a.Longitude},
		Timestamp: a.DetectedAt,
	}
}
