const fromEnv = Number(import.meta.env.VITE_YANDEX_METRIKA_ID);
const fromRuntime = Number(typeof window !== 'undefined' ? window.__APP_CONFIG__?.YANDEX_METRIKA_ID : 0);

/** Counter id from index.html — Direct goals are bound to this counter. */
export const YANDEX_METRIKA_INLINE_ID = 111658785;

export const YANDEX_METRIKA_ID = fromEnv || fromRuntime || 0;

export const YANDEX_METRIKA_PAID_GOAL = 'premium_paid';

export const YANDEX_METRIKA_DEFERRED_PATHS = ['/paid'] as const;
