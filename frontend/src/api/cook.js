import axios from './client';

const API_PATH = '/cook';

export const getAllCooks = async (params) => {
    const response = await axios.get(API_PATH, { params });
    return response.data;
};

export const getCookById = async (id) => {
    const response = await axios.get(`${API_PATH}/${id}`);
    return response.data;
};

export const createCook = async (data) => {
    const response = await axios.post(API_PATH, data);
    return response.data;
};

export const updateCook = async (id, data) => {
    const response = await axios.put(`${API_PATH}/${id}`, data);
    return response.data;
};

export const deleteCook = async (id) => {
    const response = await axios.delete(`${API_PATH}/${id}`);
    return response.data;
};

export const exportCook = async (format = 'excel') => {
    return axios.get(`${API_PATH}/export?format=${format}`, {
        responseType: 'blob',
    });
};
