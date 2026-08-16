import { useEffect, useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Trans, useTranslation } from 'react-i18next';
import { getLegalDocumentPath } from '@/utils/legalDocumentUtils';
import { COOKIE_PREFERENCES_EVENT, getCookieConsent, setCookieConsent, type CookieConsentValue } from '@/utils/cookieConsentUtils';
import { clearYandexMetrikaCookies, hitYandexMetrika, initYandexMetrika } from '@/utils/yandexMetrika';

const CookieConsent: React.FC = () => {
	const { t } = useTranslation();
	const { pathname, search } = useLocation();
	const [consent, setConsent] = useState<CookieConsentValue | null>(() => getCookieConsent());
	const [bannerOpen, setBannerOpen] = useState(() => getCookieConsent() === null);

	useEffect(() => {
		const onOpenPreferences = () => setBannerOpen(true);

		window.addEventListener(COOKIE_PREFERENCES_EVENT, onOpenPreferences);

		return () => window.removeEventListener(COOKIE_PREFERENCES_EVENT, onOpenPreferences);
	}, []);

	useEffect(() => {
		if (consent !== 'accepted') {
			return;
		}

		initYandexMetrika();
		hitYandexMetrika(`${pathname}${search}`);
	}, [consent, pathname, search]);

	const onChoose = (value: CookieConsentValue) => {
		setCookieConsent(value);
		setConsent(value);
		setBannerOpen(false);

		if (value === 'accepted') {
			initYandexMetrika();
			return;
		}

		clearYandexMetrikaCookies();
	};

	if (!bannerOpen) {
		return null;
	}

	return (
		<div className="fixed inset-x-0 bottom-0 z-50 border-t border-[var(--divider)] bg-[var(--surface)] shadow-[0_-8px_32px_var(--shadow-soft)]" role="dialog" aria-labelledby="cookie-consent-title" aria-describedby="cookie-consent-body">
			<div className="mx-auto flex w-full max-w-7xl flex-col gap-4 px-4 py-4 pb-[max(1rem,env(safe-area-inset-bottom))] sm:px-6 lg:flex-row lg:items-center lg:justify-between lg:gap-8 lg:px-8 lg:py-5">
				<div className="min-w-0 max-w-3xl">
					<p id="cookie-consent-title" className="text-[15px] leading-6 gu-text-primary">
						<span className="font-semibold">{t('cookies.title-lead')}</span> {t('cookies.title-rest')}
					</p>
					<p id="cookie-consent-body" className="mt-2 text-[13px] leading-5 gu-text-muted">
						<Trans
							i18nKey="cookies.body"
							components={{
								policy: <Link to={getLegalDocumentPath('cookies')} className="underline" />,
							}}
						/>
					</p>
				</div>

				<div className="flex w-full shrink-0 flex-row gap-2 sm:w-auto">
					<button
						type="button"
						onClick={() => onChoose('rejected')}
						className="flex-1 cursor-pointer whitespace-nowrap rounded-full border border-[var(--divider)] bg-transparent px-5 py-2.5 text-[14px] font-medium gu-text-primary transition hover:bg-[var(--surface-muted)] sm:flex-none"
					>
						{t('cookies.reject-all')}
					</button>
					<button
						type="button"
						onClick={() => onChoose('accepted')}
						className="flex-1 cursor-pointer whitespace-nowrap rounded-full bg-[var(--text-primary)] px-5 py-2.5 text-[14px] font-medium text-[var(--surface)] transition hover:opacity-90 sm:flex-none"
					>
						{t('cookies.accept-all')}
					</button>
				</div>
			</div>
		</div>
	);
};

export default CookieConsent;
