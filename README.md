# Air Quality Monitoring Platform

A comprehensive platform for real-time air quality monitoring, with microservice architecture and React frontend.

## Project Structure

```
air-quality-platform/
├── backend/
│   ├── cmd/                  # Entry points for each microservice
│   │   ├── ingest/           # Ingest API service
│   │   ├── processor/        # Data processor service
│   │   └── notifier/         # Notifier WebSocket service
│   ├── build/                # Docker build files for each service
│   │   ├── ingest/           # Dockerfile for ingest service
│   │   ├── processor/        # Dockerfile for processor service
│   │   └── notifier/         # Dockerfile for notifier service
│   ├── db/                   # Database scripts and migrations
│   │   └── init/             # Initial schema setup
│   └── docker-compose.yml    # Backend service definitions
├── frontend/
│   ├── src/
│   │   ├── components/       # Reusable UI components
│   │   ├── pages/            # Page components 
│   │   ├── hooks/            # Custom React hooks
│   │   ├── services/         # API service functions
│   │   └── utils/            # Utility functions
│   ├── public/               # Static assets
│   ├── Dockerfile            # Frontend Docker configuration
│   └── nginx.conf            # Nginx configuration for the frontend
├── docker-compose.yml        # Root compose file that combines frontend and backend
└── README.md                 # You are here
```



## Quick Start

```bash
# Build and run all services
docker-compose up --build

# Access:
# - Frontend: http://localhost:80
# - Ingest API: http://localhost:8080
# - Notifier API: http://localhost:8081
```

## Testing the Application

1. **Send test data to the ingest API**:

```bash
curl -X POST http://localhost:8080/api/data \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": 41.015,
    "longitude": 28.979,
    "parameter": "PM2.5",
    "value": 90.0,
    "timestamp": "2023-05-02T13:45:00Z"
  }'
```

2. **Check recent anomalies**:

```bash
curl http://localhost:8081/api/anomalies
```

3. **View any anomalies in the frontend anomaly panel**

## Architecture

- **Frontend**: React.js with Mapbox for visualization
- **Backend**: Go microservices
  - **Ingest Service**: Receives air quality data via REST API
  - **Processor Service**: Analyzes data for anomalies
  - **Notifier Service**: Provides WebSocket real-time updates
- **Data Flow**: REST API → Kafka → TimescaleDB → WebSocket

## Features

### Backend Services

- **Ingest Service**: Accepts incoming air quality data from sensors
- **Processor Service**: Internal service for anomaly detection 
- **Notifier Service**: Provides WebSocket connections and historical anomaly data

### Frontend Features

- **Air Quality Map or Table**: Shows air pollution levels on a world map or in a table view
- **Historical Charts**: Line charts showing time-series data for a selected region
- **Real-time Alerts**: Displays alerts when anomalies are detected via WebSocket
- **Region Details**: Provides detailed stats when a region is selected

## Development

### Running Individual Services

To run the backend services only:
```
cd backend
docker-compose up
```

To run the frontend in development mode:
```
cd frontend
npm install
npm start
```

## Technologies Used

- **Backend**: Go with microservices architecture
- **Frontend**: React.js, Mapbox GL JS, Chart.js
- **Communication**: RESTful APIs, WebSockets
- **Data Storage**: TimescaleDB (PostgreSQL for time-series data)
- **Message Queue**: Kafka
- **Containerization**: Docker, Docker Compose

## Troubleshooting

### Common Issues

- **Missing Mapbox Token**: The frontend will show a table view instead of a map
- **Connection Refused Errors**: Make sure all services are running by checking `docker-compose ps`
- **Empty Data**: Send some test data using the curl commands above

## License

MIT 
