import type { AppLanguageCode } from '@/constants/Language.constant';
import { formatCurrency } from '@/utils/currencyUtils';

export const formatPlanAmount = (amount: number, currency: string, language?: string): string => {
	return formatCurrency(amount, currency, language, { maximumFractionDigits: 0 });
};

export const getLocalizedPlanText = (values: Record<string, string> | undefined, language: AppLanguageCode): string => {
	if (!values) {
		return '';
	}

	return values[language] ?? values.en ?? Object.values(values)[0] ?? '';
};

export const getLocalizedPlanFeatures = (values: Record<string, string[]> | undefined, language: AppLanguageCode): string[] => {
	if (!values) {
		return [];
	}

	return values[language] ?? values.en ?? Object.values(values)[0] ?? [];
};
