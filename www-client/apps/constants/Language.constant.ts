export const APP_NAME = 'Paylist';

export type AppLanguageCode = 'ru' | 'en' | 'de' | 'es';

export const SUPPORTED_LANGUAGES = [
	{ code: 'ru', labelKey: 'account.language-ru' },
	{ code: 'en', labelKey: 'account.language-en' },
	{ code: 'de', labelKey: 'account.language-de' },
	{ code: 'es', labelKey: 'account.language-es' },
] as const satisfies readonly { code: AppLanguageCode; labelKey: string }[];

export const LANGs = SUPPORTED_LANGUAGES.map((language) => language.code);

export const DEFAULT_LANGUAGE: AppLanguageCode = 'ru';

export const getAppLanguage = (language: string): AppLanguageCode => {
	const code = language.slice(0, 2) as AppLanguageCode;

	return LANGs.includes(code) ? code : DEFAULT_LANGUAGE;
};
