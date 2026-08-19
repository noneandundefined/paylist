import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { ROUTES, buildRoute } from '@/constants/constants';
import AccentBadge from '@/components/common/AccentBadge/AccentBadge';
import SubscriptionIcon from '@/components/common/TrackedSubscription/SubscriptionIcon';
import { getRenewalBadgeDays, isSubscriptionOverdue } from '@/utils/SubscriptionRenewalUtils';
import type { TrackedSubscriptionResponse } from '@/rest/trackedSubscriptionAPI';
import { formatSubscriptionPrice, getRenewOnDateLabel, getSubscriptionShareAmount, formatSubscriptionName } from '@/utils/TrackedSubscriptionDisplayUtils';

interface TrackedSubscriptionCardProps {
	subscription: TrackedSubscriptionResponse;
}

const TrackedSubscriptionCard: React.FC<TrackedSubscriptionCardProps> = ({ subscription }) => {
	const { t, i18n } = useTranslation();
	const overdue = isSubscriptionOverdue(subscription.date_pay);
	const renewalBadgeDays = overdue ? null : getRenewalBadgeDays(subscription.date_pay);
	const renewalBadgeLabel = overdue ? t('home.overdue-badge') : renewalBadgeDays === null ? null : renewalBadgeDays === 0 ? t('home.renew-badge-today') : t('home.renew-badge-in-days', { days: renewalBadgeDays });

	return (
		<Link
			to={buildRoute(ROUTES.SUBSCRIPTION_DETAIL, { id: subscription.id })}
			className={`gu-glass-card flex min-w-0 w-full items-center gap-3 overflow-hidden px-4 py-5 text-left no-underline transition hover:no-underline ${overdue ? 'gu-overdue-surface' : 'hover:bg-[var(--surface-muted)] gu-text-primary'}`}
		>
			<SubscriptionIcon name={subscription.name} categories={subscription.categories} />

			<div className="min-w-0 flex-1 overflow-hidden">
				<div className="grid min-w-0 grid-cols-[minmax(0,1fr)_auto] items-center gap-2">
					<p className={`overflow-hidden text-ellipsis whitespace-nowrap text-[17px] font-bold ${overdue ? 'gu-overdue-title' : 'gu-text-primary'}`}>{formatSubscriptionName(subscription.name, subscription.tariff, t)}</p>
					{renewalBadgeLabel && <AccentBadge className={`shrink-0 rounded-full px-2 py-0.5 text-[11px] font-semibold ${overdue ? 'gu-overdue-badge' : ''}`}>{renewalBadgeLabel}</AccentBadge>}
				</div>
				<p className={`mt-0.5 overflow-hidden text-ellipsis whitespace-nowrap text-[13px] ${overdue ? 'gu-overdue-muted' : 'gu-text-muted'}`}>{getRenewOnDateLabel(subscription.date_pay, t, i18n.language)}</p>
			</div>

			<div className="shrink-0 text-right">
				<p className={`text-[17px] font-bold ${overdue ? 'gu-overdue-title' : ''}`}>{formatSubscriptionPrice(getSubscriptionShareAmount(subscription), subscription.currency, i18n.language)}</p>
				<p className={`text-[13px] capitalize ${overdue ? 'gu-overdue-muted' : 'gu-text-muted'}`}>{t(`home.period-${subscription.period}`)}</p>
			</div>
		</Link>
	);
};

export default TrackedSubscriptionCard;
