import React, { useState } from 'react';
import AirQualityMap from '../components/AirQualityMap';
import AnomalyAlertPanel from '../components/AnomalyAlertPanel';
import RegionDetail from '../components/RegionDetail';
import HistoricalChart from '../components/HistoricalChart';
import { useWebSocket } from '../hooks/useWebSocket';

const Home = () => {
  const [selectedRegion, setSelectedRegion] = useState(null);
  const [alerts, setAlerts] = useState([]);
  
  // Connect to WebSocket for real-time anomaly alerts
  const { lastMessage } = useWebSocket(
    process.env.REACT_APP_WS_URL || 'ws://localhost:8081/ws/alerts',
    (message) => {
      // Add new alert to state
      const newAlert = {
        ...message,
        id: Date.now(), // Use timestamp as id
        createdAt: new Date()
      };
      
      setAlerts(prevAlerts => [newAlert, ...prevAlerts]);
    }
  );

  const handleMapClick = (region) => {
    setSelectedRegion(region);
  };

  return (
    <div className="container">
      <div className="row mb-4">
        <div className="col-md-8">
          <div className="card">
            <div className="card-header">
              <h5>Real-time Air Quality Map</h5>
            </div>
            <div className="card-body">
              <AirQualityMap onRegionSelect={handleMapClick} />
            </div>
          </div>
        </div>
        
        <div className="col-md-4">
          <div className="card">
            <div className="card-header">
              <h5>Anomaly Alerts</h5>
            </div>
            <div className="card-body alert-panel">
              <AnomalyAlertPanel alerts={alerts} />
            </div>
          </div>
        </div>
      </div>
      
      {selectedRegion && (
        <div className="row">
          <div className="col-md-5">
            <div className="card">
              <div className="card-header">
                <h5>Region Details</h5>
              </div>
              <div className="card-body">
                <RegionDetail region={selectedRegion} />
              </div>
            </div>
          </div>
          
          <div className="col-md-7">
            <div className="card">
              <div className="card-header">
                <h5>Historical Data</h5>
              </div>
              <div className="card-body">
                <HistoricalChart region={selectedRegion} />
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Home; 