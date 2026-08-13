import axiosClient from './axios';
import { apiPost } from '@/rest/apiClient';

const apiPath = '/device';

export interface DeviceAuthSessionResponse {
	created_at: string;
	device_id: string;
	session_id: string;
	confirmed: boolean;
}

export interface DeviceAuthValidResponse {
	valid: boolean;
	message: string;
}

export const basicDeviceSessionGetBySessionId = async (sessionId: string): Promise<DeviceAuthValidResponse> => {
	const response = await axiosClient.get(`${apiPath}/session/${sessionId}`, {
		skipErrorHandler: true,
	} as Parameters<typeof axiosClient.get>[1] & { skipErrorHandler?: boolean });

	return response.data.message;
};

export const basicDeviceConfirm = async (session_id: string): Promise<DeviceAuthSessionResponse> => apiPost(`${apiPath}/confirm`, { session_id });
