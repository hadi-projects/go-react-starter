import axios from './client';

const API_PATH = '/arsip';

export const getAllArsips = async (params) => {
    const response = await axios.get(API_PATH, { params });
    return response.data;
};

export const getArsipById = async (id) => {
    const response = await axios.get(`${API_PATH}/${id}`);
    return response.data;
};

export const createArsip = async (data) => {
    const response = await axios.post(API_PATH, data);
    return response.data;
};

export const updateArsip = async (id, data) => {
    const response = await axios.put(`${API_PATH}/${id}`, data);
    return response.data;
};

export const deleteArsip = async (id) => {
    const response = await axios.delete(`${API_PATH}/${id}`);
    return response.data;
};

export const exportArsip = async (format = 'excel') => {
    return axios.get(`${API_PATH}/export?format=${format}`, {
        responseType: 'blob',
    });
};
