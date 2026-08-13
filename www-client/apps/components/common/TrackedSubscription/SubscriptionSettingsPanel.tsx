import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { ROUTES } from '@/constants/constants';
import SubscriptionSettingRow from '@/components/common/TrackedSubscription/SubscriptionSettingRow';
import Poll from '@/components/@icons/poll';
import Reload from '@/components/@icons/reload';
import BellOutline from '@/components/@icons/bell-outline';

interface SubscriptionSettingsPanelProps {
	autoRenewal: boolean;
	notification: boolean;
	includeInAnalytics: boolean;
	onAutoRenewalChange: (checked: boolean) => void;
	onNotificationChange: (checked: boolean) => void;
	onIncludeInAnalyticsChange: (checked: boolean) => void;
	canUseNotification: boolean;
	onPremiumRequired: () => void;
}

const SubscriptionSettingsPanel: React.FC<SubscriptionSettingsPanelProps> = ({
	autoRenewal,
	notification,
	includeInAnalytics,
	onAutoRenewalChange,
	onNotificationChange,
	onIncludeInAnalyticsChange,
	canUseNotification,
	onPremiumRequired,
}) => {
	const { t } = useTranslation();

	const handleNotificationChange = (checked: boolean) => {
		if (checked && !canUseNotification) {
			onPremiumRequired();
			return;
		}

		onNotificationChange(checked);
	};

	return (
		<div className="gu-glass-card divide-y gu-divide overflow-hidden">
			<SubscriptionSettingRow icon={<Reload fill="currentColor" size={21} />} label={t('subscription.auto-renewal')} checked={autoRenewal} onChange={onAutoRenewalChange} />

			<SubscriptionSettingRow
				icon={<BellOutline fill="currentColor" size={21} />}
				label={t('subscription.notifications')}
				checked={notification}
				onChange={handleNotificationChange}
				canUse={canUseNotification}
				hint={
					canUseNotification ? (
						<>
							{t('subscription.notifications-account-hint')}{' '}
							<Link to={ROUTES.ACCOUNT} className="font-semibold text-[#0085FF] no-underline hover:no-underline">
								{t('subscription.notifications-account-link')}
							</Link>
						</>
					) : undefined
				}
			/>

			<SubscriptionSettingRow icon={<Poll fill="currentColor" size={21} />} label={t('subscription.include-in-analytics')} checked={includeInAnalytics} onChange={onIncludeInAnalyticsChange} canUse={true} />
		</div>
	);
};

export default SubscriptionSettingsPanel;
