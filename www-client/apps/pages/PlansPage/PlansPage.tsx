import Check from '@/components/@icons/check';
import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { getLegalDocumentPath } from '@/utils/legalDocumentUtils';
import { ROUTES } from '@/constants/constants';
import { basicPlanList } from '@/rest/planAPI';
import { APP_NAME, getAppLanguage } from '@/constants/Language.constant';
import GUIButton from '@/components/ui/Button/GUIButton';
import { basicPaymentCheckout } from '@/rest/paymentAPI';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { useLoginState } from '@/hooks/useLoginState';
import { QUERY_KEYS } from '@/constants/QueryKeys.constants';
import { PREMIUM_GRADIENT, PREMIUM_GRADIENT_FEATURE } from '@/constants/designTokens';
import PageHeader from '@/components/common/PageHeader/PageHeader';
import AccentBadge from '@/components/common/AccentBadge/AccentBadge';
import { formatPlanAmount, getLocalizedPlanFeatures, getLocalizedPlanText } from '@/utils/planDisplayUtils';
import Fallback from '@/components/Fallback/Fallback';

interface PlansHeroCardProps {
	planName: string;
	amountLabel: string;
}

const PlansHeroCard = ({ planName, amountLabel }: PlansHeroCardProps) => {
	return (
		<div className="relative mx-auto mt-2 w-full max-w-[320px] rotate-[-2deg]">
			<div className="absolute -left-3 top-8 h-16 w-16 rounded-full bg-[#0085FF]/15 blur-2xl" />
			<div className="absolute -right-2 bottom-2 h-20 w-20 rounded-full bg-[#60D1FF]/20 blur-2xl" />

			<div className="relative overflow-hidden rounded-[22px] border border-white/80 p-5 shadow-[0_16px_40px_rgba(0,133,255,0.22)]" style={{ background: PREMIUM_GRADIENT }}>
				<div className="mb-4 flex items-center justify-between">
					<span className="rounded-full bg-white/20 px-3 py-1 text-[11px] font-semibold uppercase tracking-wide text-white">{APP_NAME}</span>
					<span className="text-[11px] font-medium uppercase tracking-[0.18em] text-white/75">{planName}</span>
				</div>

				<div className="space-y-2">
					<div className="h-2.5 w-24 rounded-full bg-white/25" />
					<div className="h-2.5 w-40 rounded-full bg-white/20" />
					<div className="h-2.5 w-32 rounded-full bg-white/20" />
				</div>

				<div className="mt-6 flex items-end justify-between">
					<div>
						<p className="text-[11px] uppercase tracking-wide text-white/70">Total / month</p>
						<p className="text-[28px] font-bold leading-none text-white">{amountLabel}</p>
					</div>

					<div className="flex gap-1.5">
						<AccentBadge className="inline-flex h-8 w-8 items-center justify-center rounded-full text-[11px] font-bold">₽</AccentBadge>
						<span className="inline-flex h-8 w-8 items-center justify-center rounded-full bg-white/20 text-[11px] font-bold text-white">★</span>
					</div>
				</div>
			</div>
		</div>
	);
};

const PlansPage = () => {
	const { t, i18n } = useTranslation();
	const navigate = useNavigate();
	const language = getAppLanguage(i18n.language);

	const { isPremium, loading: loginLoading } = useLoginState();
	const { data: plans, loading: plansLoading } = useHandleServer([QUERY_KEYS.planList], () => basicPlanList());
	const [acceptedOffer, setAcceptedOffer] = useState(false);
	const offerPath = getLegalDocumentPath('offer');

	const plan = plans?.[0];
	const loading = loginLoading || plansLoading;

	const planAmount = plan ? formatPlanAmount(plan.amount, plan.currency, i18n.language) : '';
	const priceLabel = plan ? t('plans.price-per-month', { amount: planAmount }) : '';
	const planDescription = plan ? getLocalizedPlanText(plan.description, language) : '';
	const planFeatures = plan ? getLocalizedPlanFeatures(plan.features, language) : [];

	const onClose = () => {
		if (window.history.length > 1) {
			navigate(-1);
			return;
		}

		navigate(ROUTES.ACCOUNT);
	};

	const onSubscribe = async () => {
		if (!plan) {
			return;
		}

		const checkout = await basicPaymentCheckout({ plan_name: plan.plan_name });
		window.location.assign(checkout.confirmation_url);
	};

	if (loading) {
		return (
			<div className="min-h-screen gu-page-bg">
				<Fallback text={t('message.page-loading')} />
			</div>
		);
	}

	if (!plan) {
		return (
			<div className="min-h-screen gu-page-bg px-4 py-5">
				<div className="mx-auto flex w-full max-w-lg flex-col">
					<PageHeader variant="close" onClose={onClose} backLabel={t('action.close')} title="" />
					<p className="mt-16 text-center text-[15px] gu-text-muted">{t('plans.not-available')}</p>
				</div>
			</div>
		);
	}

	return (
		<div className="min-h-screen gu-page-bg gu-text-primary">
			<div className="mx-auto flex min-h-screen w-full max-w-7xl flex-col px-4 pb-8 pt-5">
				<PageHeader variant="close" onClose={onClose} backLabel={t('action.close')} title="" />

				<PlansHeroCard planName={plan.plan_name} amountLabel={planAmount} />

				<div className="mt-8 text-center">
					<h1 className="text-[34px] font-bold leading-tight tracking-tight gu-text-primary">
						{APP_NAME} <span className="gu-premium-gradient-text font-serif italic">premium</span>
					</h1>
					<p className="mt-2 text-[15px] gu-text-muted">{planDescription || t('plans.subtitle')}</p>
				</div>

				<div className="gu-glass-card mt-8 rounded-[18px] border-2 border-[#0085FF] p-4">
					<div className="flex items-center justify-between gap-3">
						<div>
							<p className="text-[18px] font-semibold gu-text-primary">{t('plans.monthly')}</p>
							<p className="mt-1 text-[13px] gu-text-muted">{t('plans.monthly-description', { days: plan.duration_days })}</p>
						</div>

						<div className="flex items-center gap-3">
							<p className="text-[15px] font-medium gu-text-secondary">{priceLabel}</p>
							<span className="inline-flex h-5 w-5 shrink-0 items-center justify-center rounded-full border-[5px] border-[#0085FF] bg-[var(--surface)]" aria-hidden />
						</div>
					</div>
				</div>

				{planFeatures.length > 0 && (
					<section className="mt-8">
						<h2 className="text-[18px] font-semibold gu-text-primary">{t('plans.benefits-title')}</h2>

						<ul className="mt-4 space-y-4">
							{planFeatures.map((feature) => (
								<li key={feature} className="flex items-start gap-3">
									<span className="inline-flex h-10 w-10 shrink-0 items-center justify-center rounded-2xl shadow-[0_8px_20px_rgba(0,133,255,0.18)]" style={{ background: PREMIUM_GRADIENT_FEATURE }}>
										<Check fill="#fff" size={18} />
									</span>

									<p className="min-w-0 pt-2 text-[15px] font-medium leading-snug gu-text-secondary">{feature}</p>
								</li>
							))}
						</ul>
					</section>
				)}

				<div className="mt-auto pt-10">
					{!isPremium ? (
						<label className="mb-4 flex items-start gap-3 text-left text-[13px] leading-relaxed gu-text-muted">
							<input type="checkbox" className="mt-1 h-4 w-4 shrink-0" checked={acceptedOffer} onChange={(event) => setAcceptedOffer(event.target.checked)} />
							<span>
								{t('plans.offer-accept-prefix')}{' '}
								<Link to={offerPath} target="_blank" rel="noopener noreferrer" className="underline gu-text-primary">
									{t('plans.offer-title')}
								</Link>
								{t('plans.offer-accept-suffix', { price: planAmount, days: plan.duration_days })}
							</span>
						</label>
					) : null}

					<GUIButton
						type="button"
						disabled={isPremium || !acceptedOffer}
						onClick={onSubscribe}
						className="w-full rounded-[18px] px-6 py-4 text-[16px] font-semibold text-white shadow-[0_10px_30px_rgba(0,133,255,0.25)] transition hover:brightness-105 disabled:cursor-not-allowed disabled:opacity-50"
						style={{ background: PREMIUM_GRADIENT }}
					>
						{isPremium ? t('plans.already-subscribed') : t('plans.subscribe-monthly')}
					</GUIButton>

					<p className="mt-3 text-center text-[12px] leading-relaxed gu-text-muted">{t('plans.disclaimer', { price: planAmount, days: plan.duration_days })}</p>
				</div>
			</div>
		</div>
	);
};

export default PlansPage;
