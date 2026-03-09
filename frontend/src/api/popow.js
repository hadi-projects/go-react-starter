import axios from './client';

const API_PATH = '/api/v1/popow';

export const getAllPopows = async (params) => {
    const response = await axios.get(API_PATH, { params });
    return response.data;
};

export const getPopowById = async (id) => {
    const response = await axios.get(`${API_PATH}/${id}`);
    return response.data;
};

export const createPopow = async (data) => {
    const response = await axios.post(API_PATH, data);
    return response.data;
};

export const updatePopow = async (id, data) => {
    const response = await axios.put(`${API_PATH}/${id}`, data);
    return response.data;
};

export const deletePopow = async (id) => {
    const response = await axios.delete(`${API_PATH}/${id}`);
    return response.data;
};
