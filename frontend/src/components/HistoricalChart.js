import React, { useEffect, useState } from 'react';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
} from 'chart.js';
import { notifierService } from '../services/api';
import { getColorForValue } from '../utils/airQualityUtils';

// Register required Chart.js components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

const HistoricalChart = ({ region }) => {
  const [chartData, setChartData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [timeRange, setTimeRange] = useState(24); // Default to 24 hours

  // When region changes or time range changes, fetch historical data
  useEffect(() => {
    if (!region) return;
    
    const fetchHistoricalData = async () => {
      setLoading(true);
      setError(null);
      
      try {
        // Fetch anomalies for the selected region
        const anomalies = await notifierService.getRecentAnomalies(timeRange);
        
        // Filter anomalies for the selected region and parameter
        const filteredAnomalies = anomalies.filter(anomaly => (
          anomaly.parameter === region.parameter &&
          Math.abs(anomaly.latitude - region.latitude) < 0.01 &&
          Math.abs(anomaly.longitude - region.longitude) < 0.01
        ));
        
        if (filteredAnomalies.length === 0) {
          // If no matching anomalies, create mock data for demonstration
          // In a real app, you would need to fetch actual historical data for this location
          createMockData(region);
        } else {
          prepareChartData(filteredAnomalies);
        }
        
        setLoading(false);
      } catch (err) {
        console.error('Error fetching historical data:', err);
        setError('Failed to load historical data');
        setLoading(false);
      }
    };
    
    fetchHistoricalData();
  }, [region, timeRange]);

  // Prepare data for Chart.js
  const prepareChartData = (data) => {
    // Sort data by timestamp
    const sortedData = [...data].sort((a, b) => 
      new Date(a.detected_at) - new Date(b.detected_at)
    );
    
    const labels = sortedData.map(item => {
      const date = new Date(item.detected_at);
      return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    });
    
    const values = sortedData.map(item => item.value);
    const parameterKey = region.parameter?.replace('.', '') || 'PM25';
    
    const chartData = {
      labels,
      datasets: [
        {
          label: region.parameter,
          data: values,
          borderColor: getColorForValue(parameterKey, Math.max(...values)),
          backgroundColor: 'rgba(255, 255, 255, 0.2)',
          fill: false,
          tension: 0.4,
        }
      ]
    };
    
    setChartData(chartData);
  };
  
  // Create mock data for demonstration purposes
  const createMockData = (region) => {
    const now = new Date();
    const labels = [];
    const values = [];
    const baseValue = region.value || 50;
    
    // Generate 24 data points for the last 24 hours
    for (let i = timeRange; i >= 0; i--) {
      const time = new Date(now);
      time.setHours(time.getHours() - i);
      
      labels.push(time.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }));
      
      // Generate slightly random values around the base value
      const randomVariation = (Math.random() - 0.5) * 20;
      values.push(Math.max(0, baseValue + randomVariation));
    }
    
    const parameterKey = region.parameter?.replace('.', '') || 'PM25';
    
    const chartData = {
      labels,
      datasets: [
        {
          label: region.parameter,
          data: values,
          borderColor: getColorForValue(parameterKey, Math.max(...values)),
          backgroundColor: 'rgba(255, 255, 255, 0.2)',
          fill: false,
          tension: 0.4,
        }
      ]
    };
    
    setChartData(chartData);
  };

  const handleTimeRangeChange = (e) => {
    setTimeRange(parseInt(e.target.value, 10));
  };

  if (!region) {
    return (
      <div className="empty-state">
        <p>Select a region on the map to view historical data</p>
      </div>
    );
  }

  return (
    <div className="historical-chart">
      <div className="chart-controls">
        <div className="form-group">
          <label htmlFor="timeRange">Time Range: </label>
          <select 
            id="timeRange" 
            value={timeRange} 
            onChange={handleTimeRangeChange}
            className="form-control"
          >
            <option value={6}>Last 6 hours</option>
            <option value={12}>Last 12 hours</option>
            <option value={24}>Last 24 hours</option>
            <option value={48}>Last 48 hours</option>
            <option value={72}>Last 72 hours</option>
          </select>
        </div>
      </div>
      
      {loading && <div className="loading">Loading chart data...</div>}
      {error && <div className="error">{error}</div>}
      
      {chartData && (
        <div className="chart-container">
          <Line 
            data={chartData} 
            options={{
              responsive: true,
              plugins: {
                legend: {
                  position: 'top',
                },
                title: {
                  display: true,
                  text: `Historical ${region.parameter} Values`,
                },
                tooltip: {
                  callbacks: {
                    label: function(context) {
                      return `${context.dataset.label}: ${context.parsed.y.toFixed(1)}`;
                    }
                  }
                }
              },
              scales: {
                y: {
                  beginAtZero: true,
                  title: {
                    display: true,
                    text: 'Value'
                  }
                },
                x: {
                  title: {
                    display: true,
                    text: 'Time'
                  }
                }
              }
            }}
          />
        </div>
      )}
      
      <div className="chart-info">
        <p><small>Note: Chart shows {chartData?.labels?.length || 0} data points over the selected time period</small></p>
      </div>
    </div>
  );
};

export default HistoricalChart; 