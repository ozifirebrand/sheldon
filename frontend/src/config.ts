export const API_BASE = typeof import.meta.env?.VITE_API_URL === 'string'
  ? import.meta.env.VITE_API_URL.replace(/\/$/, '')
  : ''
