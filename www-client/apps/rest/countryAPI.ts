import { apiGet } from '@/rest/apiClient';
import i18n from '@/utils/i18n';
import { CACHEKEYs } from '@/constants/CacheKeys.constants';
import { getCountryLabel } from '@/utils/countryDisplayUtils';
import { DAY_MS, withTtlCache } from '@/utils/ttlStorage';

export interface CountryItem {
	code: string;
	name: string;
	inflation_rate: number;
}

const isCountryList = (value: unknown): value is CountryItem[] =>
	Array.isArray(value) &&
	value.length > 0 &&
	value.every(
		(item) =>
			Boolean(item) &&
			typeof item === 'object' &&
			typeof (item as CountryItem).code === 'string' &&
			typeof (item as CountryItem).name === 'string' &&
			typeof (item as CountryItem).inflation_rate === 'number'
	);

export const fetchCountries = async (): Promise<CountryItem[]> =>
	withTtlCache(CACHEKEYs.COUNTRY_LIST, DAY_MS, async () => (await apiGet<CountryItem[]>('/country/countries')) ?? [], isCountryList);

export const loadCountrySelectOptions = async () => {
	const items = await fetchCountries();
	const t = i18n.t.bind(i18n);

	return items
		.map((item) => {
			const label = getCountryLabel(item.code, item.name, t);

			return {
				value: item.code,
				label,
				description: `( ${item.code} · ~${item.inflation_rate}% )`,
				keywords: `${item.code} ${item.name} ${label}`,
			};
		})
		.sort((left, right) => left.label.localeCompare(right.label, i18n.language));
};

export const getCountryInflationRate = (countries: CountryItem[], code: string): number => {
	const match = countries.find((item) => item.code.toUpperCase() === code.toUpperCase());
	return match?.inflation_rate ?? 5;
};
