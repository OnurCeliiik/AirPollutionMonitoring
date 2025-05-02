# Real-Time Air Quality Monitoring Platform

This is a backend system for monitoring air pollution in real-time. It receives environmental data, processes it asynchronously, performs anomaly detection, and exposes data through RESTful APIs and WebSockets.

## Architecture

```
┌─────────────┐     HTTP     ┌──────────────┐      ┌────────────┐
│ Data Sources│────POST─────►│ Ingest Service│────►│            │
└─────────────┘              └──────────────┘      │            │
                                   │               │   Kafka    │
                                   │               │            │
                                   ▼               │            │
┌─────────────┐               ┌──────────────┐     │            │
│  WebSocket  │◄──────────────┤Notifier Service◄───┤            │
│   Clients   │  Alerts       └──────────────┘     └────────────┘
└─────────────┘                     ▲                     ▲
                                    │                     │
                                    │                     │
                                    │               ┌──────────────┐
                                    └───────────────┤ Processor    │
                                                    │   Service    │
                                                    └──────────────┘
                                                          │
                                                          │
                                                          ▼
                                                    ┌────────────┐
                                                    │TimescaleDB │
                                                    └────────────┘
```

The platform consists of three microservices:

1. **API Gateway / Ingest Service** - Accepts pollution data from scripts or UI and publishes raw data to Kafka.
2. **Data Processor Service** - Subscribes to raw air data, validates, detects anomalies, and inserts data into TimescaleDB.
3. **Anomaly Notifier Service** - Subscribes to anomaly alerts, stores anomalies in DB, and sends WebSocket notifications to frontend clients.

## Tech Stack

- **Language**: Go
- **Architecture**: Microservices
- **Message Queue**: Apache Kafka (for async communication)
- **Database**: TimescaleDB (time-series PostgreSQL)
- **Containerization**: Docker
- **API Documentation**: Swagger (via swaggo/gin-swagger)
- **WebSocket**: For real-time anomaly alerts

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21 or higher (for local development)
- Bash (for running test scripts)

### Running with Docker

From the root project directory:

```bash
docker-compose up --build
```

This will start all services including TimescaleDB, Kafka, and the three microservices.

### API Endpoints

#### Ingest Service (port 8080)

- `POST /api/data` - Submit air quality data
  ```json
  {
    "latitude": 41.015,
    "longitude": 28.979,
    "parameter": "PM2.5",
    "value": 90.0,
    "timestamp": "2025-05-02T13:45:00Z"
  }
  ```
- `GET /health` - Health check endpoint
- `GET /swagger/index.html` - Swagger UI for API documentation

#### Notifier Service (port 8081)

- `GET /api/anomalies` - Get recent anomalies
- `GET /health` - Health check endpoint
- `WS /ws/alerts` - WebSocket endpoint for real-time anomaly alerts

### Testing

#### Running Tests

```bash
# Run unit tests
cd backend
go test ./...

# Run integration tests (requires running services)
go test -tags=integration ./...
```

#### Using Test Scripts

You can test the platform using the provided scripts in a Docker container:

```bash
# Run manual input test
docker-compose run test ./manual-input.sh 41.015 28.979 "PM2.5" 90.0

# Run automated test
docker-compose run test ./auto-test.sh --duration=60 --rate=10 --anomaly-chance=20
```

Parameters for auto-test.sh:
- `duration`: Test duration in seconds (default: 60)
- `rate`: Requests per minute (default: 10)
- `anomaly-chance`: Percentage chance of generating anomalous values (default: 20)

## Database Schema

### TimescaleDB Tables

#### air_quality_data
| Column     | Type       |
|------------|------------|
| id         | UUID       |
| latitude   | FLOAT      |
| longitude  | FLOAT      |
| parameter  | TEXT       |
| value      | FLOAT      |
| timestamp  | TIMESTAMPTZ|

#### anomalies
| Column     | Type       |
|------------|------------|
| id         | UUID       |
| type       | TEXT       |
| parameter  | TEXT       |
| value      | FLOAT      |
| latitude   | FLOAT      |
| longitude  | FLOAT      |
| detected_at| TIMESTAMPTZ|

## Anomaly Detection Logic

The platform applies one or more of the following anomaly detection methods:

1. **Thresholds** - Compares values against WHO limits (e.g., PM2.5 > 15 μg/m³)
2. **Z-score / Statistical Outliers** - Identifies values that are statistically significant outliers
3. **Spike Detection** - Identifies sudden increases (>50% higher than 24h average)
4. **Geographic Inconsistency** - Identifies values that differ significantly from nearby readings

## Troubleshooting

### Common Issues

#### Kafka Not Starting
- Verify Zookeeper is running properly: `docker-compose logs zookeeper`
- Check Kafka logs: `docker-compose logs kafka`
- Ensure ports 9092 and 2181 are not in use by other applications

#### TimescaleDB Issues
- Verify the database is running: `docker-compose logs timescaledb`
- Check connection string in service environment variables
- Ensure port 5432 is not in use by another PostgreSQL instance

#### Services Failing to Start
- Check that all services can connect to Kafka: `docker-compose logs ingest`
- Verify TimescaleDB is healthy: `docker-compose exec timescaledb pg_isready`
- Check service health endpoints: `curl http://localhost:8080/health`

### Resetting the System
To completely reset the system and start fresh:
```bash
docker-compose down -v
docker-compose up --build
```

## License

MIT 