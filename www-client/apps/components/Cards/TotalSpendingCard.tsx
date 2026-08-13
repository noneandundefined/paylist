import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import ChevronUp from '@/components/@icons/chevron-up';
import ChevronDown from '@/components/@icons/chevron-down';
import { getCurrencyDisplayParts } from '@/utils/currencyUtils';
import SubscriptionIcon from '@/components/common/TrackedSubscription/SubscriptionIcon';
import type { TrackedSubscriptionResponse, TrackedSubscriptionSummaryResponse } from '@/rest/trackedSubscriptionAPI';

interface TotalSpendingCardProps {
	summary: TrackedSubscriptionSummaryResponse;
	subscriptions: TrackedSubscriptionResponse[];
}

const TotalSpendingCard = ({ summary, subscriptions }: TotalSpendingCardProps) => {
	const { t, i18n } = useTranslation();

	const priceParts = useMemo(() => getCurrencyDisplayParts(summary.total_amount, summary.display_currency, i18n.language), [summary.total_amount, summary.display_currency, i18n.language]);

	const previewSubscriptions = useMemo(
		() => summary.preview_subscription_ids.map((id) => subscriptions.find((subscription) => subscription.id === id)).filter((subscription): subscription is TrackedSubscriptionResponse => Boolean(subscription)),
		[summary.preview_subscription_ids, subscriptions]
	);

	const comparisonMonthLabel = useMemo(() => {
		const previousMonth = new Date();
		previousMonth.setMonth(previousMonth.getMonth() - 1);

		return previousMonth.toLocaleDateString(i18n.language, {
			month: 'long',
			year: 'numeric',
		});
	}, [i18n.language]);

	const comparisonLabel =
		summary.comparison_direction === 'less'
			? t('home.spending-less-than', {
					percent: summary.comparison_percent,
					month: comparisonMonthLabel,
				})
			: t('home.spending-more-than', {
					percent: summary.comparison_percent,
					month: comparisonMonthLabel,
				});

	const isLess = summary.comparison_direction === 'less';
	const hasFraction = Boolean(priceParts.fraction);

	return (
		<section className="gu-glass-card mb-6 p-5">
			<div className="flex items-start justify-between gap-4">
				<div>
					<p className="text-[11px] font-semibold uppercase tracking-[0.08em] gu-text-secondary">{t('home.total-spending')}</p>
					<div className="mt-1 flex items-start gu-text-primary">
						{hasFraction ? (
							<>
								{priceParts.leadingSymbol ? <span className="mt-2 mr-1 text-[30px] font-semibold leading-none">{priceParts.leadingSymbol}</span> : null}
								{priceParts.trailingSymbol ? <span className="mt-2 mr-1 text-[30px] font-semibold leading-none">{priceParts.trailingSymbol}</span> : null}
								<span className="text-[42px] font-bold leading-none tracking-tight">{priceParts.whole}</span>
								<span className="mt-3 text-[22px] font-medium leading-none gu-text-muted">
									{priceParts.fractionSeparator}
									{priceParts.fraction}
								</span>
							</>
						) : (
							<span className="text-[42px] font-bold leading-none tracking-tight">{priceParts.formatted}</span>
						)}
					</div>
				</div>
			</div>

			<div className="mt-5 flex items-end justify-between gap-4">
				<span className={`inline-flex max-w-[65%] items-center gap-1 rounded-full px-2.5 py-1 text-[12px] font-medium leading-tight ${isLess ? 'bg-[#ecfccb] text-[#3f6212]' : 'bg-red-50 text-red-700'}`}>
					{isLess ? <ChevronUp fill="#3f6212" size={16} /> : <ChevronDown fill="#b91c1c" size={16} />}
					{comparisonLabel}
				</span>

				<div className="text-right">
					<div className="mb-1 flex justify-end pl-2">
						{previewSubscriptions.map((subscription, index) => (
							<div key={subscription.id} className={`relative ${index > 0 ? '-ml-2' : ''}`} style={{ zIndex: previewSubscriptions.length - index }}>
								<SubscriptionIcon name={subscription.name} size="xs" className="border-2 border-[var(--surface)]" />
							</div>
						))}
					</div>
					<p className="text-[14px] font-semibold gu-text-primary">{t('home.active-count', { count: summary.active_count })}</p>
					<p className="text-[12px] gu-text-muted">{t('home.subscriptions-label')}</p>
				</div>
			</div>
		</section>
	);
};

export default TotalSpendingCard;
