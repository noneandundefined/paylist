import { CACHEKEYs } from '@/constants/CacheKeys.constants';

const REF_QUERY = 'ref';

export const captureReferralCodeFromSearch = (search: string) => {
	const params = new URLSearchParams(search.startsWith('?') ? search.slice(1) : search);
	const code = params.get(REF_QUERY)?.trim();

	if (code) {
		sessionStorage.setItem(CACHEKEYs.REFERRAL_CODE, code);
	}
};

export const getStoredReferralCode = (): string => sessionStorage.getItem(CACHEKEYs.REFERRAL_CODE)?.trim() ?? '';

export const clearStoredReferralCode = () => {
	sessionStorage.removeItem(CACHEKEYs.REFERRAL_CODE);
};
