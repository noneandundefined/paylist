import type { TFunction } from 'i18next';

export const getCountryLabel = (code: string, fallbackName: string, t: TFunction): string => {
	return t(`country.${code}`, fallbackName);
};
