import { CACHEKEYs } from '@/constants/CacheKeys.constants';

export type Theme = 'light' | 'dark';

export const getStoredTheme = (): Theme | null => {
	const stored = localStorage.getItem(CACHEKEYs.THEME);

	if (stored === 'light' || stored === 'dark') {
		return stored;
	}

	return null;
};

export const getSystemTheme = (): Theme => {
	if (typeof window === 'undefined') {
		return 'light';
	}

	return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
};

export const resolveTheme = (theme?: Theme | null): Theme => theme ?? getStoredTheme() ?? getSystemTheme();

export const applyTheme = (theme: Theme) => {
	if (typeof document === 'undefined') {
		return;
	}

	document.documentElement.classList.toggle('dark', theme === 'dark');
	document.documentElement.style.colorScheme = theme;
};

export const persistTheme = (theme: Theme) => {
	localStorage.setItem(CACHEKEYs.THEME, theme);
	applyTheme(theme);
};

applyTheme(resolveTheme());
