import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { getLegalDocumentPath } from '@/utils/legalDocumentUtils';

const linkClassName = 'gu-text-muted underline hover:text-[var(--text-primary)]';

const AuthTermsFooter: React.FC = () => {
	const { t } = useTranslation();

	return (
		<div className="space-y-2 text-center">
			<p className="text-[12px] leading-relaxed gu-text-muted">
				{t('auth.terms-prefix')}{' '}
				<Link to={getLegalDocumentPath('terms')} target="_blank" rel="noopener noreferrer" className={linkClassName}>
					{t('auth.terms-of-service')}
				</Link>{' '}
				{t('auth.terms-and')}{' '}
				<Link to={getLegalDocumentPath('privacy')} target="_blank" rel="noopener noreferrer" className={linkClassName}>
					{t('auth.privacy-policy')}
				</Link>
			</p>
		</div>
	);
};

export default AuthTermsFooter;
