import client from "./client";

const API = "/storage";
const PUBLIC_API = "/public";

// ── File management (authenticated) ──────────────────────────────────────────

export const uploadFile = (formData, onProgress) =>
  client.post(`${API}/upload`, formData, {
    headers: { "Content-Type": "multipart/form-data" },
    onUploadProgress: (e) =>
      onProgress && onProgress(Math.round((e.loaded * 100) / e.total)),
  });

export const getFiles = (params) => client.get(API, { params });
export const getFileById = (id) => client.get(`${API}/${id}`);
export const deleteFile = (id) => client.delete(`${API}/${id}`);
export const downloadOwnFile = (id) =>
  client.get(`${API}/${id}/download`, { responseType: "blob" });

// ── Share link management (authenticated) ─────────────────────────────────────

export const createShareLink = (fileId, data) =>
  client.post(`${API}/${fileId}/share`, data);
export const getShareLinks = (fileId) => client.get(`${API}/${fileId}/shares`);
export const updateShareLink = (shareId, data) =>
  client.put(`${API}/shares/${shareId}`, data);
export const revokeShareLink = (shareId) =>
  client.delete(`${API}/shares/${shareId}`);
export const getShareLinkLogs = (shareId) =>
  client.get(`${API}/shares/${shareId}/logs`);

// ── Public access (no auth required) ─────────────────────────────────────────

export const getPublicFileInfo = (token) =>
  client.get(`${PUBLIC_API}/share/${token}`);
export const downloadViaShareLink = (token, password = "") =>
  client.get(`${PUBLIC_API}/share/${token}/download`, {
    params: password ? { password } : {},
    responseType: "blob",
  });
