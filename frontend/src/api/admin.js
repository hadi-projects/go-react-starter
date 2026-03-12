import apiClient from './client';

export const getUsers = async (page = 1, limit = 10, search = '') => {
    const response = await apiClient.get(`/users?page=${page}&limit=${limit}&search=${search}`);
    return response.data;
};
export const exportUsers = async (format = 'excel') => {
    return apiClient.get(`/users/export?format=${format}`, {
        responseType: 'blob',
    });
};
export const getMe = async () => {
    const response = await apiClient.get('/users/me');
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

export const getRoles = async (page = 1, limit = 10, search = '') => {
    const response = await apiClient.get(`/roles?page=${page}&limit=${limit}&search=${search}`);
    return response.data;
};

export const createRole = async (data) => {
    const response = await apiClient.post('/roles', data);
    return response.data;
};

export const updateRole = async (id, data) => {
    const response = await apiClient.put(`/roles/${id}`, data);
    return response.data;
};

export const deleteRole = async (id) => {
    const response = await apiClient.delete(`/roles/${id}`);
    return response.data;
};

export const exportRoles = async (format = 'excel') => {
    return apiClient.get(`/roles/export?format=${format}`, {
        responseType: 'blob',
    });
};

export const getPermissions = async (page = 1, limit = 10, search = '') => {
    const response = await apiClient.get(`/permissions?page=${page}&limit=${limit}&search=${search}`);
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

export const exportPermissions = async (format = 'excel') => {
    return apiClient.get(`/permissions/export?format=${format}`, {
        responseType: 'blob',
    });
};

// Cache management
export const getCacheStatus = async () => {
    const response = await apiClient.get('/cache/status');
    return response.data;
};

export const getHealthStatus = async () => {
    const response = await apiClient.get('/health/status');
    return response.data;
};

export const clearCache = async () => {
    const response = await apiClient.delete('/cache/clear');
    return response.data;
};

// Module Generator
export const generateModule = async (data) => {
    const response = await apiClient.post('/generator', data);
    return response.data;
};

export const exportAdmin = async (format = 'excel') => {
    return apiClient.get(`/admin/module/export?format=${format}`, {
        responseType: 'blob',
    });
};
