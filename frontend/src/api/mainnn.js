import axios from './client';

const API_PATH = '/mainnn';

export const getAllMainnns = async (params) => {
    const response = await axios.get(API_PATH, { params });
    return response.data;
};

export const getMainnnById = async (id) => {
    const response = await axios.get(`${API_PATH}/${id}`);
    return response.data;
};

export const createMainnn = async (data) => {
    const response = await axios.post(API_PATH, data);
    return response.data;
};

export const updateMainnn = async (id, data) => {
    const response = await axios.put(`${API_PATH}/${id}`, data);
    return response.data;
};

export const deleteMainnn = async (id) => {
    const response = await axios.delete(`${API_PATH}/${id}`);
    return response.data;
};

export const exportMainnn = async (format = 'excel') => {
    return axios.get(`${API_PATH}/export?format=${format}`, {
        responseType: 'blob',
    });
};
