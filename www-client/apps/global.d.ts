interface AppConfig {
	TECHNICAL_WORK: boolean;
	YANDEX_METRIKA_ID?: number | string;
}

interface Window {
	__APP_CONFIG__: AppConfig;
	ym?: ((...args: unknown[]) => void) & { a?: unknown[][]; l?: number };
}
