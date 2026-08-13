import { notify } from '@/components/Notification/notify';
import { apiDelete, apiGet, apiPost, apiPut } from '@/rest/apiClient';
import { config as configClient } from '@/.config/config.client';
import type { TrackedSubscriptionPeriod } from '@/interface/trackedSubscription/trackedSubscriptionCreateRequest.interface';
import type { TrackedSubscriptionEditRequest } from '@/interface/trackedSubscription/trackedSubscriptionEditRequest.interface';
import type { TrackedSubscriptionCreateRequest } from '@/interface/trackedSubscription/trackedSubscriptionCreateRequest.interface';

const apiPath = '/tracked-subscription';

export type { TrackedSubscriptionPeriod };

export interface TrackedSubscriptionSummaryResponse {
	display_currency: string;
	total_amount: number;
	active_count: number;
	preview_subscription_ids: number[];
	comparison_percent: number;
	comparison_direction: 'less' | 'more';
}

export interface TrackedSubscriptionResponse {
	id: number;
	created_at: string;
	updated_at: string;
	user_uuid: string;
	name: string;
	price: number;
	currency: string;
	period: TrackedSubscriptionPeriod;
	date_pay: string;
	auto_renewal: boolean;
	notification: boolean;
	include_in_analytics: boolean;
	note?: string | null;
	categories?: string[];
}

export interface TrackedSubscriptionDetailResponse extends TrackedSubscriptionResponse {
	categories: string[];
}

export interface SubscriptionCategoryResponse {
	id: number;
	created_at: string;
	slug: string;
	label?: string | null;
	is_custom?: boolean;
}

export const basicTrackedSubscriptionSummary = async (): Promise<TrackedSubscriptionSummaryResponse> => apiGet(`${apiPath}/summary`);

export const basicSubscriptionCategoryList = async (): Promise<SubscriptionCategoryResponse[]> => {
	const categories = await apiGet<SubscriptionCategoryResponse[]>(`${apiPath}/categories`);
	return categories ?? [];
};

export const basicSubscriptionCategoryCreate = async (label: string): Promise<SubscriptionCategoryResponse> => apiPost(`${apiPath}/categories`, { label });

export const basicSubscriptionCategoryDelete = async (categoryId: number): Promise<void> => {
	const message = await apiDelete<string>(`${apiPath}/categories/${categoryId}`);
	notify.success(message);
};

export const basicTrackedSubscriptionList = async (search?: string): Promise<TrackedSubscriptionResponse[]> => {
	const subscriptions = await apiGet<TrackedSubscriptionResponse[]>(apiPath, {
		params: search ? { search } : undefined,
	});

	return subscriptions ?? [];
};

export const basicTrackedSubscriptionGetById = async (id: number): Promise<TrackedSubscriptionDetailResponse> => apiGet(`${apiPath}/${id}`);

export const basicTrackedSubscriptionCreate = async (payload: TrackedSubscriptionCreateRequest): Promise<string> => apiPost(apiPath, payload);

export const basicTrackedSubscriptionUpdate = async (id: number, payload: TrackedSubscriptionEditRequest): Promise<string> => apiPut(`${apiPath}/${id}`, payload);

export const basicTrackedSubscriptionDelete = async (id: number): Promise<string> => apiDelete(`${apiPath}/${id}`);

export const getTrackedSubscriptionImageUrl = (name: string): string => {
	const baseURL = configClient.type.release === 'dev' ? configClient.links.URL_BACKEND_DEV : configClient.links.URL_BACKEND_PROD;

	return `${baseURL}${apiPath}/images/w350?name=${encodeURIComponent(name)}`;
};
