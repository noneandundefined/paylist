import { apiGet } from '@/rest/apiClient';

export interface PlanResponse {
	id: number;
	created_at: string;
	plan_name: string;
	amount: number;
	currency: string;
	duration_days: number;
	max_total_subscriptions?: number | null;
	notification_subscriptions: boolean;
	auto_find_subscriptions: boolean;
	description: Record<string, string>;
	features: Record<string, string[]>;
}

export const basicPlanList = async (): Promise<PlanResponse[]> => {
	const plans = await apiGet<PlanResponse[]>('/plans');
	return plans ?? [];
};
