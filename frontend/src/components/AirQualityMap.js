import React, { useEffect, useState, useRef } from 'react';
import mapboxgl from 'mapbox-gl';
import 'mapbox-gl/dist/mapbox-gl.css';
import { notifierService } from '../services/api';
import { getColorForValue } from '../utils/airQualityUtils';

// Set your Mapbox token here
mapboxgl.accessToken = process.env.REACT_APP_MAPBOX_TOKEN || '';

const AirQualityMap = ({ onRegionSelect }) => {
  const mapContainer = useRef(null);
  const map = useRef(null);
  const [airQualityData, setAirQualityData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [noToken, setNoToken] = useState(!mapboxgl.accessToken);

  // Initialize map when component mounts
  useEffect(() => {
    if (map.current) return; // initialize map only once
    
    // If no token, don't try to initialize the map
    if (!mapboxgl.accessToken) {
      setLoading(false);
      setNoToken(true);
      return;
    }
    
    map.current = new mapboxgl.Map({
      container: mapContainer.current,
      style: 'mapbox://styles/mapbox/light-v11',
      center: [0, 20], // Default center (middle of the map)
      zoom: 1.5
    });
    
    // Add navigation controls
    map.current.addControl(new mapboxgl.NavigationControl(), 'top-right');
    
    // Add event listener for when the map has finished loading
    map.current.on('load', () => {
      // After map loads, fetch air quality data
      fetchAnomalyData();
    });
    
    return () => {
      if (map.current) {
        map.current.remove();
        map.current = null;
      }
    };
  }, []);

  // Fetch anomaly data from the API
  const fetchAnomalyData = async () => {
    setLoading(true);
    try {
      const data = await notifierService.getRecentAnomalies(24);
      setAirQualityData(data);
      if (map.current && map.current.loaded()) {
        updateMapData(data);
      }
      setLoading(false);
    } catch (err) {
      console.error('Error fetching air quality data:', err);
      setError('Failed to load air quality data');
      setLoading(false);
    }
  };

  // Update map with new data
  const updateMapData = (data) => {
    if (!map.current || !map.current.loaded()) return;
    
    // Remove previous layers and sources if they exist
    if (map.current.getSource('air-quality-data')) {
      map.current.removeLayer('air-quality-heat');
      map.current.removeLayer('air-quality-point');
      map.current.removeSource('air-quality-data');
    }
    
    // Prepare GeoJSON data for the map
    const geoJsonData = {
      type: 'FeatureCollection',
      features: data.map(item => ({
        type: 'Feature',
        properties: {
          id: item.id,
          parameter: item.parameter,
          value: item.value,
          color: getColorForValue(item.parameter.replace('.', ''), item.value),
          detected_at: item.detected_at
        },
        geometry: {
          type: 'Point',
          coordinates: [item.longitude, item.latitude]
        }
      }))
    };
    
    // Add a new source with the data
    map.current.addSource('air-quality-data', {
      type: 'geojson',
      data: geoJsonData
    });
    
    // Add a heatmap layer
    map.current.addLayer({
      id: 'air-quality-heat',
      type: 'heatmap',
      source: 'air-quality-data',
      paint: {
        'heatmap-weight': 1,
        'heatmap-intensity': 1,
        'heatmap-color': [
          'interpolate',
          ['linear'],
          ['heatmap-density'],
          0, 'rgba(0, 255, 0, 0)',
          0.2, 'rgba(0, 255, 0, 0.5)',
          0.4, 'rgba(255, 255, 0, 0.5)',
          0.6, 'rgba(255, 128, 0, 0.5)',
          0.8, 'rgba(255, 0, 0, 0.5)',
          1, 'rgba(153, 0, 0, 0.5)'
        ],
        'heatmap-radius': 30,
        'heatmap-opacity': 0.7
      }
    });
    
    // Add a circle layer for point data
    map.current.addLayer({
      id: 'air-quality-point',
      type: 'circle',
      source: 'air-quality-data',
      paint: {
        'circle-radius': 6,
        'circle-color': ['get', 'color'],
        'circle-opacity': 0.9,
        'circle-stroke-width': 1,
        'circle-stroke-color': '#fff'
      }
    });
    
    // Add click event to points
    map.current.on('click', 'air-quality-point', (e) => {
      if (e.features.length > 0) {
        const feature = e.features[0];
        const props = feature.properties;
        const coordinates = feature.geometry.coordinates.slice();
        
        // Find full data object for selected point
        const fullDataPoint = data.find(item => item.id === props.id);
        
        if (fullDataPoint && onRegionSelect) {
          onRegionSelect({
            ...fullDataPoint,
            coordinates
          });
        }
      }
    });

    // Change cursor to pointer when hovering over a point
    map.current.on('mouseenter', 'air-quality-point', () => {
      map.current.getCanvas().style.cursor = 'pointer';
    });
    
    map.current.on('mouseleave', 'air-quality-point', () => {
      map.current.getCanvas().style.cursor = '';
    });
  };

  // Refresh data periodically
  useEffect(() => {
    if (noToken) {
      // Even without a map, we can still fetch data
      fetchAnomalyData();
    }
    
    const interval = setInterval(() => {
      if (!noToken || airQualityData.length > 0) {
        fetchAnomalyData();
      }
    }, 60000); // Refresh every minute
    
    return () => clearInterval(interval);
  }, [noToken, airQualityData]);

  // If no token, render a table instead of a map
  if (noToken) {
    return (
      <div>
        <div className="no-token-message">
          <p><strong>Mapbox token not found!</strong> Add a token to your .env file to see the map.</p>
          <p>For testing, here's a list of the latest air quality data:</p>
        </div>
        
        {loading ? (
          <div className="loading">Loading data...</div>
        ) : error ? (
          <div className="error">{error}</div>
        ) : airQualityData.length === 0 ? (
          <div className="empty-state">No air quality data available</div>
        ) : (
          <div className="data-table-container">
            <table className="data-table">
              <thead>
                <tr>
                  <th>Parameter</th>
                  <th>Value</th>
                  <th>Location</th>
                  <th>Time</th>
                  <th>Action</th>
                </tr>
              </thead>
              <tbody>
                {airQualityData.map(item => (
                  <tr key={item.id}>
                    <td>{item.parameter}</td>
                    <td style={{ color: getColorForValue(item.parameter.replace('.', ''), item.value) }}>
                      {item.value.toFixed(1)}
                    </td>
                    <td>
                      [{item.latitude.toFixed(3)}, {item.longitude.toFixed(3)}]
                    </td>
                    <td>{new Date(item.detected_at).toLocaleString()}</td>
                    <td>
                      <button 
                        onClick={() => onRegionSelect(item)}
                        className="select-button"
                      >
                        Select
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    );
  }

  return (
    <div>
      <div ref={mapContainer} className="map-container" />
      {loading && <div className="loading">Loading map data...</div>}
      {error && <div className="error">{error}</div>}
    </div>
  );
};

export default AirQualityMap; 