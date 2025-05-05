# Air Quality Monitoring Platform Frontend

A React-based frontend for the Air Quality Monitoring Platform. This application provides a real-time visualization of air quality data including a heatmap, historical charts, and anomaly alerts.

## Features

- 🌍 **Interactive Map**: Shows air pollution levels on a world map with color-coded intensity based on pollution levels.
- 📈 **Historical Charts**: Line charts showing time-series data for a selected region.
- 🚨 **Real-time Alerts**: Displays alerts when anomalies are detected via WebSocket.
- 🔍 **Region Details**: Provides detailed stats when a region is selected on the map.

## Prerequisites

- Node.js (v14+)
- npm or yarn
- Mapbox account for map visualization (get a free token at https://mapbox.com)

## Environment Variables

Create a `.env` file in the frontend directory with the following variables:

```
REACT_APP_INGEST_API_URL=http://localhost:8080
REACT_APP_NOTIFIER_API_URL=http://localhost:8081
REACT_APP_WEBSOCKET_URL=ws://localhost:8081/ws/alerts
REACT_APP_MAPBOX_TOKEN=your_mapbox_access_token_here
```

## Setup and Running (Standalone)

1. Install dependencies:
   ```
   npm install
   ```

2. Start the development server:
   ```
   npm start
   ```

3. Access the application at http://localhost:3000

## Building for Production

```
npm run build
```

This creates a `build` folder with optimized production files.

## Running with Docker

### Build and run only the frontend:

```
docker build -t air-quality-frontend .
docker run -p 80:80 air-quality-frontend
```

### Run the entire stack (frontend + backend):

From the project root directory:

```
docker-compose up -d
```

This will start both the frontend and backend services, with the frontend accessible at http://localhost:80

## Project Structure

- `src/components/`: React components
  - `AirQualityMap.js`: Mapbox-based heatmap
  - `AnomalyAlertPanel.js`: Real-time anomaly alert display
  - `HistoricalChart.js`: Chart.js time-series charts
  - `RegionDetail.js`: Detailed view for selected regions
- `src/services/`: API service functions
- `src/hooks/`: Custom React hooks
  - `useWebSocket.js`: Hook for WebSocket connection
- `src/utils/`: Utility functions
- `public/`: Static files

## Technologies Used

- React.js
- Mapbox GL JS for mapping
- Chart.js for data visualization
- WebSocket for real-time communication
- Axios for HTTP requests
- Docker for containerization 