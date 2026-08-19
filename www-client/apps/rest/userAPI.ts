import { notify } from '@/components/Notification/notify';
import { apiDelete, apiGet, apiPatch, apiPost } from '@/rest/apiClient';
import type { UserMeUpdateRequest } from '@/interface/user/userMeUpdateRequest.interface';
import { compressImageForUpload } from '@/utils/imageUploadUtils';

const apiPath = '/users';

export interface UserLoginStateResponse {
	id: number;
	created_at: string;
	email: string;
	email_confirmed: boolean;
	first_name: string | null;
	last_name: string | null;
	avatars: string | null;
	plan_name: string;
	valid_to: string | null;
	amount: number;
	currency: string;
	notification_subscriptions: boolean;
	max_total_subscriptions: number | null;
	auto_find_subscriptions: boolean;
	is_admin?: boolean;
}

export interface UserSessionResponse {
	session_id: string;
	platform: string;
	device_id?: string;
	created_at: string;
	current: boolean;
}

export interface UserSettingsResponse {
	display_currency?: string | null;
	country?: string | null;
	telegram_connected?: boolean;
	telegram_username?: string | null;
	max_connected?: boolean;
	max_username?: string | null;
}

export interface TelegramLinkResponse {
	bot_url: string;
}

export interface UserReferralRank {
	level: number;
	key: string;
	min_count: number;
	max_count?: number | null;
	reward_days: number;
	current: boolean;
}

export interface UserReferralResponse {
	code: string;
	site_url: string;
	bot_url: string;
	referral_count: number;
	rank: number;
	ranks: UserReferralRank[];
}

export interface UserPublicProfile {
	email: string;
	first_name?: string | null;
	last_name?: string | null;
	avatars?: string | null;
}

export interface AdminMessageRecipient {
	user_uuid: string;
	email: string;
	first_name?: string | null;
	last_name?: string | null;
	telegram_connected: boolean;
	max_connected: boolean;
}

export interface AdminSendMessageResponse {
	sent: number;
	skipped: number;
	failed: number;
}

export const basicAdminRecipients = async (): Promise<AdminMessageRecipient[]> => apiGet(`${apiPath}/admin/recipients`);

export const basicAdminSendMessage = async (payload: { channel: 'email' | 'telegram' | 'max'; user_uuid?: string | null; text: string }): Promise<AdminSendMessageResponse> => apiPost(`${apiPath}/admin/messages`, payload);

export const basicUserLoginState = async (): Promise<UserLoginStateResponse> => apiGet(`${apiPath}/login-state`);

export const basicUserSearchByEmail = async (email: string): Promise<UserPublicProfile[]> => {
	const users = await apiGet<UserPublicProfile[]>(`${apiPath}/search`, {
		params: { email },
	});

	return users ?? [];
};

export const basicUserSettingsGet = async (): Promise<UserSettingsResponse> => apiGet(`${apiPath}/settings`);

export const basicUserReferralGet = async (): Promise<UserReferralResponse> => apiGet(`${apiPath}/referral`);

export const basicUserSettingsUpdate = async (payload: { display_currency?: string; country?: string }): Promise<void> => {
	const message = await apiPatch<string>(`${apiPath}/settings`, payload);
	notify.success(message);
};

export const basicUserSessionsGetList = async (): Promise<UserSessionResponse[]> => apiGet(`${apiPath}/sessions`);

export const basicUserSessionsDisconnect = async (sessionId: string): Promise<void> => {
	const message = await apiPost<string>(`${apiPath}/sessions/disconnect`, { session_id: sessionId });
	notify.success(message);
};

export const basicUserProfileUpdate = async (payload: UserMeUpdateRequest): Promise<void> => {
	const message = await apiPatch<string>(`${apiPath}/me`, payload);
	notify.success(message);
};

export const basicUserAvatarUpdate = async (file: File): Promise<void> => {
	const payload = new FormData();
	payload.append('avatar', await compressImageForUpload(file));

	const message = await apiPost<string>(`${apiPath}/me/avatar`, payload);
	notify.success(message);
};

export const basicUserAccountDelete = async (): Promise<void> => {
	const message = await apiDelete<string>(`${apiPath}/me`);
	notify.success(message);
};

export const basicUserTelegramLink = async (): Promise<TelegramLinkResponse> => apiPost(`${apiPath}/telegram/link`);

export const basicUserTelegramDisconnect = async (): Promise<void> => {
	const message = await apiDelete<string>(`${apiPath}/telegram`);
	notify.success(message);
};

export const basicUserMaxLink = async (): Promise<TelegramLinkResponse> => apiPost(`${apiPath}/max/link`);

export const basicUserMaxDisconnect = async (): Promise<void> => {
	const message = await apiDelete<string>(`${apiPath}/max`);
	notify.success(message);
};
