import { notify } from '@/components/Notification/notify';
import { apiDelete, apiGet, apiPatch, apiPost } from '@/rest/apiClient';
import type { UserMeUpdateRequest } from '@/interface/user/userMeUpdateRequest.interface';

const apiPath = '/users';

export interface UserLoginStateResponse {
	id: number;
	created_at: string;
	email: string;
	email_confirmed: boolean;
	first_name: string | null;
	last_name: string | null;
	plan_name: string;
	valid_to: string | null;
	amount: number;
	currency: string;
	notification_subscriptions: boolean;
	max_total_subscriptions: number | null;
	auto_find_subscriptions: boolean;
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
}

export interface TelegramLinkResponse {
	bot_url: string;
}

export const basicUserLoginState = async (): Promise<UserLoginStateResponse> => apiGet(`${apiPath}/login-state`);

export const basicUserSettingsGet = async (): Promise<UserSettingsResponse> => apiGet(`${apiPath}/settings`);

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

export const basicUserAccountDelete = async (): Promise<void> => {
	const message = await apiDelete<string>(`${apiPath}/me`);
	notify.success(message);
};

export const basicUserTelegramLink = async (): Promise<TelegramLinkResponse> => apiPost(`${apiPath}/telegram/link`);

export const basicUserTelegramDisconnect = async (): Promise<void> => {
	const message = await apiDelete<string>(`${apiPath}/telegram`);
	notify.success(message);
};
