import client from './client';

const logApi = {
    getLogs: async (params) => {
        const response = await client.get('/logs', { params });
        return response.data;
    },
    getHttpLogs: async (params) => {
        const response = await client.get('/logs/http', { params });
        return response.data;
    },
};

export default logApi;
