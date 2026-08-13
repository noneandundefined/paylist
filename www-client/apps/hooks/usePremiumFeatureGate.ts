import { useTranslation } from 'react-i18next';
import { useLoginState } from '@/hooks/useLoginState';
import { requirePremiumFeature } from '@/utils/premiumUtils';

export const usePremiumFeatureGate = () => {
	const { t } = useTranslation();
	const { loginState, loading, isPremium, canUseNotification, displayName, initials } = useLoginState();

	const requirePremium = (allowed: boolean) => requirePremiumFeature(allowed, t);

	return {
		loginState,
		loading,
		isPremium,
		canUseNotification,
		displayName,
		initials,
		requirePremium,
	};
};
