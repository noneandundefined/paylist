import { buildRoute, ROUTES } from '@/constants/constants';
import type { LegalDocumentType } from '@/content/legal';

export type { LegalDocumentType };

export const getLegalDocumentPath = (type: LegalDocumentType): string => buildRoute(ROUTES.LEGAL, { type });
