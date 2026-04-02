import apiClient from './client';

export const getSettingsByCategory = async (category) => {
  const response = await apiClient.get(`/settings/${category}`);
  return response.data;
};

export const updateSettings = async (settings) => {
  const response = await apiClient.put('/settings', { settings });
  return response.data;
};
