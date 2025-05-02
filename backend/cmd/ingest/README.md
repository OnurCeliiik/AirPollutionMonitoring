# API Gateway / Ingest Service

This service acts as the entry point for air quality data, accepting measurements via REST API and publishing them to Kafka for further processing.

## Responsibilities

- Accept HTTP POST requests with air quality measurements
- Validate incoming data
- Publish valid data to the `raw-air-data` Kafka topic
- Provide Swagger documentation for API endpoints

## Configuration

The service can be configured using the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | HTTP port to listen on | 8080 |
| KAFKA_BROKERS | Comma-separated list of Kafka brokers | localhost:9092 |
| ENVIRONMENT | Environment (development/production) | development |
| LOG_LEVEL | Logging level (DEBUG, INFO, WARN, ERROR, FATAL) | INFO |
| ALLOWED_ORIGINS | CORS allowed origins | * |

## API Endpoints

### POST /api/data

Submits a new air quality data point.

**Request:**
```json
{
  "latitude": 41.015,
  "longitude": 28.979,
  "parameter": "PM2.5",
  "value": 90.0,
  "timestamp": "2025-05-02T13:45:00Z"
}
```

**Response:**
```json
{
  "message": "Data received and queued for processing",
  "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
}
```

### GET /health

Health check endpoint.

**Response:**
```json
{
  "status": "up",
  "service": "ingest"
}
```

### GET /swagger/index.html

Swagger UI documentation (available in development mode).

## Main Components

- `main.go`: Service entry point that configures and starts the HTTP server
- `internal/api/ingest_handler.go`: Handler for the data ingest API endpoint

## See Also

See the [main README](../../README.md) for system-wide instructions and architecture overview. 