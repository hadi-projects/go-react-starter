import apiClient from './client';

export const login = async (credentials) => {
    const response = await apiClient.post('/auth/login', credentials);
    return response.data;
};

export const register = async (userData) => {
    const response = await apiClient.post('/auth/register', userData);
    return response.data;
};

export const forgotPassword = async (email) => {
    const response = await apiClient.post('/auth/forgot-password', { email });
    return response.data;
};

export const resetPassword = async (data) => {
    const response = await apiClient.post('/auth/reset-password', data);
    return response.data;
};

export const logoutApi = async (reason) => {
    try {
        await apiClient.post('/auth/logout', { reason });
    } catch (error) {
        // Silent fail on logging
    }
};

export const refreshTokenApi = async (refreshToken) => {
    const response = await apiClient.post('/auth/refresh', { refresh_token: refreshToken });
    return response.data;
};
