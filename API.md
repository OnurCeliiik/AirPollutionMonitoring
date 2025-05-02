# Air Pollution Monitoring Platform API Documentation

This document provides detailed information about the Air Pollution Monitoring Platform's API endpoints, WebSocket connections, and data models to assist frontend developers.

## Table of Contents

- [Service Architecture](#service-architecture)
- [Ingest Service](#ingest-service)
- [Notifier Service](#notifier-service)
- [Data Models](#data-models)
- [Error Handling](#error-handling)
- [WebSocket Communication](#websocket-communication)

## Service Architecture

The platform consists of three microservices:

1. **Ingest Service** (Port 8080): Accepts incoming air quality data from sensors
2. **Processor Service**: Internal service for anomaly detection (no external API)
3. **Notifier Service** (Port 8081): Provides websocket connections and historical anomaly data

## Ingest Service

Base URL: `http://localhost:8080` (development) or your production domain

### Submit Air Quality Data

Submits new air quality measurement data from sensors.

- **URL**: `/api/data`
- **Method**: `POST`
- **Content-Type**: `application/json`

**Request Body**:
```json
{
  "latitude": 41.015,
  "longitude": 28.979,
  "parameter": "PM2.5",
  "value": 90.0,
  "timestamp": "2023-05-02T13:45:00Z"
}
```

**Parameters**:
| Field | Type | Description | Required |
|-------|------|-------------|----------|
| latitude | float | Latitude coordinate (range: -90 to 90) | Yes |
| longitude | float | Longitude coordinate (range: -180 to 180) | Yes |
| parameter | string | Measurement parameter (PM2.5, PM10, O3, etc.) | Yes |
| value | float | Measurement value | Yes |
| timestamp | string (ISO8601) | Time of measurement | Yes |

**Success Response**:
- **Code**: 201 CREATED
- **Content**:
```json
{
  "message": "Data received and queued for processing",
  "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
}
```

**Error Responses**:
- **Code**: 400 BAD REQUEST
  - Invalid data format or missing required fields
- **Code**: 500 INTERNAL SERVER ERROR
  - Server processing error

### Health Check

Check if the ingest service is operational.

- **URL**: `/health`
- **Method**: `GET`

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "status": "up",
  "service": "ingest"
}
```

## Notifier Service

Base URL: `http://localhost:8081` (development) or your production domain

### Get Recent Anomalies

Retrieve historical anomaly data for a specified time period.

- **URL**: `/api/anomalies`
- **Method**: `GET`

**Query Parameters**:
| Parameter | Type | Description | Required | Default |
|-----------|------|-------------|----------|---------|
| hours | integer | Number of hours of history to retrieve | No | 24 |

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
[
  {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "type": "ThresholdExceeded",
    "parameter": "PM2.5",
    "value": 90.0,
    "latitude": 41.015,
    "longitude": 28.979,
    "detected_at": "2023-05-02T13:45:00Z"
  },
  {
    "id": "7ca8c921-0eae-22e2-91b5-11d15fe541d9",
    "type": "StatisticalOutlier",
    "parameter": "O3",
    "value": 120.5,
    "latitude": 41.025,
    "longitude": 28.990,
    "detected_at": "2023-05-02T14:10:00Z"
  }
]
```

**Error Response**:
- **Code**: 500 INTERNAL SERVER ERROR
  - Database connection error or other server issue

### Health Check

Check if the notifier service is operational.

- **URL**: `/health`
- **Method**: `GET`

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "status": "up",
  "service": "notifier"
}
```

## WebSocket Communication

The platform provides real-time anomaly notifications via WebSocket.

### Connect to Anomaly WebSocket

- **URL**: `ws://localhost:8081/ws/alerts` (development) or your production domain with `wss://`
- **Protocol**: WebSocket

**Connection Process**:
1. Establish WebSocket connection
2. Listen for incoming messages

**Message Format** (received from server):
```json
{
  "parameter": "PM2.5",
  "value": 90.0,
  "type": "ThresholdExceeded",
  "location": [41.015, 28.979],
  "timestamp": "2023-05-02T13:45:00Z"
}
```

### Anomaly Types

| Type | Description |
|------|-------------|
| ThresholdExceeded | Value exceeds WHO guidelines |
| StatisticalOutlier | Value is a statistical outlier (z-score) |
| SpikeDetected | Sudden increase in pollutant levels |
| GeographicInconsistency | Reading inconsistent with nearby sensors |

## Data Models

### Air Quality Data

```json
{
  "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "parameter": "PM2.5",
  "value": 90.0,
  "latitude": 41.015,
  "longitude": 28.979,
  "timestamp": "2023-05-02T13:45:00Z"
}
```

### Anomaly Alert

```json
{
  "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "type": "ThresholdExceeded",
  "parameter": "PM2.5",
  "value": 90.0,
  "latitude": 41.015,
  "longitude": 28.979,
  "detected_at": "2023-05-02T13:45:00Z",
  "air_quality_data_id": "7ca8c921-0eae-22e2-91b5-11d15fe541d9"
}
```

## Error Handling

All API endpoints return standard HTTP status codes:

- 200: Success
- 201: Created
- 400: Bad Request (client error, invalid input)
- 404: Not Found (endpoint does not exist)
- 500: Internal Server Error

Error responses contain a JSON body with error details:

```json
{
  "error": "Error description message"
}
```

## CORS Support

The API supports Cross-Origin Resource Sharing with the following headers:

```
Access-Control-Allow-Origin: * (configurable in production)
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```

## Security Notes

- In production environments, use HTTPS/WSS instead of HTTP/WS
- Configure ALLOWED_ORIGINS environment variable to restrict CORS access
- WebSocket connections do not require authentication in the current version 