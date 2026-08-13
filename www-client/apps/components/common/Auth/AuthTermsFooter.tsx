import { useTranslation } from 'react-i18next';
import { getLegalDocumentUrl } from '@/utils/legalDocumentUtils';

const linkClassName = 'gu-text-muted underline hover:text-[var(--text-primary)]';

const AuthTermsFooter: React.FC = () => {
	const { t, i18n } = useTranslation();

	const termsUrl = getLegalDocumentUrl('terms', i18n.language);
	const privacyUrl = getLegalDocumentUrl('privacy', i18n.language);

	return (
		<p className="text-center text-[12px] leading-relaxed gu-text-muted">
			{t('auth.terms-prefix')}{' '}
			<a href={termsUrl} target="_blank" rel="noopener noreferrer" className={linkClassName}>
				{t('auth.terms-of-service')}
			</a>{' '}
			{t('auth.terms-and')}{' '}
			<a href={privacyUrl} target="_blank" rel="noopener noreferrer" className={linkClassName}>
				{t('auth.privacy-policy')}
			</a>
		</p>
	);
};

export default AuthTermsFooter;
