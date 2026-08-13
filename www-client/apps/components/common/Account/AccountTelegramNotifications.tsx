import BellOutline from '@/components/@icons/bell-outline';
import GUIButton from '@/components/ui/Button/GUIButton';
import PremiumGatedSection from '@/components/common/Account/PremiumGatedSection';
import { basicUserTelegramDisconnect, basicUserTelegramLink } from '@/rest/userAPI';
import { useTranslation } from 'react-i18next';

interface AccountTelegramNotificationsProps {
	isPremium: boolean;
	canUseNotification: boolean;
	connected: boolean;
	username?: string | null;
	onChanged: () => void | Promise<void>;
}

const AccountTelegramNotifications: React.FC<AccountTelegramNotificationsProps> = ({ isPremium, canUseNotification, connected, username, onChanged }) => {
	const { t } = useTranslation();

	const onConnect = async () => {
		if (!canUseNotification) {
			return;
		}

		const response = await basicUserTelegramLink();
		window.open(response.bot_url, '_blank', 'noopener,noreferrer');
	};

	const onDisconnect = async () => {
		await basicUserTelegramDisconnect();
		await onChanged();
	};

	return (
		<PremiumGatedSection title={t('account.notifications-settings')} isPremium={isPremium && canUseNotification}>
			<div className="space-y-3">
				{connected ? (
					<div className="flex items-center justify-between gap-3 py-3">
						<div className="flex min-w-0 items-center gap-3">
							<span className="inline-flex h-9 w-9 shrink-0 items-center justify-center gu-text-primary">
								<BellOutline fill="currentColor" size={21} />
							</span>
							<div className="min-w-0">
								<p className="text-[15px] font-medium gu-text-primary">{t('account.telegram-connected')}</p>
								<p className="truncate text-[13px] gu-text-muted">{username ? `@${username}` : t('account.telegram-connected-desc')}</p>
							</div>
						</div>
						<GUIButton type="button" onClick={onDisconnect} className="shrink-0 text-[13px]">
							{t('account.telegram-disconnect')}
						</GUIButton>
					</div>
				) : (
					<div className="space-y-2">
						<p className="text-[13px] gu-text-muted">{t('account.telegram-connect-desc')}</p>
						<GUIButton type="button" onClick={onConnect} className="w-full py-3">
							{t('account.telegram-connect')}
						</GUIButton>
					</div>
				)}
			</div>
		</PremiumGatedSection>
	);
};

export default AccountTelegramNotifications;
