import type { TrackedSubscriptionResponse } from '@/rest/trackedSubscriptionAPI';

export type AnalyticsPeriod = '1M' | '3M' | '6M' | '1Y';

export interface AnalyticsRecommendation {
	id: string;
	type: 'yearly-save' | 'concentration' | 'excluded' | 'cluster' | 'small-subs' | 'upcoming-heavy';
	titleKey: string;
	descKey: string;
	descValues?: Record<string, string | number>;
	subscriptionId?: number;
}

export interface MonthlyProjectionPoint {
	label: string;
	monthKey: string;
	amount: number;
}

export interface SubscriptionShare {
	subscription: TrackedSubscriptionResponse;
	monthlyAmount: number;
	sharePercent: number;
}

export interface CategoryShare {
	slug: string;
	monthlyAmount: number;
	sharePercent: number;
	color: string;
}

export const CATEGORY_CHART_COLORS = ['#F97316', '#EAB308', '#14B8A6', '#0085FF', '#8B5CF6', '#64748B'] as const;

export interface AnalyticsSnapshot {
	displayCurrency: string;
	subscriptionCount: number;
	analyticsCount: number;
	monthlyTotal: number;
	yearlyTotal: number;
	dailyAverage: number;
	weeklyAverage: number;
	next30DaysOutflow: number;
	next90DaysOutflow: number;
	monthlyCount: number;
	yearlyCount: number;
	excludedCount: number;
	topSubscriptions: SubscriptionShare[];
	categoryBreakdown: CategoryShare[];
	projections: MonthlyProjectionPoint[];
	recommendations: AnalyticsRecommendation[];
	upcomingPayments: Array<{ subscription: TrackedSubscriptionResponse; amount: number; daysUntil: number }>;
}

export const getMonthlyAmount = (price: number, period: string): number => {
	if (period === 'yearly') {
		return price / 12;
	}

	return price;
};

export const getYearlyAmount = (price: number, period: string): number => {
	if (period === 'yearly') {
		return price;
	}

	return price * 12;
};

export const convertWithRates = (amount: number, from: string, to: string, base: string, rates: Record<string, number>): number => {
	const source = from.toUpperCase();
	const target = to.toUpperCase();
	const ratesBase = base.toUpperCase();

	if (source === target) {
		return amount;
	}

	if (ratesBase === target && rates[source]) {
		return amount / rates[source];
	}

	if (ratesBase === source && rates[target]) {
		return amount * rates[target];
	}

	return amount;
};

const monthKey = (date: Date) => `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`;

const addMonths = (date: Date, count: number) => {
	const next = new Date(date);
	next.setMonth(next.getMonth() + count);
	return next;
};

const normalizeDate = (value: string) => {
	const date = new Date(value.slice(0, 10));
	date.setHours(0, 0, 0, 0);
	return date;
};

export const getNextPaymentDate = (datePay: string, period: string): Date => {
	const anchor = normalizeDate(datePay);
	const today = new Date();
	today.setHours(0, 0, 0, 0);

	const next = new Date(anchor);

	while (next.getTime() < today.getTime()) {
		if (period === 'yearly') {
			next.setFullYear(next.getFullYear() + 1);
		} else {
			next.setMonth(next.getMonth() + 1);
		}
	}

	return next;
};

const paymentAmountInMonth = (subscription: TrackedSubscriptionResponse, target: Date): number => {
	const anchor = normalizeDate(subscription.date_pay);

	if (subscription.period === 'monthly') {
		return subscription.price;
	}

	if (target.getMonth() === anchor.getMonth()) {
		return subscription.price;
	}

	return 0;
};

const getDaysUntil = (datePay: string, period: string): number => {
	const today = new Date();
	today.setHours(0, 0, 0, 0);

	const target = getNextPaymentDate(datePay, period);
	const diffMs = target.getTime() - today.getTime();

	return Math.round(diffMs / (1000 * 60 * 60 * 24));
};

const outflowWithinDays = (subscriptions: TrackedSubscriptionResponse[], days: number, convert: (amount: number, from: string) => number): number => {
	const today = new Date();
	today.setHours(0, 0, 0, 0);
	const horizon = new Date(today);
	horizon.setDate(horizon.getDate() + days);

	let total = 0;

	for (const subscription of subscriptions) {
		if (!subscription.include_in_analytics) {
			continue;
		}

		let next = getNextPaymentDate(subscription.date_pay, subscription.period);

		while (next.getTime() <= horizon.getTime()) {
			if (next.getTime() >= today.getTime()) {
				total += convert(subscription.price, subscription.currency);
			}

			if (subscription.period === 'yearly') {
				next = new Date(next);
				next.setFullYear(next.getFullYear() + 1);
			} else {
				next = new Date(next);
				next.setMonth(next.getMonth() + 1);
			}
		}
	}

	return total;
};

const buildProjections = (subscriptions: TrackedSubscriptionResponse[], months: number, displayCurrency: string, rates: Record<string, number>, ratesBase: string, locale: string): MonthlyProjectionPoint[] => {
	const today = new Date();
	today.setDate(1);
	today.setHours(0, 0, 0, 0);

	const points: MonthlyProjectionPoint[] = [];

	for (let index = 0; index < months; index += 1) {
		const target = addMonths(today, index);
		let amount = 0;

		for (const subscription of subscriptions) {
			if (!subscription.include_in_analytics) {
				continue;
			}

			const nativeAmount = paymentAmountInMonth(subscription, target);
			if (nativeAmount <= 0) {
				continue;
			}

			amount += convertWithRates(nativeAmount, subscription.currency, displayCurrency, ratesBase, rates);
		}

		points.push({
			monthKey: monthKey(target),
			label: target.toLocaleDateString(locale, { month: 'short' }),
			amount: Math.round(amount * 100) / 100,
		});
	}

	return points;
};

const buildRecommendations = (subscriptions: TrackedSubscriptionResponse[], monthlyTotal: number, displayCurrency: string, convert: (amount: number, from: string) => number): AnalyticsRecommendation[] => {
	const recommendations: AnalyticsRecommendation[] = [];

	for (const subscription of subscriptions) {
		if (subscription.period !== 'monthly' || !subscription.include_in_analytics) {
			continue;
		}

		const yearlyIfMonthly = subscription.price * 12;
		const estimatedYearlyPlan = subscription.price * 10;
		const savings = yearlyIfMonthly - estimatedYearlyPlan;

		if (savings >= subscription.price * 1.5) {
			recommendations.push({
				id: `yearly-${subscription.id}`,
				type: 'yearly-save',
				titleKey: 'analytics.rec-yearly-title',
				descKey: 'analytics.rec-yearly-desc',
				descValues: {
					name: subscription.name,
					percent: Math.round((savings / yearlyIfMonthly) * 100),
					amount: Math.round(convert(savings, subscription.currency)),
					currency: displayCurrency,
				},
				subscriptionId: subscription.id,
			});
		}
	}

	if (monthlyTotal > 0) {
		const sorted = [...subscriptions]
			.filter((item) => item.include_in_analytics)
			.map((item) => ({
				item,
				monthly: convert(getMonthlyAmount(item.price, item.period), item.currency),
			}))
			.sort((left, right) => right.monthly - left.monthly);

		const top = sorted[0];
		if (top && top.monthly / monthlyTotal >= 0.4) {
			recommendations.push({
				id: `concentration-${top.item.id}`,
				type: 'concentration',
				titleKey: 'analytics.rec-concentration-title',
				descKey: 'analytics.rec-concentration-desc',
				descValues: {
					name: top.item.name,
					percent: Math.round((top.monthly / monthlyTotal) * 100),
				},
				subscriptionId: top.item.id,
			});
		}
	}

	const excluded = subscriptions.filter((item) => !item.include_in_analytics);
	if (excluded.length > 0) {
		recommendations.push({
			id: 'excluded',
			type: 'excluded',
			titleKey: 'analytics.rec-excluded-title',
			descKey: 'analytics.rec-excluded-desc',
			descValues: { count: excluded.length },
		});
	}

	const dueSoon = subscriptions.filter((item) => {
		const days = getDaysUntil(item.date_pay, item.period);
		return days >= 0 && days <= 7;
	});

	if (dueSoon.length >= 3) {
		const total = dueSoon.reduce((sum, item) => sum + convert(item.period === 'yearly' ? item.price : item.price, item.currency), 0);
		recommendations.push({
			id: 'cluster',
			type: 'cluster',
			titleKey: 'analytics.rec-cluster-title',
			descKey: 'analytics.rec-cluster-desc',
			descValues: {
				count: dueSoon.length,
				amount: Math.round(total * 100) / 100,
				currency: displayCurrency,
			},
		});
	}

	const smallSubs = subscriptions.filter((item) => {
		if (!item.include_in_analytics) {
			return false;
		}

		const monthly = convert(getMonthlyAmount(item.price, item.period), item.currency);
		return monthly > 0 && monthly / monthlyTotal < 0.05;
	});

	if (smallSubs.length >= 4 && monthlyTotal > 0) {
		const total = smallSubs.reduce((sum, item) => sum + convert(getMonthlyAmount(item.price, item.period), item.currency), 0);
		recommendations.push({
			id: 'small-subs',
			type: 'small-subs',
			titleKey: 'analytics.rec-small-title',
			descKey: 'analytics.rec-small-desc',
			descValues: {
				count: smallSubs.length,
				amount: Math.round(total * 100) / 100,
				currency: displayCurrency,
			},
		});
	}

	const next30 = outflowWithinDays(
		subscriptions.filter((item) => item.include_in_analytics),
		30,
		convert
	);

	if (next30 > monthlyTotal * 1.35 && monthlyTotal > 0) {
		recommendations.push({
			id: 'upcoming-heavy',
			type: 'upcoming-heavy',
			titleKey: 'analytics.rec-upcoming-title',
			descKey: 'analytics.rec-upcoming-desc',
			descValues: {
				amount: Math.round(next30 * 100) / 100,
				currency: displayCurrency,
			},
		});
	}

	return recommendations.slice(0, 6);
};

const buildCategoryBreakdown = (subscriptions: TrackedSubscriptionResponse[], monthlyTotal: number, convert: (amount: number, from: string) => number): CategoryShare[] => {
	const buckets = new Map<string, number>();

	for (const subscription of subscriptions) {
		if (!subscription.include_in_analytics) {
			continue;
		}

		const monthlyAmount = convert(getMonthlyAmount(subscription.price, subscription.period), subscription.currency);
		const categories = subscription.categories?.length ? subscription.categories : ['uncategorized'];
		const share = monthlyAmount / categories.length;

		for (const slug of categories) {
			buckets.set(slug, (buckets.get(slug) ?? 0) + share);
		}
	}

	const sorted = [...buckets.entries()]
		.map(([slug, amount]) => ({
			slug,
			monthlyAmount: Math.round(amount * 100) / 100,
		}))
		.sort((left, right) => right.monthlyAmount - left.monthlyAmount);

	if (sorted.length === 0) {
		return [];
	}

	const topItems = sorted.slice(0, 4);
	const otherAmount = sorted.slice(4).reduce((sum, item) => sum + item.monthlyAmount, 0);

	const items: CategoryShare[] = topItems.map((item, index) => ({
		slug: item.slug,
		monthlyAmount: item.monthlyAmount,
		sharePercent: monthlyTotal > 0 ? Math.round((item.monthlyAmount / monthlyTotal) * 1000) / 10 : 0,
		color: CATEGORY_CHART_COLORS[index % CATEGORY_CHART_COLORS.length],
	}));

	if (otherAmount > 0) {
		items.push({
			slug: 'other',
			monthlyAmount: Math.round(otherAmount * 100) / 100,
			sharePercent: monthlyTotal > 0 ? Math.round((otherAmount / monthlyTotal) * 1000) / 10 : 0,
			color: CATEGORY_CHART_COLORS[5],
		});
	}

	return items;
};

export const buildAnalyticsSnapshot = (subscriptions: TrackedSubscriptionResponse[], displayCurrency: string, rates: Record<string, number>, ratesBase: string, period: AnalyticsPeriod, locale: string): AnalyticsSnapshot => {
	const convert = (amount: number, from: string) => convertWithRates(amount, from, displayCurrency, ratesBase, rates);

	const analyticsSubs = subscriptions.filter((item) => item.include_in_analytics);

	let monthlyTotal = 0;

	for (const subscription of analyticsSubs) {
		monthlyTotal += convert(getMonthlyAmount(subscription.price, subscription.period), subscription.currency);
	}

	monthlyTotal = Math.round(monthlyTotal * 100) / 100;
	const yearlyTotal = Math.round(monthlyTotal * 12 * 100) / 100;
	const dailyAverage = Math.round((monthlyTotal / 30.437) * 100) / 100;
	const weeklyAverage = Math.round(dailyAverage * 7 * 100) / 100;

	const months = period === '1M' ? 1 : period === '3M' ? 3 : period === '6M' ? 6 : 12;
	const projections = buildProjections(analyticsSubs, months, displayCurrency, rates, ratesBase, locale);

	const next30DaysOutflow = outflowWithinDays(analyticsSubs, 30, convert);
	const next90DaysOutflow = outflowWithinDays(analyticsSubs, 90, convert);

	const topSubscriptions: SubscriptionShare[] = analyticsSubs
		.map((subscription) => {
			const monthlyAmount = convert(getMonthlyAmount(subscription.price, subscription.period), subscription.currency);
			return {
				subscription,
				monthlyAmount,
				sharePercent: monthlyTotal > 0 ? Math.round((monthlyAmount / monthlyTotal) * 1000) / 10 : 0,
			};
		})
		.sort((left, right) => right.monthlyAmount - left.monthlyAmount)
		.slice(0, 5);

	const upcomingPayments = analyticsSubs
		.map((subscription) => ({
			subscription,
			amount: convert(subscription.price, subscription.currency),
			daysUntil: getDaysUntil(subscription.date_pay, subscription.period),
		}))
		.filter((item) => item.daysUntil >= 0)
		.sort((left, right) => left.daysUntil - right.daysUntil)
		.slice(0, 6);

	const recommendations = buildRecommendations(subscriptions, monthlyTotal, displayCurrency, convert);
	const categoryBreakdown = buildCategoryBreakdown(analyticsSubs, monthlyTotal, convert);

	return {
		displayCurrency,
		subscriptionCount: subscriptions.length,
		analyticsCount: analyticsSubs.length,
		monthlyTotal,
		yearlyTotal,
		dailyAverage,
		weeklyAverage,
		next30DaysOutflow: Math.round(next30DaysOutflow * 100) / 100,
		next90DaysOutflow: Math.round(next90DaysOutflow * 100) / 100,
		monthlyCount: subscriptions.filter((item) => item.period === 'monthly').length,
		yearlyCount: subscriptions.filter((item) => item.period === 'yearly').length,
		excludedCount: subscriptions.filter((item) => !item.include_in_analytics).length,
		topSubscriptions,
		categoryBreakdown,
		projections,
		recommendations,
		upcomingPayments,
	};
};

export const periodMonths: Record<AnalyticsPeriod, number> = {
	'1M': 1,
	'3M': 3,
	'6M': 6,
	'1Y': 12,
};

export const buildInflationProjections = (points: MonthlyProjectionPoint[], annualInflationPercent: number): MonthlyProjectionPoint[] => {
	const monthlyRate = annualInflationPercent / 100 / 12;

	return points.map((point, index) => ({
		...point,
		amount: Math.round(point.amount * Math.pow(1 + monthlyRate, index) * 100) / 100,
	}));
};
