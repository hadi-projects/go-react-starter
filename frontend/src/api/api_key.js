import api from './client';

export const getApiKeys = (params) => {
    return api.get('/apikeys', { params });
};

export const createApiKey = (data) => {
    return api.post('/apikeys', data);
};

export const deleteApiKey = (id) => {
    return api.delete(`/apikeys/${id}`);
};
