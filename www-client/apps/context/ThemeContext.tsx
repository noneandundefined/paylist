import { CACHEKEYs } from '@/constants/CacheKeys.constants';
import { applyTheme, getStoredTheme, resolveTheme, type Theme } from '@/utils/theme';
import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from 'react';

interface ThemeContextType {
	theme: Theme;
	isDark: boolean;
	setTheme: (theme: Theme) => void;
	toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextType | null>(null);

export const ThemeProvider = ({ children }: { children: ReactNode }) => {
	const [theme, setThemeState] = useState<Theme>(() => resolveTheme());

	useEffect(() => {
		applyTheme(theme);
	}, [theme]);

	useEffect(() => {
		const media = window.matchMedia('(prefers-color-scheme: dark)');

		const handleChange = () => {
			if (getStoredTheme()) {
				return;
			}

			setThemeState(media.matches ? 'dark' : 'light');
		};

		media.addEventListener('change', handleChange);
		return () => media.removeEventListener('change', handleChange);
	}, []);

	const setTheme = useCallback((nextTheme: Theme) => {
		localStorage.setItem(CACHEKEYs.THEME, nextTheme);
		setThemeState(nextTheme);
		applyTheme(nextTheme);
	}, []);

	const toggleTheme = useCallback(() => {
		setTheme(theme === 'dark' ? 'light' : 'dark');
	}, [setTheme, theme]);

	const value = useMemo(
		() => ({
			theme,
			isDark: theme === 'dark',
			setTheme,
			toggleTheme,
		}),
		[setTheme, theme, toggleTheme]
	);

	return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>;
};

export const useTheme = () => {
	const context = useContext(ThemeContext);

	if (!context) {
		throw new Error('useTheme must be used inside ThemeProvider');
	}

	return context;
};
