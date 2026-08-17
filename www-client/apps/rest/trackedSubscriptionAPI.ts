import { notify } from '@/components/Notification/notify';
import { apiDelete, apiDownload, apiGet, apiPost, apiPut } from '@/rest/apiClient';
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
	tariff?: string;
	price: number;
	currency: string;
	period: TrackedSubscriptionPeriod;
	date_pay: string;
	auto_renewal: boolean;
	notification: boolean;
	include_in_analytics: boolean;
	note?: string | null;
	categories?: string[];
	share_percent?: number;
	share_price?: number;
	is_owner?: boolean;
}

export type SharedSubscriptionMemberStatus = 'pending' | 'accepted' | 'declined';
export type SharedSubscriptionMemberRole = 'owner' | 'member' | 'observer';

export interface SharedSubscriptionMember {
	id: number;
	created_at: string;
	updated_at: string;
	tracked_subscription_id: number;
	user_uuid?: string | null;
	email: string;
	role: SharedSubscriptionMemberRole;
	share_percent: number;
	notification: boolean;
	include_in_analytics: boolean;
	status: SharedSubscriptionMemberStatus;
	invite_expires_at?: string | null;
	first_name?: string | null;
	last_name?: string | null;
	avatars?: string | null;
}

export interface SharedSubscriptionShareItem {
	proposal_id: number;
	member_id: number;
	share_percent: number;
}

export interface SharedSubscriptionShareVote {
	proposal_id: number;
	user_uuid: string;
	accepted: boolean;
	created_at: string;
}

export interface SharedSubscriptionShareProposal {
	id: number;
	proposed_by_user_uuid: string;
	status: string;
	items: SharedSubscriptionShareItem[];
	votes: SharedSubscriptionShareVote[];
	my_vote: boolean | null;
}

export interface SharedSubscriptionMembersResponse {
	members: SharedSubscriptionMember[];
	pending_proposal: SharedSubscriptionShareProposal | null;
}

export interface SharedSubscriptionInvitePreview {
	subscription_id: number;
	subscription_name: string;
	owner_name: string;
	email: string;
	share_percent: number;
	role?: SharedSubscriptionMemberRole;
	status: SharedSubscriptionMemberStatus;
	invite_expires_at?: string | null;
}

export interface SharedSubscriptionInviteAcceptResponse {
	message: string;
	subscription_id: number;
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

export interface ServiceResponse {
	id: number;
	slug: string;
	name: string;
	category: string;
	aliases: string[];
}

const servicesByName = new Map<string, ServiceResponse>();

export const fetchServices = async (): Promise<ServiceResponse[]> => {
	const items = await apiGet<ServiceResponse[]>(`${apiPath}/services`);
	return items ?? [];
};

export const loadServiceSelectOptions = async () => {
	const items = await fetchServices();

	servicesByName.clear();

	return items.map((item) => {
		servicesByName.set(item.name, item);

		return {
			value: item.name,
			label: item.name,
			keywords: `${item.name} ${item.slug} ${item.aliases.join(' ')}`,
		};
	});
};

export const getServiceByName = (name: string): ServiceResponse | undefined => servicesByName.get(name);

export const basicTrackedSubscriptionSummary = async (): Promise<TrackedSubscriptionSummaryResponse> => apiGet(`${apiPath}/summary`);

export type AnalyticsRecommendationType = 'yearly-save' | 'concentration' | 'excluded' | 'cluster' | 'small-subs' | 'upcoming-heavy' | 'overlap' | 'family-share' | 'crowd-overpay' | 'downgrade' | 'expensive-tariff';

export interface AnalyticsRecommendation {
	id: string;
	type: AnalyticsRecommendationType;
	title_key: string;
	desc_key: string;
	desc_values?: Record<string, string | number>;
	subscription_id?: number;
}

export interface AnalyticsRecommendationsResponse {
	recommendations: AnalyticsRecommendation[];
}

export const basicTrackedSubscriptionAnalytics = async (): Promise<AnalyticsRecommendationsResponse> => {
	const payload = await apiGet<AnalyticsRecommendationsResponse>(`${apiPath}/analytics`);
	return payload ?? { recommendations: [] };
};

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

export const basicTrackedSubscriptionExport = async (): Promise<void> => {
	await apiDownload(`${apiPath}/export`, 'paylist-subscriptions.csv');
};

export const basicTrackedSubscriptionGetById = async (id: number): Promise<TrackedSubscriptionDetailResponse> => apiGet(`${apiPath}/${id}`);

export const basicTrackedSubscriptionCreate = async (payload: TrackedSubscriptionCreateRequest): Promise<string> => apiPost(apiPath, payload);

export const basicTrackedSubscriptionUpdate = async (id: number, payload: TrackedSubscriptionEditRequest): Promise<string> => apiPut(`${apiPath}/${id}`, payload);

export const basicTrackedSubscriptionDelete = async (id: number): Promise<string> => apiDelete(`${apiPath}/${id}`);

export const basicTrackedSubscriptionMembers = async (id: number): Promise<SharedSubscriptionMembersResponse> => apiGet(`${apiPath}/${id}/members`);

export const basicTrackedSubscriptionInvite = async (id: number, email: string, sharePercent: number, role: SharedSubscriptionMemberRole = 'member'): Promise<string> =>
	apiPost(`${apiPath}/${id}/members`, { email, share_percent: sharePercent, role });

export const basicTrackedSubscriptionRemoveMember = async (id: number, memberId: number): Promise<string> => apiDelete(`${apiPath}/${id}/members/${memberId}`);

export const basicTrackedSubscriptionLeave = async (id: number): Promise<string> => apiDelete(`${apiPath}/${id}/members/me`);

export const basicTrackedSubscriptionProposeShares = async (id: number, shares: Array<{ member_id: number; share_percent: number }>): Promise<string> => apiPost(`${apiPath}/${id}/shares`, { shares });

export const basicTrackedSubscriptionVoteShares = async (id: number, proposalId: number, accept: boolean): Promise<string> => apiPost(`${apiPath}/${id}/shares/${proposalId}/vote`, { accept });

export const basicTrackedSubscriptionInvitePreview = async (token: string): Promise<SharedSubscriptionInvitePreview> => apiGet(`${apiPath}/invites`, { params: { token } });

export const basicTrackedSubscriptionAcceptInvite = async (token: string): Promise<SharedSubscriptionInviteAcceptResponse> => apiPost(`${apiPath}/invites/accept`, { token });

export const getTrackedSubscriptionImageUrl = (name: string): string => {
	const baseURL = configClient.type.release === 'dev' ? configClient.links.URL_BACKEND_DEV : configClient.links.URL_BACKEND_PROD;

	return `${baseURL}${apiPath}/images/w350?name=${encodeURIComponent(name)}`;
};
