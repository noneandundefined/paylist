import { CACHEKEYs } from '@/constants/CacheKeys.constants';

export const readAuthSession = (): boolean => Boolean(localStorage.getItem(CACHEKEYs.L_SESSION));

export const clearAuthSession = () => {
	localStorage.removeItem(CACHEKEYs.L_SESSION);
	localStorage.removeItem(CACHEKEYs.AUTH_EMAIL);
	localStorage.removeItem(CACHEKEYs.AUTH_STEP);
};
