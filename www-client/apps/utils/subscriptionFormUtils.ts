import type { TFunction } from 'i18next';

export const validateFutureDate = (value: string, t: TFunction): true | string => {
	const selected = new Date(value);
	const today = new Date();
	today.setHours(0, 0, 0, 0);
	selected.setHours(0, 0, 0, 0);

	return selected > today || t('subscription.date-invalid');
};
