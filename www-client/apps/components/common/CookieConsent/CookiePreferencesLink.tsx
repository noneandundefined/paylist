import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { COOKIE_CONSENT_CHANGE_EVENT, getCookieConsent, openCookiePreferences } from '@/utils/cookieConsentUtils';

interface CookiePreferencesLinkProps {
	className?: string;
}

const CookiePreferencesLink: React.FC<CookiePreferencesLinkProps> = ({ className = '' }) => {
	const { t } = useTranslation();
	const [visible, setVisible] = useState(() => Boolean(getCookieConsent()));

	useEffect(() => {
		const onChange = () => setVisible(Boolean(getCookieConsent()));

		window.addEventListener(COOKIE_CONSENT_CHANGE_EVENT, onChange);

		return () => window.removeEventListener(COOKIE_CONSENT_CHANGE_EVENT, onChange);
	}, []);

	if (!visible) {
		return null;
	}

	return (
		<button type="button" onClick={openCookiePreferences} className={`cursor-pointer border-0 bg-transparent p-0 text-[12px] gu-text-muted underline hover:text-[var(--text-primary)] ${className}`}>
			{t('cookies.preferences')}
		</button>
	);
};

export default CookiePreferencesLink;
