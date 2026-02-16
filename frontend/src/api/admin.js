import apiClient from './client';

export const getUsers = async (page = 1, limit = 10) => {
    const response = await apiClient.get(`/users?page=${page}&limit=${limit}`);
    return response.data;
};

export const createUser = async (data) => {
    const response = await apiClient.post('/auth/register', data);
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
