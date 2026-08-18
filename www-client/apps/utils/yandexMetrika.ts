import { YANDEX_METRIKA_DEFERRED_PATHS, YANDEX_METRIKA_ID, YANDEX_METRIKA_INLINE_ID, YANDEX_METRIKA_PAID_GOAL } from '@/constants/analytics.constants';
import { CACHEKEYs } from '@/constants/CacheKeys.constants';

const SCRIPT_SRC = 'https://mc.yandex.ru/metrika/tag.js';
const YANDEX_COOKIE_PREFIXES = ['_ym_', 'ymex', '_yasc'];

type YmStub = ((...args: unknown[]) => void) & { a?: unknown[][]; l?: number };

let initialized = false;
let lastHitUrl = '';
let activeCounterId = 0;
let paidConversionUnlocked = false;

const getYm = (): YmStub | undefined => window.ym as YmStub | undefined;

const metrikaCounterId = (): number => YANDEX_METRIKA_ID || YANDEX_METRIKA_INLINE_ID;

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
	if (document.querySelector(`script[src="${SCRIPT_SRC}"]`) || document.querySelector(`script[src^="${SCRIPT_SRC}"]`)) {
		return;
	}

	const script = document.createElement('script');
	script.async = true;
	script.src = SCRIPT_SRC;
	document.head.appendChild(script);
};

const readPaidGoalSent = (): boolean => {
	try {
		return sessionStorage.getItem(CACHEKEYs.PAID_PAGE_GOAL) === '1';
	} catch {
		return false;
	}
};

const writePaidGoalSent = (): void => {
	try {
		sessionStorage.setItem(CACHEKEYs.PAID_PAGE_GOAL, '1');
	} catch {
		return;
	}
};

export const isYandexMetrikaDeferredPath = (pathname: string): boolean => {
	return (YANDEX_METRIKA_DEFERRED_PATHS as readonly string[]).includes(pathname) && !paidConversionUnlocked;
};

export const unlockPaidPageMetrika = (): void => {
	paidConversionUnlocked = true;
};

export const initYandexMetrika = (url?: string, counterId = YANDEX_METRIKA_ID): void => {
	if (initialized || !counterId) {
		return;
	}

	if (!url && isYandexMetrikaDeferredPath(window.location.pathname)) {
		return;
	}

	initialized = true;
	activeCounterId = counterId;
	lastHitUrl = url ?? `${window.location.pathname}${window.location.search}`;

	ensureYmStub();
	injectScript();
	getYm()?.(counterId, 'init', {
		clickmap: true,
		trackLinks: true,
		accurateTrackBounce: true,
		webvisor: false,
		url: lastHitUrl,
	});
};

export const hitYandexMetrika = (url: string): void => {
	if (!initialized || !activeCounterId || url === lastHitUrl) {
		return;
	}

	lastHitUrl = url;
	getYm()?.(activeCounterId, 'hit', url);
};

export const trackPaidPageConversion = (): void => {
	if (readPaidGoalSent()) {
		unlockPaidPageMetrika();
		return;
	}

	unlockPaidPageMetrika();
	writePaidGoalSent();

	if (!initialized) {
		initYandexMetrika('/paid', metrikaCounterId());
	} else {
		hitYandexMetrika('/paid');
	}

	getYm()?.(activeCounterId || metrikaCounterId(), 'reachGoal', YANDEX_METRIKA_PAID_GOAL);
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
