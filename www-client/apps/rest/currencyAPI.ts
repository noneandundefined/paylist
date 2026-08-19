import { apiGet } from '@/rest/apiClient';
import { CACHEKEYs } from '@/constants/CacheKeys.constants';
import { DAY_MS, withTtlCache } from '@/utils/ttlStorage';

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

const isCurrencyList = (value: unknown): value is CurrencyItem[] =>
	Array.isArray(value) && value.length > 0 && value.every((item) => Boolean(item) && typeof item === 'object' && typeof (item as CurrencyItem).code === 'string' && typeof (item as CurrencyItem).name === 'string');

export const fetchCurrencies = async (): Promise<CurrencyItem[]> => withTtlCache(CACHEKEYs.CURRENCY_LIST, DAY_MS, async () => (await apiGet<CurrencyItem[]>('/currency/currencies')) ?? [], isCurrencyList);

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
