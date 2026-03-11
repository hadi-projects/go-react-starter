import axios from './client';

const API_PATH = '/testsaja';

export const getAllTestsajas = async (params) => {
    const response = await axios.get(API_PATH, { params });
    return response.data;
};

export const getTestsajaById = async (id) => {
    const response = await axios.get(`${API_PATH}/${id}`);
    return response.data;
};

export const createTestsaja = async (data) => {
    const response = await axios.post(API_PATH, data);
    return response.data;
};

export const updateTestsaja = async (id, data) => {
    const response = await axios.put(`${API_PATH}/${id}`, data);
    return response.data;
};

export const deleteTestsaja = async (id) => {
    const response = await axios.delete(`${API_PATH}/${id}`);
    return response.data;
};
