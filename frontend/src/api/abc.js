import axios from './client';

const API_PATH = '/api/v1/abc';

export const getAllabcs = async (params) => {
    const response = await axios.get(API_PATH, { params });
    return response.data;
};

export const getabcById = async (id) => {
    const response = await axios.get(`${API_PATH}/${id}`);
    return response.data;
};

export const createabc = async (data) => {
    const response = await axios.post(API_PATH, data);
    return response.data;
};

export const updateabc = async (id, data) => {
    const response = await axios.put(`${API_PATH}/${id}`, data);
    return response.data;
};

export const deleteabc = async (id) => {
    const response = await axios.delete(`${API_PATH}/${id}`);
    return response.data;
};
