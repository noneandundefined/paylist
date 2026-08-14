import axiosClient from './axios';
import { ROUTES } from '@/constants/constants';
import { notify } from '@/components/Notification/notify';
import { CACHEKEYs } from '@/constants/CacheKeys.constants';
import { clearAuthSession } from '@/utils/authSessionUtils';
import { apiGet, apiPost } from '@/rest/apiClient';
import type { AuthSigninRequest } from '@/interface/auth/authSigninRequest.interface';
import type { AuthSignupRequest } from '@/interface/auth/authSignupRequest.interface';

const apiPath = '/auth';

type AuthConfirmResponse = {
	status: 'success' | 'invalid' | 'error';
	message: string;
};

export type AuthStatusResponse = {
	status: 'signed_in' | 'sent' | 'password';
	message: string;
};

export const basicAuthSignIn = async (payload: AuthSigninRequest): Promise<AuthStatusResponse> => {
	const result = await apiPost<AuthStatusResponse>(`${apiPath}/signin`, payload);

	if (result.status === 'signed_in') {
		localStorage.setItem(CACHEKEYs.L_SESSION, result.message);
	}

	return result;
};

export const basicAuthSignUp = async (payload: AuthSignupRequest): Promise<string> => apiPost(`${apiPath}/signup`, payload);

export const basicAuthConfirmPending = async (email: string): Promise<{ pending: boolean }> => {
	const response = await axiosClient.get(`${apiPath}/confirm/pending`, {
		params: { email },
		skipErrorHandler: true,
	} as Parameters<typeof axiosClient.get>[1] & { skipErrorHandler?: boolean });

	return response.data.message;
};

export const basicAuthConfirmEmail = async (exp: unknown, sig: unknown, uuid: unknown): Promise<AuthConfirmResponse> => apiGet(`${apiPath}/confirm?exp=${exp}&sig=${sig}&uuid=${uuid}`);

export const basicAuthSignOut = async (): Promise<void> => {
	try {
		const message = await apiPost<string>(`${apiPath}/signout`);
		notify.success(message);
	} finally {
		clearAuthSession();
		window.location.replace(ROUTES.SIGNIN);
	}
};
