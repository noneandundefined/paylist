export const SUBSCRIPTION_TARIFF_NONE = 'none';

export const SUBSCRIPTION_TARIFFS = ['none', 'basic', 'standard', 'plus', 'pro', 'premium', 'max', 'lite', 'mini', 'student', 'duo', 'family', 'individual', 'business'] as const;

export type SubscriptionTariff = (typeof SUBSCRIPTION_TARIFFS)[number];

export const hasSubscriptionTariff = (tariff?: string | null): boolean => Boolean(tariff && tariff !== SUBSCRIPTION_TARIFF_NONE);
