import axios from './client';

const API_PATH = '/blog';

export const getAllBlogs = async (params) => {
    const response = await axios.get(API_PATH, { params });
    return response.data;
};

export const getBlogById = async (id) => {
    const response = await axios.get(`${API_PATH}/${id}`);
    return response.data;
};

export const createBlog = async (data) => {
    const response = await axios.post(API_PATH, data);
    return response.data;
};

export const updateBlog = async (id, data) => {
    const response = await axios.put(`${API_PATH}/${id}`, data);
    return response.data;
};

export const deleteBlog = async (id) => {
    const response = await axios.delete(`${API_PATH}/${id}`);
    return response.data;
};

export const exportBlog = async (format = 'excel') => {
    return axios.get(`${API_PATH}/export?format=${format}`, {
        responseType: 'blob',
    });
};
