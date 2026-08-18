export const PLACEHOLDER_ROUTE_ID = ':id';

export const ROUTES = {
	HOME: '/',
	START: '/start',
	NOT_FOUND: '*',
	SIGNIN: '/sign-in',
	SIGNUP: '/sign-up',
	SUBSCRIPTIONS: '/subscriptions',
	SUBSCRIPTION_CREATE: '/subscriptions/create',
	SUBSCRIPTION_DETAIL: `/subscriptions/${PLACEHOLDER_ROUTE_ID}`,
	ACCOUNT: '/account',
	ANALYTICS: '/analytics',
	PLANS: '/plans',
	PAID: '/paid',
	LEGAL: '/legal/:type',

	AUTH_AUTHORIZE_DEVICE: '/authorize',
	AUTH_CONFIRM_EMAIL: '/paylist-confirm-email',
	AUTH_SENT_CONFIRM_EMAIL: '/paylist-sent-confirm-email',
	AUTH_FORGOT_PASSWORD: '/forgot-password',
	AUTH_RESET_PASSWORD: '/paylist-reset-password',
	SUBSCRIPTION_INVITE: '/paylist-subscription-invite',
};

type Params = Record<string, string | number>;

export const buildRoute = (route: string, params: Params): string => {
	let result = route;

	Object.entries(params).forEach(([key, value]) => {
		result = result.replace(`:${key}`, encodeURIComponent(String(value)));
	});

	return result;
};
