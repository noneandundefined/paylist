import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Helmet } from 'react-helmet-async';
import { useTranslation } from 'react-i18next';
import Check from '@/components/@icons/check';
import GUIButton from '@/components/ui/Button/GUIButton';
import { ROUTES } from '@/constants/constants';
import { basicPaymentConfirm } from '@/rest/paymentAPI';
import { hasPaidPageGrant, writePaidPageGrant } from '@/utils/paidPageAccessUtils';
import { trackPaidPageConversion } from '@/utils/yandexMetrika';

const POLL_ATTEMPTS = 12;
const POLL_INTERVAL_MS = 1500;

const sleep = (ms: number) => new Promise((resolve) => window.setTimeout(resolve, ms));

const PaidPage = () => {
	const { t } = useTranslation();
	const navigate = useNavigate();
	const [searchParams] = useSearchParams();
	const paymentId = searchParams.get('payment_id')?.trim() || '';
	const [ready, setReady] = useState(() => hasPaidPageGrant(paymentId || null));

	useEffect(() => {
		if (ready) {
			return;
		}

		if (!paymentId) {
			navigate(ROUTES.HOME, { replace: true });
			return;
		}

		let cancelled = false;

		const confirmPayment = async () => {
			for (let attempt = 0; attempt < POLL_ATTEMPTS; attempt += 1) {
				try {
					const result = await basicPaymentConfirm(paymentId);

					if (cancelled) {
						return;
					}

					if (result.paid) {
						writePaidPageGrant(paymentId);
						trackPaidPageConversion();
						setReady(true);
						navigate(ROUTES.PAID, { replace: true });
						return;
					}

					if (!result.status || result.status === 'canceled' || result.status === 'failed') {
						navigate(ROUTES.HOME, { replace: true });
						return;
					}
				} catch {
					if (cancelled) {
						return;
					}

					if (attempt === POLL_ATTEMPTS - 1) {
						navigate(ROUTES.HOME, { replace: true });
						return;
					}
				}

				await sleep(POLL_INTERVAL_MS);
			}

			if (!cancelled) {
				navigate(ROUTES.HOME, { replace: true });
			}
		};

		void confirmPayment();

		return () => {
			cancelled = true;
		};
	}, [navigate, paymentId, ready]);

	if (!ready) {
		return <div className="min-h-screen bg-[var(--surface)]" />;
	}

	return (
		<div className="flex min-h-screen flex-col bg-[var(--surface)] px-6">
			<Helmet>
				<meta name="robots" content="noindex, nofollow" />
			</Helmet>

			<div className="h-screen flex flex-col items-center justify-center w-full space-y-[5rem]">
				<div className="mx-auto flex w-full max-w-md flex-col items-center text-center">
					<div className="flex h-[88px] w-[88px] items-center justify-center rounded-full bg-[#22c55e]" aria-hidden>
						<Check fill="#fff" size={44} />
					</div>

					<h1 className="mt-8 text-[32px] font-bold leading-tight tracking-tight gu-text-primary">{t('paid.title')}</h1>
					<p className="mt-3 max-w-xs text-[16px] leading-relaxed gu-text-muted">{t('paid.subtitle')}</p>
				</div>

				<div className="w-[30rem]">
					<GUIButton type="button" onClick={() => navigate(ROUTES.HOME)} className="w-full rounded-full bg-[var(--surface-muted)] py-4 text-[16px] font-semibold gu-text-primary transition hover:opacity-90">
						{t('paid.cta')}
					</GUIButton>
				</div>
			</div>
		</div>
	);
};

export default PaidPage;
