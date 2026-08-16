import type { AppLanguageCode } from '@/constants/Language.constant';
import { cookiesDocuments } from './cookies';
import { offerDocuments } from './offer';
import { privacyDocuments } from './privacy';
import { termsDocuments } from './terms';
import type { LegalDocument, LegalDocumentType } from './types';

export type { LegalDocument, LegalDocumentType, LegalSection } from './types';
export { isLegalDocumentType, LEGAL_DOCUMENT_TYPES } from './types';

const documents: Record<LegalDocumentType, Record<AppLanguageCode, LegalDocument>> = {
	terms: termsDocuments,
	privacy: privacyDocuments,
	offer: offerDocuments,
	cookies: cookiesDocuments,
};

export const getLegalDocument = (type: LegalDocumentType, language: AppLanguageCode): LegalDocument => documents[type][language];
