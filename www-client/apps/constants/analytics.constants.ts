const fromEnv = Number(import.meta.env.VITE_YANDEX_METRIKA_ID);
const fromRuntime = Number(typeof window !== 'undefined' ? window.__APP_CONFIG__?.YANDEX_METRIKA_ID : 0);

export const YANDEX_METRIKA_ID = fromEnv || fromRuntime || 0;
