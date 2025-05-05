import React, { useEffect, useState } from 'react';
import { timeAgo, getColorForValue, getQualityDescription } from '../utils/airQualityUtils';

const AnomalyAlertPanel = ({ alerts }) => {
  const [filteredAlerts, setFilteredAlerts] = useState([]);

  // Filter alerts to only keep those less than 1 hour old
  useEffect(() => {
    const oneHourAgo = new Date();
    oneHourAgo.setHours(oneHourAgo.getHours() - 1);
    
    const filtered = alerts.filter(alert => {
      const alertTime = new Date(alert.createdAt || alert.timestamp || alert.detected_at);
      return alertTime >= oneHourAgo;
    });
    
    setFilteredAlerts(filtered);
  }, [alerts]);

  // Get severity class based on anomaly type
  const getSeverityClass = (type) => {
    switch (type) {
      case 'ThresholdExceeded':
      case 'SpikeDetected':
        return 'danger';
      case 'StatisticalOutlier':
      case 'GeographicInconsistency':
        return 'warning';
      default:
        return '';
    }
  };

  // Get human-readable description of anomaly type
  const getAnomalyDescription = (type) => {
    switch (type) {
      case 'ThresholdExceeded':
        return 'Exceeded WHO guidelines';
      case 'StatisticalOutlier':
        return 'Statistical anomaly detected';
      case 'SpikeDetected':
        return 'Sudden increase in pollutant';
      case 'GeographicInconsistency':
        return 'Inconsistent with nearby readings';
      default:
        return type;
    }
  };

  // If no alerts, show empty state
  if (filteredAlerts.length === 0) {
    return (
      <div className="empty-state">
        <p>No active alerts in the last hour</p>
      </div>
    );
  }

  return (
    <div>
      {filteredAlerts.map(alert => {
        const alertTime = new Date(alert.createdAt || alert.timestamp || alert.detected_at);
        const parameterKey = alert.parameter?.replace('.', '') || 'PM25';
        
        return (
          <div 
            key={alert.id} 
            className={`alert-item ${getSeverityClass(alert.type)}`}
            style={{ borderLeftColor: getColorForValue(parameterKey, alert.value) }}
          >
            <div className="alert-header">
              <strong>{alert.parameter || 'Unknown'}: {alert.value?.toFixed(1)}</strong>
              <span className="badge" style={{ 
                backgroundColor: getColorForValue(parameterKey, alert.value),
                color: alert.value > 100 ? 'white' : 'black'
              }}>
                {getQualityDescription(parameterKey, alert.value)}
              </span>
            </div>
            
            <div className="alert-body">
              <p>{getAnomalyDescription(alert.type)}</p>
              <p className="location">
                <small>
                  Location: [{alert.latitude?.toFixed(3) || alert.location?.[0]?.toFixed(3)}, 
                  {alert.longitude?.toFixed(3) || alert.location?.[1]?.toFixed(3)}]
                </small>
              </p>
            </div>
            
            <div className="alert-footer">
              <small className="text-muted">
                {timeAgo(alertTime)}
              </small>
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default AnomalyAlertPanel; 