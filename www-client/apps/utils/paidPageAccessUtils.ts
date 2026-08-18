import { CACHEKEYs } from '@/constants/CacheKeys.constants';

export const readPaidPageGrant = (): string | null => {
	try {
		return sessionStorage.getItem(CACHEKEYs.PAID_PAGE_GRANT);
	} catch {
		return null;
	}
};

export const writePaidPageGrant = (paymentId: string): void => {
	try {
		sessionStorage.setItem(CACHEKEYs.PAID_PAGE_GRANT, paymentId);
	} catch {
		return;
	}
};

export const hasPaidPageGrant = (paymentId?: string | null): boolean => {
	const grant = readPaidPageGrant();

	if (!grant) {
		return false;
	}

	if (!paymentId) {
		return true;
	}

	return grant === paymentId;
};
