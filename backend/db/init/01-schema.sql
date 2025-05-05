-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- Create air quality data table
CREATE TABLE IF NOT EXISTS air_quality_data (
    id UUID,
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL,
    parameter TEXT NOT NULL,
    value FLOAT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (id, timestamp)
);

-- Convert to TimescaleDB hypertable
SELECT create_hypertable('air_quality_data', 'timestamp', if_not_exists => TRUE);

-- Create anomalies table
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

-- Convert to TimescaleDB hypertable
SELECT create_hypertable('anomalies', 'detected_at', if_not_exists => TRUE);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_air_quality_location ON air_quality_data (latitude, longitude);
CREATE INDEX IF NOT EXISTS idx_air_quality_parameter ON air_quality_data (parameter);
CREATE INDEX IF NOT EXISTS idx_anomalies_type ON anomalies (type);
CREATE INDEX IF NOT EXISTS idx_anomalies_parameter ON anomalies (parameter); 