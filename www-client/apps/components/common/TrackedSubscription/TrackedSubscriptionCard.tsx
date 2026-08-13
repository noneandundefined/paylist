import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { ROUTES, buildRoute } from '@/constants/constants';
import AccentBadge from '@/components/common/AccentBadge/AccentBadge';
import SubscriptionIcon from '@/components/common/TrackedSubscription/SubscriptionIcon';
import { getRenewalBadgeDays, isSubscriptionOverdue } from '@/utils/SubscriptionRenewalUtils';
import type { TrackedSubscriptionResponse } from '@/rest/trackedSubscriptionAPI';
import { formatSubscriptionPrice, getRenewOnDateLabel } from '@/utils/TrackedSubscriptionDisplayUtils';

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
			className={`gu-glass-card flex w-full items-center gap-3 px-4 py-5 text-left no-underline transition hover:no-underline ${overdue ? 'gu-overdue-surface' : 'hover:bg-[var(--surface-muted)] gu-text-primary'}`}
		>
			<SubscriptionIcon name={subscription.name} />

			<div className="min-w-0 flex-1">
				<div className="flex min-w-0 items-center gap-2">
					<p className={`truncate text-[17px] -tracking-wide font-bold ${overdue ? 'gu-overdue-title' : 'gu-text-primary'}`}>{subscription.name}</p>
					{renewalBadgeLabel && <AccentBadge className={`shrink-0 rounded-full px-2 py-0.5 text-[11px] font-semibold ${overdue ? 'gu-overdue-badge' : ''}`}>{renewalBadgeLabel}</AccentBadge>}
				</div>
				<p className={`mt-0.5 truncate text-[13px] ${overdue ? 'gu-overdue-muted' : 'gu-text-muted'}`}>{getRenewOnDateLabel(subscription.date_pay, t, i18n.language)}</p>
			</div>

			<div className="shrink-0 text-right">
				<p className={`text-[17px] font-bold ${overdue ? 'gu-overdue-title' : ''}`}>{formatSubscriptionPrice(subscription.price, subscription.currency, i18n.language)}</p>
				<p className={`text-[13px] capitalize ${overdue ? 'gu-overdue-muted' : 'gu-text-muted'}`}>{t(`home.period-${subscription.period}`)}</p>
			</div>
		</Link>
	);
};

export default TrackedSubscriptionCard;
