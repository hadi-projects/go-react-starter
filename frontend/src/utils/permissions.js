export const PERMS = {
  // Core
  CREATE_USER: 1n << 0n,
  DELETE_USER: 1n << 3n,
  EDIT_USER: 1n << 4n,
  GET_USER: 1n << 5n,
  CREATE_ROLE: 1n << 6n,
  DELETE_ROLE: 1n << 7n,
  EDIT_ROLE: 1n << 8n,
  GET_ROLE: 1n << 9n,
  CREATE_PERMISSION: 1n << 10n,
  DELETE_PERMISSION: 1n << 11n,
  EDIT_PERMISSION: 1n << 12n,
  GET_PERMISSION: 1n << 13n,
  MANAGE_CACHE: 1n << 14n,
  GET_ALL_LOGS: 1n << 15n,
  CREATE_MODULE: 1n << 16n,
  GET_AUDIT_LOG: 1n << 17n,
  GET_AUTH_LOG: 1n << 18n,
  GET_OWN_LOGS: 1n << 19n,
  GET_PROFILE: 1n << 20n,

  // Logs
  GET_HTTP_LOG: 1n << 29n,

  // Produk
  GET_PRODUK: 1n << 30n,
  CREATE_PRODUK: 1n << 31n,
  UPDATE_PRODUK: 1n << 32n,
  DELETE_PRODUK: 1n << 33n,

  // Settings
  GET_SETTING: 1n << 34n,
  EDIT_SETTING: 1n << 35n,

  // Storage
  UPLOAD_FILE: 1n << 58n,
  GET_FILE: 1n << 59n,
  DELETE_FILE: 1n << 60n,
  SHARE_FILE: 1n << 61n,
  MANAGE_STORAGE: 1n << 62n,
};
