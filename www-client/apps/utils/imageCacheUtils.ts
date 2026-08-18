import { CACHEKEYs } from '@/constants/CacheKeys.constants';

const CACHE_NAME = 'paylist-images-v1';
const MAX_MISSES = 400;

const objectUrls = new Map<string, string>();
const ready = new Set<string>();
const inflight = new Map<string, Promise<string>>();
const misses = new Set<string>(readMisses());

const isBrowser = typeof window !== 'undefined';

function readMisses(): string[] {
	if (typeof localStorage === 'undefined') {
		return [];
	}

	try {
		const raw = localStorage.getItem(CACHEKEYs.IMAGE_MISS);
		const parsed = raw ? (JSON.parse(raw) as unknown) : [];

		return Array.isArray(parsed) ? parsed.filter((item): item is string => typeof item === 'string') : [];
	} catch {
		return [];
	}
}

function persistMisses(): void {
	if (typeof localStorage === 'undefined') {
		return;
	}

	try {
		localStorage.setItem(CACHEKEYs.IMAGE_MISS, JSON.stringify([...misses].slice(-MAX_MISSES)));
	} catch {
		return;
	}
}

const rememberMiss = (src: string): void => {
	if (!src || misses.has(src)) {
		return;
	}

	misses.add(src);
	persistMisses();
};

export const isImageFailed = (src?: string | null): boolean => Boolean(src && misses.has(src));

export const getCachedObjectUrl = (src?: string | null): string | null => {
	if (!src) {
		return null;
	}

	return objectUrls.get(src) ?? null;
};

export const isImageReady = (src?: string | null): boolean => Boolean(src && (ready.has(src) || objectUrls.has(src)));

const preloadHtmlImage = (src: string): Promise<boolean> =>
	new Promise((resolve) => {
		const image = new Image();
		image.decoding = 'async';
		image.onload = () => resolve(true);
		image.onerror = () => resolve(false);
		image.src = src;
	});

const cacheResponse = async (src: string): Promise<string> => {
	const cache = await caches.open(CACHE_NAME);
	let response = await cache.match(src);

	if (!response) {
		response = await fetch(src, { mode: 'cors', credentials: 'omit', cache: 'force-cache' });

		if (!response.ok) {
			rememberMiss(src);
			throw new Error('image-http');
		}

		await cache.put(src, response.clone());
	}

	const blob = await response.blob();
	const objectUrl = URL.createObjectURL(blob);
	objectUrls.set(src, objectUrl);
	ready.add(src);

	return objectUrl;
};

export const resolveCachedImage = (src: string): Promise<string> => {
	if (!src) {
		return Promise.reject(new Error('empty'));
	}

	if (misses.has(src)) {
		return Promise.reject(new Error('miss'));
	}

	const cached = objectUrls.get(src);
	if (cached) {
		return Promise.resolve(cached);
	}

	const pending = inflight.get(src);
	if (pending) {
		return pending;
	}

	const task = (async () => {
		if (isBrowser && 'caches' in window) {
			try {
				return await cacheResponse(src);
			} catch (error) {
				if (misses.has(src)) {
					throw error;
				}
			}
		}

		const ok = await preloadHtmlImage(src);
		if (!ok) {
			rememberMiss(src);
			throw new Error('image-load');
		}

		ready.add(src);
		return src;
	})().finally(() => {
		inflight.delete(src);
	});

	inflight.set(src, task);
	return task;
};

export const markImageFailed = (src?: string | null): void => {
	if (src) {
		rememberMiss(src);
	}
};
