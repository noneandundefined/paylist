import { useTranslation } from 'react-i18next';
import { formatSubscriptionPrice } from '@/utils/TrackedSubscriptionDisplayUtils';
import type { AnalyticsSnapshot } from '@/utils/analyticsUtils';

interface AnalyticsStatRowProps {
	snapshot: AnalyticsSnapshot;
}

const statCardClass = 'gu-glass-card p-4 sm:p-5';
const primaryLabelClass = 'text-[11px] font-semibold uppercase tracking-[0.08em] gu-text-secondary sm:text-xs lg:text-sm';
const primaryValueClass = 'mt-1 text-[22px] font-bold leading-tight gu-text-primary sm:text-[28px] lg:text-[34px]';
const hintClass = 'mt-0.5 text-[11px] gu-text-muted sm:text-xs lg:text-sm';
const secondaryLabelClass = 'text-[11px] gu-text-muted sm:text-xs lg:text-sm';
const secondaryValueClass = 'mt-1 text-[16px] font-semibold leading-tight gu-text-primary sm:text-lg lg:text-2xl';

const AnalyticsStatRow: React.FC<AnalyticsStatRowProps> = ({ snapshot }) => {
	const { t, i18n } = useTranslation();
	const { displayCurrency } = snapshot;

	return (
		<div className="space-y-3">
			<div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
				<article className={statCardClass}>
					<p className={primaryLabelClass}>{t('analytics.stat-monthly')}</p>
					<p className={primaryValueClass}>{formatSubscriptionPrice(snapshot.monthlyTotal, displayCurrency, i18n.language)}</p>
				</article>

				<article className={statCardClass}>
					<p className={primaryLabelClass}>{t('analytics.stat-yearly')}</p>
					<p className={primaryValueClass}>{formatSubscriptionPrice(snapshot.yearlyTotal, displayCurrency, i18n.language)}</p>
				</article>

				<article className={statCardClass}>
					<p className={primaryLabelClass}>{t('analytics.stat-subscriptions')}</p>
					<p className={primaryValueClass}>{snapshot.subscriptionCount}</p>
					<p className={hintClass}>{t('analytics.stat-in-analytics', { count: snapshot.analyticsCount })}</p>
				</article>
			</div>

			<div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
				<article className={statCardClass}>
					<p className={secondaryLabelClass}>{t('analytics.stat-daily')}</p>
					<p className={secondaryValueClass}>{formatSubscriptionPrice(snapshot.dailyAverage, displayCurrency, i18n.language)}</p>
				</article>

				<article className={statCardClass}>
					<p className={secondaryLabelClass}>{t('analytics.stat-weekly')}</p>
					<p className={secondaryValueClass}>{formatSubscriptionPrice(snapshot.weeklyAverage, displayCurrency, i18n.language)}</p>
				</article>

				<article className={statCardClass}>
					<p className={secondaryLabelClass}>{t('analytics.stat-next-30')}</p>
					<p className={secondaryValueClass}>{formatSubscriptionPrice(snapshot.next30DaysOutflow, displayCurrency, i18n.language)}</p>
				</article>

				<article className={statCardClass}>
					<p className={secondaryLabelClass}>{t('analytics.stat-next-90')}</p>
					<p className={secondaryValueClass}>{formatSubscriptionPrice(snapshot.next90DaysOutflow, displayCurrency, i18n.language)}</p>
				</article>
			</div>
		</div>
	);
};

export default AnalyticsStatRow;
