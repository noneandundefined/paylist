import { getAppLanguage, type AppLanguageCode } from '@/constants/Language.constant';

const CURRENCY_COUNTRY: Record<string, string> = {
	AUD: 'au',
	BGN: 'bg',
	BRL: 'br',
	CAD: 'ca',
	CHF: 'ch',
	CNY: 'cn',
	CZK: 'cz',
	DKK: 'dk',
	EUR: 'eu',
	GBP: 'gb',
	HKD: 'hk',
	HUF: 'hu',
	IDR: 'id',
	ILS: 'il',
	INR: 'in',
	ISK: 'is',
	JPY: 'jp',
	KRW: 'kr',
	MXN: 'mx',
	MYR: 'my',
	NOK: 'no',
	NZD: 'nz',
	PHP: 'ph',
	PLN: 'pl',
	RON: 'ro',
	RUB: 'ru',
	SEK: 'se',
	SGD: 'sg',
	THB: 'th',
	TRY: 'tr',
	USD: 'us',
	ZAR: 'za',
};

const APP_LOCALE: Record<AppLanguageCode, string> = {
	ru: 'ru-RU',
	en: 'en-US',
	de: 'de-DE',
	es: 'es-ES',
};

const CURRENCY_FORMAT_LOCALE: Partial<Record<string, string>> = {
	USD: 'en-US',
	EUR: 'de-DE',
	GBP: 'en-GB',
	RUB: 'ru-RU',
	SGD: 'en-SG',
	JPY: 'ja-JP',
	KRW: 'ko-KR',
	CNY: 'zh-CN',
	CHF: 'de-CH',
	AUD: 'en-AU',
	CAD: 'en-CA',
	HKD: 'en-HK',
	INR: 'en-IN',
	BRL: 'pt-BR',
	MXN: 'es-MX',
	TRY: 'tr-TR',
	PLN: 'pl-PL',
	SEK: 'sv-SE',
	NOK: 'nb-NO',
	DKK: 'da-DK',
	THB: 'th-TH',
	IDR: 'id-ID',
	PHP: 'en-PH',
	MYR: 'ms-MY',
	NZD: 'en-NZ',
	ZAR: 'en-ZA',
	ILS: 'he-IL',
	CZK: 'cs-CZ',
	HUF: 'hu-HU',
	RON: 'ro-RO',
	BGN: 'bg-BG',
	ISK: 'is-IS',
};

/** Overrides where Intl narrowSymbol is ambiguous (e.g. SGD → $ instead of S$). */
const CURRENCY_SYMBOL_OVERRIDE: Partial<Record<string, string>> = {
	SGD: 'S$',
	AUD: 'A$',
	CAD: 'CA$',
	HKD: 'HK$',
	NZD: 'NZ$',
	MXN: 'MX$',
};

export interface CurrencyDisplayParts {
	leadingSymbol: string;
	whole: string;
	fractionSeparator: string;
	fraction: string;
	trailingSymbol: string;
	formatted: string;
}

const getCurrencyFormatLocale = (currency: string, language?: string): string => {
	const code = currency.toUpperCase();

	return CURRENCY_FORMAT_LOCALE[code] ?? APP_LOCALE[getAppLanguage(language ?? '')] ?? 'en-US';
};

type CurrencyFractionDigits = Pick<Intl.NumberFormatOptions, 'minimumFractionDigits' | 'maximumFractionDigits'>;

const createCurrencyFormatter = (_: number, currency: string, language?: string, fractionDigits?: CurrencyFractionDigits) => {
	const code = currency.toUpperCase();
	const formatLocale = getCurrencyFormatLocale(code, language);

	return new Intl.NumberFormat(formatLocale, {
		style: 'currency',
		currency: code,
		currencyDisplay: 'narrowSymbol',
		...(fractionDigits ?? {}),
	});
};

const applyCurrencySymbolOverride = (parts: Intl.NumberFormatPart[], currency: string): string => {
	const override = CURRENCY_SYMBOL_OVERRIDE[currency.toUpperCase()];

	if (!override) {
		return parts.map((part) => part.value).join('');
	}

	return parts.map((part) => (part.type === 'currency' ? override : part.value)).join('');
};

export const formatCurrency = (amount: number, currency: string, language?: string, fractionDigits?: CurrencyFractionDigits): string => {
	const code = currency.toUpperCase();
	const formatter = createCurrencyFormatter(amount, code, language, fractionDigits);

	return applyCurrencySymbolOverride(formatter.formatToParts(amount), code);
};

export const getCurrencyDisplayParts = (amount: number, currency: string, language?: string): CurrencyDisplayParts => {
	const code = currency.toUpperCase();
	const formatter = createCurrencyFormatter(amount, code, language);
	const parts = formatter.formatToParts(amount);
	const override = CURRENCY_SYMBOL_OVERRIDE[code];

	let leadingSymbol = '';
	let whole = '';
	let fractionSeparator = '';
	let fraction = '';
	let trailingSymbol = '';
	let seenNumber = false;

	for (const part of parts) {
		if (part.type === 'currency') {
			const symbol = override ?? part.value;

			if (!seenNumber) {
				leadingSymbol += symbol;
			} else {
				trailingSymbol += symbol;
			}

			continue;
		}

		if (part.type === 'integer' || part.type === 'group') {
			seenNumber = true;
			whole += part.value;
			continue;
		}

		if (part.type === 'decimal') {
			fractionSeparator = part.value;
			continue;
		}

		if (part.type === 'fraction') {
			fraction = part.value;
		}
	}

	return {
		leadingSymbol,
		whole,
		fractionSeparator,
		fraction,
		trailingSymbol,
		formatted: applyCurrencySymbolOverride(parts, code),
	};
};

export const getCurrencyFlagUrl = (code: string): string | null => {
	const countryCode = CURRENCY_COUNTRY[code.toUpperCase()];

	if (!countryCode) {
		return null;
	}

	return `https://flagcdn.com/w40/${countryCode}.png`;
};
