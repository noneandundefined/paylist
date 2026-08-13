import { CACHEKEYs } from '@/constants/CacheKeys.constants';

export const readAuthSession = (): boolean => Boolean(localStorage.getItem(CACHEKEYs.L_SESSION));
