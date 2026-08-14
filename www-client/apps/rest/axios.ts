import Qs from 'qs';
import i18next from 'i18next';
import { PoWDDosDecision } from '@/utils/PowDDosUtils';
import { notify } from '@/components/Notification/notify';
import { CACHEKEYs } from '@/constants/CacheKeys.constants';
import { clearAuthSession } from '@/utils/authSessionUtils';
import axios, { AxiosError, AxiosRequestConfig } from 'axios';
import { config as configClient } from '@/.config/config.client';

interface AxiosRequestConfigWithRetry extends AxiosRequestConfig {
	_retry?: boolean;
	skipErrorHandler?: boolean;
	actionKey?: string;
	skipIdempotency?: boolean;
}

/** Axios base url */
const axiosClient = axios.create({
	baseURL: configClient.type.release == 'dev' ? configClient.links.URL_BACKEND_DEV : configClient.links.URL_BACKEND_PROD,
	paramsSerializer: (params) => Qs.stringify(params, { arrayFormat: 'comma' }),
});

/** Axios helper for request */
axiosClient.interceptors.request.use(
	(config) => {
		/** Auth token */
		const token = localStorage.getItem(CACHEKEYs.L_SESSION);
		if (token) {
			config.headers = config.headers || {};
			config.headers['Authorization'] = `Bearer ${token}`;
		}

		/** Language header */
		const currentLanguage = i18next.language || 'en';
		config.headers['Accept-Language'] = currentLanguage;

		const xDebug = configClient.type.release == 'dev' ? '0' : '0';
		config.headers['X-Debug'] = xDebug;

		/** RequestId header */
		const requestID = localStorage.getItem(CACHEKEYs.PAYLIST_X_REQ_ID);
		if (requestID) {
			config.headers['X-Request-ID'] = requestID;
		}

		return config;
	},
	(error) => Promise.reject(error)
);

/** Axios helper for response */
axiosClient.interceptors.response.use(
	(response) => {
		const requestID = response.headers['x-request-id'];
		if (requestID) {
			localStorage.setItem(CACHEKEYs.PAYLIST_X_REQ_ID, requestID);
		}

		return response;
	},
	async (error) => {
		if (axios.isCancel(error) || error.code === 'ERR_CANCELED') {
			return Promise.reject(error);
		}

		// internet
		if (error.message === 'Network Error' || error.code === 'ERR_NETWORK' || error.message.includes('Network request failed')) {
			notify.error(i18next.t('message.internet-error'));
			return Promise.reject(error);
		}

		if (error instanceof AxiosError) {
			/** Pow DDoS */
			if (error?.response?.status == 429 && error.response.data?.error?.challenge) {
				const { challenge, difficulty } = error.response.data.error;

				const nonce = await PoWDDosDecision(challenge, difficulty);

				const config = {
					...error.config,
					headers: {
						...error.config?.headers,
						'Pow-Challenge': challenge,
						'Pow-Nonce': nonce,
					},
				};

				return axiosClient(config);
			}

			const originalRequest = error.config as AxiosRequestConfigWithRetry;

			if (originalRequest.skipErrorHandler) {
				return Promise.reject(error);
			}

			if (error.response && error.response.status) {
				if (error.response.status === 401) {
					try {
						clearAuthSession();
					} catch (e) {
						console.error(e);
					}
				} else {
					notify.error(error.response.data?.error || i18next.t('message.server-error'));
				}
			} else if (error.request) {
				notify.error(i18next.t('message.request-error'));
			} else {
				notify.error(i18next.t('message.unknown-request-error'));
			}
		} else {
			notify.error(i18next.t('message.try-again-error'));
		}

		return Promise.reject(error);
	}
);

export default axiosClient;
