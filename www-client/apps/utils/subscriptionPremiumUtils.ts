export interface PremiumSubscriptionFlagsInput {
	autoRenewal: boolean;
	notification: boolean;
	canUseNotification: boolean;
}

export const resolvePremiumSubscriptionFlags = ({ autoRenewal, notification, canUseNotification }: PremiumSubscriptionFlagsInput) => ({
	auto_renewal: autoRenewal,
	notification: canUseNotification ? notification : false,
});

export const clampPremiumSubscriptionFlags = (autoRenewal: boolean, notification: boolean, canUseNotification: boolean) => ({
	autoRenewal,
	notification: canUseNotification ? notification : false,
});
