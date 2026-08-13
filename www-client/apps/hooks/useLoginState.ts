import { QUERY_KEYS } from '@/constants/QueryKeys.constants';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { basicUserLoginState } from '@/rest/userAPI';
import { getInitialsFromName } from '@/utils/stringUtils';
import { getUserDisplayName, isPremiumPlan } from '@/utils/userDisplayUtils';

export const useLoginState = () => {
	const {
		data: loginState,
		loading,
		reload,
		updateHServer,
	} = useHandleServer([QUERY_KEYS.userLoginState], () => basicUserLoginState(), {
		staleTime: 0,
		refetchOnWindowFocus: true,
	});

	const displayName = loginState ? getUserDisplayName(loginState) : 'user';
	const initials = getInitialsFromName(displayName);
	const isPremium = isPremiumPlan(loginState?.plan_name);
	const canUseNotification = isPremium && (loginState?.notification_subscriptions ?? false);

	return {
		loginState,
		loading,
		reload,
		updateHServer,
		displayName,
		initials,
		isPremium,
		canUseNotification,
	};
};
