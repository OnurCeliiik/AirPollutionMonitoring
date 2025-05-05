import React from 'react';
import { formatTimestamp, getQualityDescription, getColorForValue } from '../utils/airQualityUtils';

const RegionDetail = ({ region }) => {
  if (!region) {
    return (
      <div className="empty-state">
        <p>Select a region on the map to view details</p>
      </div>
    );
  }

  const parameterKey = region.parameter?.replace('.', '') || 'PM25';
  const color = getColorForValue(parameterKey, region.value);
  const quality = getQualityDescription(parameterKey, region.value);
  
  return (
    <div className="region-detail">
      <div className="location-info">
        <h4>Location</h4>
        <p>
          Latitude: {region.latitude?.toFixed(5) || 'N/A'}<br />
          Longitude: {region.longitude?.toFixed(5) || 'N/A'}
        </p>
      </div>
      
      <div className="measurement-info">
        <h4>Latest Measurement</h4>
        <div className="measurement-value" style={{ color }}>
          <strong>{region.parameter}: {region.value?.toFixed(2)}</strong>
        </div>
        <div className="quality-indicator" style={{ backgroundColor: color, color: region.value > 100 ? 'white' : 'black' }}>
          {quality}
        </div>
        <p>
          Detected at: {formatTimestamp(region.detected_at)}
        </p>
      </div>
      
      {region.type && (
        <div className="anomaly-info">
          <h4>Anomaly Information</h4>
          <p><strong>Type:</strong> {region.type}</p>
          <p><strong>ID:</strong> {region.id}</p>
          {region.air_quality_data_id && (
            <p><strong>Data Source ID:</strong> {region.air_quality_data_id}</p>
          )}
        </div>
      )}
      
      <div className="recommendations">
        <h4>Health Recommendations</h4>
        {quality === 'Good' && (
          <p>Air quality is considered satisfactory, and air pollution poses little or no risk.</p>
        )}
        {quality === 'Moderate' && (
          <p>Air quality is acceptable; however, for some pollutants there may be a moderate health concern for a very small number of people.</p>
        )}
        {quality === 'Unhealthy' && (
          <p>Members of sensitive groups may experience health effects. The general public is not likely to be affected.</p>
        )}
        {quality === 'Very Unhealthy' && (
          <p>Health warnings of emergency conditions. The entire population is more likely to be affected.</p>
        )}
        {quality === 'Hazardous' && (
          <p>Health alert: everyone may experience more serious health effects.</p>
        )}
      </div>
    </div>
  );
};

export default RegionDetail; 