import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

const apiClient = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Request interceptor to add auth token
apiClient.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// Response interceptor for error handling
apiClient.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;
        
        // If 401 and not from login and not already retried
        if (error.response?.status === 401 && !originalRequest._retry && !originalRequest.url.includes('/auth/login')) {
            originalRequest._retry = true;
            
            const refreshToken = localStorage.getItem('refresh_token');
            if (refreshToken && !originalRequest.url.includes('/auth/refresh')) {
                try {
                    const { refreshTokenApi } = await import('./auth');
                    const data = await refreshTokenApi(refreshToken);
                    const newToken = data.data.access_token;
                    
                    localStorage.setItem('token', newToken);
                    // Update header and retry
                    originalRequest.headers.Authorization = `Bearer ${newToken}`;
                    return apiClient(originalRequest);
                } catch (refreshError) {
                    // Refresh failed, clear everything and login
                    console.error('Refresh token failed', refreshError);
                }
            }

            // If we're here, either no refresh token or refresh failed
            // Attempt to log logout before clearing token
            try {
                const logoutApi = (await import('./auth')).logoutApi;
                await logoutApi('system');
            } catch (err) {
                // Ignore errors
            }
            // Clear token and redirect to login
            localStorage.removeItem('token');
            localStorage.removeItem('refresh_token');
            localStorage.removeItem('user');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

export default apiClient;
