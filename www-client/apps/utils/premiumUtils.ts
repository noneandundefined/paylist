import type { TFunction } from 'i18next';
import { notify } from '@/components/Notification/notify';

export const notifyPremiumRequired = (t: TFunction): void => {
	notify.error(t('subscription.premium-required'));
};

export const requirePremiumFeature = (allowed: boolean, t: TFunction): boolean => {
	if (!allowed) {
		notifyPremiumRequired(t);
		return false;
	}

	return true;
};

export const notifySoon = (t: TFunction, messageKey: string): void => {
	notify.success(t(messageKey));
};
