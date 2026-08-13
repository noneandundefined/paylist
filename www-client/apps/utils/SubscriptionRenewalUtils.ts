export const RENEWAL_BADGE_MAX_DAYS = 3;

export const getDaysUntilRenewal = (nextBillingAt: string): number => {
	const today = new Date();
	today.setHours(0, 0, 0, 0);

	const billingDate = new Date(nextBillingAt);
	billingDate.setHours(0, 0, 0, 0);

	return Math.round((billingDate.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
};

export const isSubscriptionOverdue = (datePay: string): boolean => getDaysUntilRenewal(datePay) < 0;

export const getOverdueDays = (datePay: string): number => Math.abs(getDaysUntilRenewal(datePay));

export const getRenewalBadgeDays = (datePay: string): number | null => {
	const days = getDaysUntilRenewal(datePay);

	if (days < 0 || days > RENEWAL_BADGE_MAX_DAYS) {
		return null;
	}

	return days;
};
