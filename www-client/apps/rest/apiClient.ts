import { AxiosError } from 'axios';
import i18next from 'i18next';
import axiosClient from './axios';
import { notify } from '@/components/Notification/notify';
import { clearAuthSession } from '@/utils/authSessionUtils';

type AxiosConfig = Parameters<typeof axiosClient.get>[1] & { skipErrorHandler?: boolean };

export const apiGet = async <T>(url: string, config?: AxiosConfig): Promise<T> => {
	const response = await axiosClient.get(url, config);
	return response.data.message as T;
};

export const apiPost = async <T>(url: string, body?: unknown, config?: AxiosConfig): Promise<T> => {
	const response = await axiosClient.post(url, body, config);
	return response.data.message as T;
};

export const apiPut = async <T>(url: string, body?: unknown, config?: AxiosConfig): Promise<T> => {
	const response = await axiosClient.put(url, body, config);
	return response.data.message as T;
};

export const apiPatch = async <T>(url: string, body?: unknown, config?: AxiosConfig): Promise<T> => {
	const response = await axiosClient.patch(url, body, config);
	return response.data.message as T;
};

export const apiDelete = async <T>(url: string, config?: AxiosConfig): Promise<T> => {
	const response = await axiosClient.delete(url, config);
	return response.data.message as T;
};

const headerString = (value: unknown): string | undefined => {
	if (typeof value === 'string' && value !== '') {
		return value;
	}

	if (Array.isArray(value) && typeof value[0] === 'string' && value[0] !== '') {
		return value[0];
	}

	return undefined;
};

const filenameFromDisposition = (disposition?: string) => {
	const utfMatch = disposition?.match(/filename\*=UTF-8''([^;]+)/i);
	if (utfMatch?.[1]) {
		return decodeURIComponent(utfMatch[1]);
	}

	const match = disposition?.match(/filename="?([^"]+)"?/i);
	return match?.[1];
};

const notifyBlobError = async (error: unknown) => {
	if (!(error instanceof AxiosError) || !(error.response?.data instanceof Blob)) {
		return;
	}

	if (error.response.status === 401) {
		try {
			clearAuthSession();
		} catch (e) {
			console.error(e);
		}
		return;
	}

	try {
		const parsed = JSON.parse(await error.response.data.text()) as { error?: string };
		notify.error(parsed.error || i18next.t('message.server-error'));
	} catch {
		notify.error(i18next.t('message.server-error'));
	}
};

export const apiDownload = async (url: string, fallbackName: string): Promise<void> => {
	try {
		const response = await axiosClient.get(url, { responseType: 'blob', skipErrorHandler: true } as AxiosConfig);
		const blob = new Blob([response.data], { type: headerString(response.headers['content-type']) || 'application/octet-stream' });
		const objectUrl = window.URL.createObjectURL(blob);
		const link = document.createElement('a');
		link.href = objectUrl;
		link.download = filenameFromDisposition(headerString(response.headers['content-disposition'])) || fallbackName;
		document.body.appendChild(link);
		link.click();
		link.remove();
		window.URL.revokeObjectURL(objectUrl);
	} catch (error) {
		await notifyBlobError(error);
		throw error;
	}
};
