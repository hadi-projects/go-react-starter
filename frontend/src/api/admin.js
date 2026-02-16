import apiClient from './client';

export const getUsers = async (page = 1, limit = 10) => {
    const response = await apiClient.get(`/users?page=${page}&limit=${limit}`);
    return response.data;
};
// User API
export const createUser = async (data) => {
    const response = await apiClient.post('/users', data);
    return response.data;
};

export const updateUser = async (id, data) => {
    const response = await apiClient.put(`/users/${id}`, data);
    return response.data;
};

export const deleteUser = async (id) => {
    const response = await apiClient.delete(`/users/${id}`);
    return response.data;
};

export const getRoles = async (page = 1, limit = 10) => {
    const response = await apiClient.get(`/roles?page=${page}&limit=${limit}`);
    return response.data;
};

export const getPermissions = async (page = 1, limit = 10) => {
    const response = await apiClient.get(`/permissions?page=${page}&limit=${limit}`);
    return response.data;
};

export const createPermission = async (data) => {
    const response = await apiClient.post('/permissions', data);
    return response.data;
};

export const updatePermission = async (id, data) => {
    const response = await apiClient.put(`/permissions/${id}`, data);
    return response.data;
};

export const deletePermission = async (id) => {
    const response = await apiClient.delete(`/permissions/${id}`);
    return response.data;
};

// Cache management
export const clearCache = async () => {
    const response = await apiClient.delete('/cache/clear');
    return response.data;
};
