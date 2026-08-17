import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import SubscriptionIcon from '@/components/common/TrackedSubscription/SubscriptionIcon';
import AccentBadge from '@/components/common/AccentBadge/AccentBadge';
import { formatSubscriptionDate, formatSubscriptionPrice, formatSubscriptionName, getSubscriptionShareAmount } from '@/utils/TrackedSubscriptionDisplayUtils';
import type { TrackedSubscriptionResponse } from '@/rest/trackedSubscriptionAPI';

interface SubscriptionSummaryCardProps {
	subscription: Pick<TrackedSubscriptionResponse, 'name' | 'tariff' | 'price' | 'currency' | 'period' | 'created_at' | 'share_price' | 'share_percent'>;
	periodLabel: string;
	overdue?: boolean;
	trailing?: ReactNode;
}

const SubscriptionSummaryCard: React.FC<SubscriptionSummaryCardProps> = ({ subscription, periodLabel, overdue = false, trailing }) => {
	const { t, i18n } = useTranslation();
	const title = trailing ? subscription.name : formatSubscriptionName(subscription.name, subscription.tariff, t);

	return (
		<div className={`gu-glass-card flex items-center gap-3 p-4 ${overdue ? 'gu-overdue-surface' : ''}`}>
			<SubscriptionIcon name={subscription.name} />

			<div className="min-w-0 flex-1">
				<div className="flex items-center gap-2">
					<p className={`min-w-0 flex-1 truncate text-[19px] font-semibold ${overdue ? 'gu-overdue-title' : 'gu-text-primary'}`}>{title}</p>
					{trailing}
				</div>
				<p className="text-[15px]">
					<AccentBadge className={`text-[19px] font-semibold ${overdue ? 'gu-overdue-badge' : ''}`}>{formatSubscriptionPrice(getSubscriptionShareAmount(subscription), subscription.currency, i18n.language)}</AccentBadge>
					<span className={overdue ? 'gu-overdue-muted' : 'gu-text-muted'}> / {periodLabel}</span>
				</p>
				<p className={`text-[13px] ${overdue ? 'gu-overdue-muted' : 'gu-text-muted'}`}>{t('subscription.active-since', { date: formatSubscriptionDate(subscription.created_at, i18n.language) })}</p>
			</div>
		</div>
	);
};

export default SubscriptionSummaryCard;
