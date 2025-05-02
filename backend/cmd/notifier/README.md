# Anomaly Notifier Service

This service is responsible for consuming anomaly alerts from Kafka and broadcasting them to WebSocket clients in real-time.

## Responsibilities

- Consume anomaly alerts from the `anomaly-alerts` Kafka topic
- Maintain WebSocket connections with clients
- Broadcast detected anomalies to connected clients in real-time
- Provide API endpoints for retrieving historical anomalies

## Configuration

The service can be configured using the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | HTTP port to listen on | 8081 |
| KAFKA_BROKERS | Comma-separated list of Kafka brokers | localhost:9092 |
| DB_CONNECTION_STRING | TimescaleDB connection string | postgres://postgres:postgres@localhost:5432/timescaledb?sslmode=disable |
| ENVIRONMENT | Environment (development/production) | development |
| LOG_LEVEL | Logging level (DEBUG, INFO, WARN, ERROR, FATAL) | INFO |
| ALLOWED_ORIGINS | CORS allowed origins | * |

## API Endpoints

### GET /api/anomalies

Retrieves recent anomalies.

**Query Parameters:**
- `hours`: Number of hours of history to return (default: 24)

**Response:**
```json
[
  {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "type": "ThresholdExceeded",
    "parameter": "PM2.5",
    "value": 90.0,
    "latitude": 41.015,
    "longitude": 28.979,
    "detected_at": "2025-05-02T13:45:00Z"
  }
]
```

### GET /health

Health check endpoint.

**Response:**
```json
{
  "status": "up",
  "service": "notifier"
}
```

### WebSocket: /ws/alerts

WebSocket endpoint for real-time anomaly alerts.

**Message Format:**
```json
{
  "parameter": "PM2.5",
  "value": 90.0,
  "type": "ThresholdExceeded",
  "location": [41.015, 28.979],
  "timestamp": "2025-05-02T13:45:00Z"
}
```

## Main Components

- `main.go`: Service entry point that sets up Kafka consumer, WebSocket hub, and HTTP server
- `internal/services/websocket/websocket.go`: WebSocket server and client management
- `internal/db/timescaledb.go`: Database access layer for TimescaleDB

## See Also

See the [main README](../../README.md) for system-wide instructions and architecture overview. 