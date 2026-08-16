import type { AppLanguageCode } from '@/constants/Language.constant';

export type LegalDocumentType = 'terms' | 'privacy' | 'offer' | 'cookies';

export interface LegalSection {
	heading: string;
	paragraphs: string[];
}

export interface LegalDocument {
	title: string;
	updated: string;
	intro: string;
	callout: string;
	sections: LegalSection[];
}

export type LegalDocumentsByLang = Record<AppLanguageCode, LegalDocument>;

export const LEGAL_DOCUMENT_TYPES: LegalDocumentType[] = ['terms', 'privacy', 'offer', 'cookies'];

export const isLegalDocumentType = (value: string | undefined): value is LegalDocumentType => LEGAL_DOCUMENT_TYPES.includes(value as LegalDocumentType);
