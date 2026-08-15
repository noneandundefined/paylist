import type { TFunction } from 'i18next';

import { getDaysUntilRenewal, isSubscriptionOverdue } from './SubscriptionRenewalUtils';
import { getInitialsFromName } from '@/utils/stringUtils';
import { formatCurrency } from '@/utils/currencyUtils';

export const getSubscriptionIconLabel = (name: string): string => getInitialsFromName(name);

export const formatSubscriptionDate = (value: string, locale?: string): string => {
	const date = new Date(value);

	if (Number.isNaN(date.getTime())) {
		return value;
	}

	return date.toLocaleDateString(locale, {
		day: 'numeric',
		month: 'short',
		year: 'numeric',
	});
};

export const formatSubscriptionPrice = (price: number, currency: string, language?: string): string => {
	return formatCurrency(price, currency, language);
};

export const getSubscriptionShareAmount = (subscription: { price: number; share_price?: number }): number => {
	if (typeof subscription.share_price === 'number' && Number.isFinite(subscription.share_price)) {
		return subscription.share_price;
	}

	return subscription.price;
};

export const getRenewOnDateLabel = (datePay: string, t: TFunction, locale?: string): string => {
	if (isSubscriptionOverdue(datePay)) {
		return t('home.overdue-on-date', {
			date: formatSubscriptionDate(datePay, locale),
		});
	}

	return t('home.renew-on-date', {
		date: formatSubscriptionDate(datePay, locale),
	});
};

export const getNextBillingLabel = (datePay: string, t: TFunction, locale?: string): string => {
	const daysUntilRenewal = getDaysUntilRenewal(datePay);
	const formattedDate = formatSubscriptionDate(datePay, locale);

	if (daysUntilRenewal < 0) {
		return t('subscription.overdue-billing', {
			date: formattedDate,
			days: Math.abs(daysUntilRenewal),
		});
	}

	if (daysUntilRenewal === 1) {
		return t('subscription.next-billing-tomorrow', { date: formattedDate });
	}

	return t('subscription.next-billing-on', { date: formattedDate });
};
