import axios from './client';

const API_PATH = '/news';

export const getAllNewss = async (params) => {
    const response = await axios.get(API_PATH, { params });
    return response.data;
};

export const getNewsById = async (id) => {
    const response = await axios.get(`${API_PATH}/${id}`);
    return response.data;
};

export const createNews = async (data) => {
    const response = await axios.post(API_PATH, data);
    return response.data;
};

export const updateNews = async (id, data) => {
    const response = await axios.put(`${API_PATH}/${id}`, data);
    return response.data;
};

export const deleteNews = async (id) => {
    const response = await axios.delete(`${API_PATH}/${id}`);
    return response.data;
};

export const exportNews = async (format = 'excel') => {
    return axios.get(`${API_PATH}/export?format=${format}`, {
        responseType: 'blob',
    });
};
