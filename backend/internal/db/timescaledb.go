package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/airpollution/internal/models"
)

// DB represents the database connection
type DB struct {
	pool *pgxpool.Pool
}

// New creates a new database connection
func New(connStr string) (*DB, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	db := &DB{pool: pool}
	if err := db.InitSchema(); err != nil {
		return nil, err
	}

	return db, nil
}

// InitSchema initializes the database schema
func (db *DB) InitSchema() error {
	ctx := context.Background()

	// Create air_quality_data table
	_, err := db.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS air_quality_data (
			id UUID,
			latitude FLOAT NOT NULL,
			longitude FLOAT NOT NULL,
			parameter TEXT NOT NULL,
			value FLOAT NOT NULL,
			timestamp TIMESTAMPTZ NOT NULL,
			PRIMARY KEY (id, timestamp)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create air_quality_data table: %w", err)
	}

	// Convert to TimescaleDB hypertable
	_, err = db.pool.Exec(ctx, `
		SELECT create_hypertable('air_quality_data', 'timestamp', if_not_exists => TRUE);
	`)
	if err != nil {
		return fmt.Errorf("failed to create hypertable: %w", err)
	}

	// Create anomalies table
	_, err = db.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS anomalies (
			id UUID,
			type TEXT NOT NULL,
			parameter TEXT NOT NULL,
			value FLOAT NOT NULL,
			latitude FLOAT NOT NULL,
			longitude FLOAT NOT NULL,
			detected_at TIMESTAMPTZ NOT NULL,
			air_quality_data_id UUID,
			air_quality_data_timestamp TIMESTAMPTZ,
			FOREIGN KEY (air_quality_data_id, air_quality_data_timestamp) REFERENCES air_quality_data(id, timestamp),
			PRIMARY KEY (id, detected_at)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create anomalies table: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() {
	db.pool.Close()
}

// InsertAirQualityData inserts a new air quality data point
func (db *DB) InsertAirQualityData(data *models.AirQualityData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.pool.Exec(ctx, `
		INSERT INTO air_quality_data (id, latitude, longitude, parameter, value, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, data.ID, data.Latitude, data.Longitude, data.Parameter, data.Value, data.Timestamp)

	if err != nil {
		return fmt.Errorf("failed to insert air quality data: %w", err)
	}

	return nil
}

// InsertAnomaly inserts a new anomaly
func (db *DB) InsertAnomaly(anomaly *models.Anomaly) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.pool.Exec(ctx, `
		INSERT INTO anomalies (id, type, parameter, value, latitude, longitude, detected_at, air_quality_data_id, air_quality_data_timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, anomaly.ID, anomaly.Type, anomaly.Parameter, anomaly.Value, anomaly.Latitude, anomaly.Longitude, anomaly.DetectedAt, anomaly.AirQualityDataID, anomaly.AirQualityDataTimestamp)

	if err != nil {
		return fmt.Errorf("failed to insert anomaly: %w", err)
	}

	return nil
}

// GetRecentDataForParameter gets the recent data for a specific parameter in a location
func (db *DB) GetRecentDataForParameter(parameter string, latitude, longitude float64, hours int) ([]models.AirQualityData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Calculate a rough distance filter (not accurate but fast)
	latDelta := 0.25 // Approximately 25km
	lonDelta := 0.25

	rows, err := db.pool.Query(ctx, `
		SELECT id, latitude, longitude, parameter, value, timestamp
		FROM air_quality_data
		WHERE parameter = $1
		AND latitude BETWEEN $2 - $3 AND $2 + $3
		AND longitude BETWEEN $4 - $5 AND $4 + $5
		AND timestamp > NOW() - INTERVAL '$6 hours'
		ORDER BY timestamp DESC
	`, parameter, latitude, latDelta, longitude, lonDelta, hours)

	if err != nil {
		return nil, fmt.Errorf("failed to query recent data: %w", err)
	}
	defer rows.Close()

	var results []models.AirQualityData
	for rows.Next() {
		var data models.AirQualityData
		if err := rows.Scan(&data.ID, &data.Latitude, &data.Longitude, &data.Parameter, &data.Value, &data.Timestamp); err != nil {
			return nil, err
		}
		results = append(results, data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// GetRecentAnomalies gets recent anomalies
func (db *DB) GetRecentAnomalies(hours int) ([]models.Anomaly, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.pool.Query(ctx, `
		SELECT id, type, parameter, value, latitude, longitude, detected_at, air_quality_data_id, air_quality_data_timestamp
		FROM anomalies
		WHERE detected_at > NOW() - INTERVAL '$1 hours'
		ORDER BY detected_at DESC
	`, hours)

	if err != nil {
		return nil, fmt.Errorf("failed to query recent anomalies: %w", err)
	}
	defer rows.Close()

	var results []models.Anomaly
	for rows.Next() {
		var anomaly models.Anomaly
		if err := rows.Scan(&anomaly.ID, &anomaly.Type, &anomaly.Parameter, &anomaly.Value,
			&anomaly.Latitude, &anomaly.Longitude, &anomaly.DetectedAt,
			&anomaly.AirQualityDataID, &anomaly.AirQualityDataTimestamp); err != nil {
			return nil, err
		}
		results = append(results, anomaly)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
