import { CACHEKEYs } from './constants/CacheKeys.constants';

export async function getAuthState() {
	return Boolean(localStorage.getItem(CACHEKEYs.L_SESSION));
}
