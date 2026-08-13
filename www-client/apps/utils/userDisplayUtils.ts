import type { UserLoginStateResponse } from '@/rest/userAPI';
import { getInitialsFromName } from '@/utils/stringUtils';

export const isPremiumPlan = (planName?: string | null): boolean => planName?.toLowerCase() === 'premium';

export const getUserDisplayName = (user: Pick<UserLoginStateResponse, 'first_name' | 'last_name' | 'email'>): string => {
	const parts = [user.first_name?.trim(), user.last_name?.trim()].filter(Boolean);

	if (parts.length > 0) {
		return parts.join(' ');
	}

	if (user.email) {
		return user.email.split('@')[0] ?? user.email;
	}

	return 'user';
};

export const getUserInitials = (displayName: string): string => getInitialsFromName(displayName);

export const getUserHandle = (email: string): string => {
	const localPart = email.split('@')[0]?.trim() || 'user';

	return `@${localPart}`;
};
