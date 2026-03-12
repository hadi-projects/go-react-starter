import axios from './client';

const API_PATH = '/produk';

export const getAllProduks = async (params) => {
    const response = await axios.get(API_PATH, { params });
    return response.data;
};

export const getProdukById = async (id) => {
    const response = await axios.get(`${API_PATH}/${id}`);
    return response.data;
};

export const createProduk = async (data) => {
    const response = await axios.post(API_PATH, data);
    return response.data;
};

export const updateProduk = async (id, data) => {
    const response = await axios.put(`${API_PATH}/${id}`, data);
    return response.data;
};

export const deleteProduk = async (id) => {
    const response = await axios.delete(`${API_PATH}/${id}`);
    return response.data;
};

export const exportProduk = async (format = 'excel') => {
    return axios.get(`${API_PATH}/export?format=${format}`, {
        responseType: 'blob',
    });
};
