import { CACHEKEYs } from '@/constants/CacheKeys.constants';

export type CookieConsentValue = 'accepted' | 'rejected';

export const COOKIE_PREFERENCES_EVENT = 'paylist:cookie-preferences';
export const COOKIE_CONSENT_CHANGE_EVENT = 'paylist:cookie-consent-change';

export const getCookieConsent = (): CookieConsentValue | null => {
	try {
		const value = localStorage.getItem(CACHEKEYs.COOKIE_CONSENT);

		if (value === 'accepted' || value === 'rejected') {
			return value;
		}
	} catch {
		return null;
	}

	return null;
};

export const setCookieConsent = (value: CookieConsentValue): void => {
	try {
		localStorage.setItem(CACHEKEYs.COOKIE_CONSENT, value);
	} catch {
		return;
	}

	window.dispatchEvent(new CustomEvent<CookieConsentValue>(COOKIE_CONSENT_CHANGE_EVENT, { detail: value }));
};

export const openCookiePreferences = (): void => {
	window.dispatchEvent(new Event(COOKIE_PREFERENCES_EVENT));
};
