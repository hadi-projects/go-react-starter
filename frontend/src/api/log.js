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
    exportLogs: async (format = 'excel') => {
        return client.get(`/logs/export?format=${format}`, {
            responseType: 'blob',
        });
    },
    exportHttpLogs: async (format = 'excel') => {
        return client.get(`/logs/http/export?format=${format}`, {
            responseType: 'blob',
        });
    },
    exportSystemLogs: async (format = 'excel') => {
        return client.get(`/logs/system/export?format=${format}`, {
            responseType: 'blob',
        });
    },
    exportAuditLogs: async (format = 'excel') => {
        return client.get(`/logs/audit/export?format=${format}`, {
            responseType: 'blob',
        });
    },
};

export default logApi;
