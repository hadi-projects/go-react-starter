import axios from './client';

const API_PATH = '/wisuda';

export const getAllWisudas = async (params) => {
    const response = await axios.get(API_PATH, { params });
    return response.data;
};

export const getWisudaById = async (id) => {
    const response = await axios.get(`${API_PATH}/${id}`);
    return response.data;
};

export const createWisuda = async (data) => {
    const response = await axios.post(API_PATH, data);
    return response.data;
};

export const updateWisuda = async (id, data) => {
    const response = await axios.put(`${API_PATH}/${id}`, data);
    return response.data;
};

export const deleteWisuda = async (id) => {
    const response = await axios.delete(`${API_PATH}/${id}`);
    return response.data;
};

export const exportWisuda = async (format = 'excel') => {
    return axios.get(`${API_PATH}/export?format=${format}`, {
        responseType: 'blob',
    });
};
