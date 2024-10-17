import axios from 'axios';

const apiClient = axios.create({
  baseURL: 'http://localhost:8080', // Change to your API base URL
  headers: {
    'Content-Type': 'application/json',
  },
});

// Optionally add interceptors
apiClient.interceptors.request.use(
  config => {
    // Add authorization token or other modifications here if needed
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = token;
    }
    return config;
  },
  error => Promise.reject(error)
);

export default apiClient;