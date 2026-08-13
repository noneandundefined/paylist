import axiosClient from './axios';

type AxiosConfig = Parameters<typeof axiosClient.get>[1];

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
