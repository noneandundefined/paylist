import { apiGet } from '@/rest/apiClient';

export interface CurrencyItem {
	code: string;
	name: string;
}

export interface CurrencyConvertResponse {
	from: string;
	to: string;
	amount: number;
	result: number;
}

export interface CurrencyRatesResponse {
	base: string;
	rates: Record<string, number>;
}

export const fetchCurrencies = async (): Promise<CurrencyItem[]> => {
	const items = await apiGet<CurrencyItem[]>('/currency/currencies');
	return items ?? [];
};

export const loadCurrencySelectOptions = async () => {
	const items = await fetchCurrencies();

	return items.map((item) => ({
		value: item.code,
		label: item.name,
		description: `( ${item.code} )`,
		keywords: `${item.code} ${item.name}`,
	}));
};

export const basicCurrencyConvert = async (from: string, to: string, amount: number): Promise<CurrencyConvertResponse> =>
	apiGet('/currency/convert', {
		params: { from, to, amount },
	});

export const basicCurrencyRates = async (base: string, symbols?: string[]): Promise<CurrencyRatesResponse> =>
	apiGet('/currency/rates', {
		params: {
			base,
			symbols: symbols?.join(','),
		},
	});
