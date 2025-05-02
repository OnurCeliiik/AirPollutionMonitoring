# Data Processor Service

This service is responsible for consuming air quality data from Kafka, performing anomaly detection, and storing the data in TimescaleDB.

## Responsibilities

- Consume air quality data from the `raw-air-data` Kafka topic
- Perform data validation
- Run anomaly detection algorithms
- Store air quality data in TimescaleDB
- Store detected anomalies in TimescaleDB
- Publish detected anomalies to the `anomaly-alerts` Kafka topic

## Configuration

The service can be configured using the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| KAFKA_BROKERS | Comma-separated list of Kafka brokers | localhost:9092 |
| DB_CONNECTION_STRING | TimescaleDB connection string | postgres://postgres:postgres@localhost:5432/timescaledb?sslmode=disable |
| ENVIRONMENT | Environment (development/production) | development |
| LOG_LEVEL | Logging level (DEBUG, INFO, WARN, ERROR, FATAL) | INFO |

## Anomaly Detection

The service uses multiple methods to detect anomalies:

1. **Threshold Exceedance**: Compares values against WHO limits.
2. **Statistical Outlier Detection**: Uses Z-score to identify statistical outliers.
3. **Spike Detection**: Identifies sudden increases in values.
4. **Geographic Inconsistency**: Identifies values that differ from nearby readings.

## Main Components

- `main.go`: Service entry point that sets up Kafka consumer and connects to the database
- `internal/services/anomaly/detector.go`: Implements anomaly detection algorithms
- `internal/db/timescaledb.go`: Database access layer for TimescaleDB

## See Also

See the [main README](../../README.md) for system-wide instructions and architecture overview. 