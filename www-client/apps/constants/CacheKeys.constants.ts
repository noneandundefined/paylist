const indexKey = 'paylist';

export const CACHEKEYs = {
	L_SESSION: `${indexKey}:l:session`,
	AUTH_EMAIL: `${indexKey}:auth:step:email`,
	AUTH_STEP: `${indexKey}:auth:step`,
	PAYLIST_X_REQ_ID: `${indexKey}:settings-client-req`,
	REDIRECT_AUTH_DEVICE: `${indexKey}:redirectUrl:authdevice`,
	SUBSCRIPTION_AD_NEXT_SHOW: `${indexKey}:subscription:ad:nextshow`,
	SUBSCRIPTION_AD_LAST_SHOW: `${indexKey}:subscription:ad:lastshow`,
	DISPLAY_CURRENCY: `${indexKey}:settings:display-currency`,
	THEME: `${indexKey}:settings:theme`,
	COOKIE_CONSENT: `${indexKey}:settings:cookie-consent`,
	IMAGE_MISS: `${indexKey}:image:miss`,
	CURRENCY_LIST: `${indexKey}:currency:list`,
	COUNTRY_LIST: `${indexKey}:country:list`,
	PAID_PAGE_GRANT: `${indexKey}:paid:grant`,
	PAID_PAGE_GOAL: `${indexKey}:paid:goal`,
	SUBSCRIPTION_INVITE_TOKEN: `${indexKey}:subscription:invite-token`,
};
