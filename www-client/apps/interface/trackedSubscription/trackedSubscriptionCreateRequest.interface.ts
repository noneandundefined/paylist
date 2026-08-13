export type TrackedSubscriptionPeriod = 'monthly' | 'yearly';

export interface TrackedSubscriptionCreateRequest {
	name: string;
	price: number;
	currency?: string;
	period?: TrackedSubscriptionPeriod;
	date_pay: string;
	auto_renewal?: boolean;
	notification?: boolean;
	include_in_analytics?: boolean;
	categories?: string[];
}
