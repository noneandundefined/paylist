import { useTranslation } from 'react-i18next';
import SubscriptionIcon from '@/components/common/TrackedSubscription/SubscriptionIcon';
import AccentBadge from '@/components/common/AccentBadge/AccentBadge';
import { formatSubscriptionDate, formatSubscriptionPrice } from '@/utils/TrackedSubscriptionDisplayUtils';
import type { TrackedSubscriptionResponse } from '@/rest/trackedSubscriptionAPI';

interface SubscriptionSummaryCardProps {
	subscription: Pick<TrackedSubscriptionResponse, 'name' | 'price' | 'currency' | 'period' | 'created_at'>;
	periodLabel: string;
	overdue?: boolean;
}

const SubscriptionSummaryCard: React.FC<SubscriptionSummaryCardProps> = ({ subscription, periodLabel, overdue = false }) => {
	const { t, i18n } = useTranslation();

	return (
		<div className={`gu-glass-card flex items-center gap-3 p-4 ${overdue ? 'gu-overdue-surface' : ''}`}>
			<SubscriptionIcon name={subscription.name} />

			<div className="min-w-0 flex-1">
				<p className={`text-[19px] font-semibold ${overdue ? 'gu-overdue-title' : 'gu-text-primary'}`}>{subscription.name}</p>
				<p className="text-[15px]">
					<AccentBadge className={`text-[19px] font-semibold ${overdue ? 'gu-overdue-badge' : ''}`}>{formatSubscriptionPrice(subscription.price, subscription.currency, i18n.language)}</AccentBadge>
					<span className={overdue ? 'gu-overdue-muted' : 'gu-text-muted'}> / {periodLabel}</span>
				</p>
				<p className={`text-[13px] ${overdue ? 'gu-overdue-muted' : 'gu-text-muted'}`}>{t('subscription.active-since', { date: formatSubscriptionDate(subscription.created_at, i18n.language) })}</p>
			</div>
		</div>
	);
};

export default SubscriptionSummaryCard;
