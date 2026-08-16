import { YANDEX_METRIKA_ID } from '@/constants/analytics.constants';

const SCRIPT_SRC = 'https://mc.yandex.ru/metrika/tag.js';
const YANDEX_COOKIE_PREFIXES = ['_ym_', 'ymex', '_yasc'];

type YmStub = ((...args: unknown[]) => void) & { a?: unknown[][]; l?: number };

let initialized = false;
let lastHitUrl = '';

const getYm = (): YmStub | undefined => window.ym as YmStub | undefined;

const ensureYmStub = (): void => {
	if (typeof window.ym === 'function') {
		return;
	}

	const ym: YmStub = (...args: unknown[]) => {
		ym.a = ym.a || [];
		ym.a.push(args);
	};

	ym.l = Date.now();
	window.ym = ym;
};

const injectScript = (): void => {
	if (document.querySelector(`script[src="${SCRIPT_SRC}"]`)) {
		return;
	}

	const script = document.createElement('script');
	script.async = true;
	script.src = SCRIPT_SRC;
	document.head.appendChild(script);
};

export const initYandexMetrika = (): void => {
	if (initialized || !YANDEX_METRIKA_ID) {
		return;
	}

	initialized = true;
	lastHitUrl = `${window.location.pathname}${window.location.search}`;

	ensureYmStub();
	injectScript();
	getYm()?.(YANDEX_METRIKA_ID, 'init', {
		clickmap: true,
		trackLinks: true,
		accurateTrackBounce: true,
		webvisor: false,
	});
};

export const hitYandexMetrika = (url: string): void => {
	if (!initialized || !YANDEX_METRIKA_ID || url === lastHitUrl) {
		return;
	}

	lastHitUrl = url;
	getYm()?.(YANDEX_METRIKA_ID, 'hit', url);
};

export const clearYandexMetrikaCookies = (): void => {
	const hostname = window.location.hostname;
	const domains = ['', hostname, `.${hostname}`];

	document.cookie.split(';').forEach((cookie) => {
		const name = cookie.split('=')[0]?.trim();

		if (!name || !YANDEX_COOKIE_PREFIXES.some((prefix) => name.startsWith(prefix))) {
			return;
		}

		domains.forEach((domain) => {
			const domainPart = domain ? `;domain=${domain}` : '';
			document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/${domainPart}`;
		});
	});
};
