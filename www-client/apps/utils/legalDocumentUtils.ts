import { getAppLanguage, type AppLanguageCode } from '@/constants/Language.constant';

export type LegalDocumentType = 'terms' | 'privacy';

export const getLegalDocumentUrl = (type: LegalDocumentType, language: string): string => {
	const lang: AppLanguageCode = getAppLanguage(language);

	return `/docs/legal/${type}-${lang}.pdf`;
};
