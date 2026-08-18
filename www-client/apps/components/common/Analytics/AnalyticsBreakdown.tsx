import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { ROUTES, buildRoute } from '@/constants/constants';
import { formatSubscriptionPrice, formatSubscriptionName } from '@/utils/TrackedSubscriptionDisplayUtils';
import SubscriptionIcon from '@/components/common/TrackedSubscription/SubscriptionIcon';
import type { SubscriptionShare } from '@/utils/analyticsUtils';

interface AnalyticsBreakdownProps {
	items: SubscriptionShare[];
	currency: string;
}

const AnalyticsBreakdown: React.FC<AnalyticsBreakdownProps> = ({ items, currency }) => {
	const { t, i18n } = useTranslation();

	if (items.length === 0) {
		return null;
	}

	return (
		<section className="gu-glass-card p-5">
			<h2 className="text-[15px] font-semibold gu-text-primary">{t('analytics.breakdown-title')}</h2>
			<p className="mt-1 text-[13px] gu-text-muted">{t('analytics.breakdown-subtitle')}</p>

			<ul className="mt-4 space-y-4">
				{items.map((item) => (
					<li key={item.subscription.id}>
						<Link to={buildRoute(ROUTES.SUBSCRIPTION_DETAIL, { id: item.subscription.id })} className="flex items-center gap-3 no-underline hover:no-underline">
							<SubscriptionIcon name={item.subscription.name} categories={item.subscription.categories} size="sm" />
							<div className="min-w-0 flex-1">
								<div className="flex items-center justify-between gap-2">
									<p className="truncate text-[14px] font-medium gu-text-primary">{formatSubscriptionName(item.subscription.name, item.subscription.tariff, t)}</p>
									<p className="shrink-0 text-[13px] font-semibold gu-text-primary">{formatSubscriptionPrice(item.monthlyAmount, currency, i18n.language)}</p>
								</div>
								<div className="mt-2 h-2 overflow-hidden rounded-full bg-[var(--surface-muted)]">
									<div className="h-full rounded-full bg-[#0085FF]" style={{ width: `${Math.min(item.sharePercent, 100)}%` }} />
								</div>
								<p className="mt-1 text-[12px] gu-text-muted">{t('analytics.share-of-budget', { percent: item.sharePercent })}</p>
							</div>
						</Link>
					</li>
				))}
			</ul>
		</section>
	);
};

export default AnalyticsBreakdown;
