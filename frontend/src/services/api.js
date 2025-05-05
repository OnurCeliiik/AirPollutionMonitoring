import axios from 'axios';

const ingestApiUrl = process.env.REACT_APP_API_URL || process.env.REACT_APP_INGEST_API_URL || 'http://localhost:8080';
const notifierApiUrl = process.env.REACT_APP_NOTIFIER_API_URL || 'http://localhost:8081';

// Create axios instances for both services
const ingestApi = axios.create({
  baseURL: ingestApiUrl,
  headers: {
    'Content-Type': 'application/json',
  },
});

const notifierApi = axios.create({
  baseURL: notifierApiUrl,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Error handler helper
const handleApiError = (error) => {
  if (error.response) {
    console.error('API Error:', error.response.data);
    return Promise.reject(error.response.data);
  }
  console.error('API Error:', error.message);
  return Promise.reject({ error: 'Network or server error' });
};

// API functions for Ingest Service
export const ingestService = {
  // Submit air quality data
  submitData: async (data) => {
    try {
      const response = await ingestApi.post('/api/data', data);
      return response.data;
    } catch (error) {
      return handleApiError(error);
    }
  },
  
  // Check health of ingest service
  healthCheck: async () => {
    try {
      const response = await ingestApi.get('/health');
      return response.data;
    } catch (error) {
      return handleApiError(error);
    }
  },
};

// API functions for Notifier Service
export const notifierService = {
  // Get recent anomalies with optional hours parameter
  getRecentAnomalies: async (hours = 24) => {
    try {
      const response = await notifierApi.get('/api/anomalies', {
        params: { hours },
      });
      return response.data;
    } catch (error) {
      return handleApiError(error);
    }
  },
  
  // Check health of notifier service
  healthCheck: async () => {
    try {
      const response = await notifierApi.get('/health');
      return response.data;
    } catch (error) {
      return handleApiError(error);
    }
  },
}; 