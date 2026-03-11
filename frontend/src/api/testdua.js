import axios from './client';

const API_PATH = '/testdua';

export const getAllTestduas = async (params) => {
    const response = await axios.get(API_PATH, { params });
    return response.data;
};

export const getTestduaById = async (id) => {
    const response = await axios.get(`${API_PATH}/${id}`);
    return response.data;
};

export const createTestdua = async (data) => {
    const response = await axios.post(API_PATH, data);
    return response.data;
};

export const updateTestdua = async (id, data) => {
    const response = await axios.put(`${API_PATH}/${id}`, data);
    return response.data;
};

export const deleteTestdua = async (id) => {
    const response = await axios.delete(`${API_PATH}/${id}`);
    return response.data;
};
