import client from './client';

export const getCacheStatus = async () => {
    try {
        const response = await client.get('/cache/status');
        return response.data;
    } catch (error) {
        throw error;
    }
};
