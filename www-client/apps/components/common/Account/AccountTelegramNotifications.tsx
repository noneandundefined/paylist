import { useTranslation } from 'react-i18next';
import { useModalContext } from '@/context/useModalContext';
import { basicUserTelegramDisconnect, basicUserTelegramLink } from '@/rest/userAPI';
import AccountMessengerRow from '@/components/common/Account/AccountMessengerRow';
import { openMessengerConnectModal } from '@/components/common/Account/AccountMessengerConnectModal';

const TELEGRAM_ICON = '/local/images/social-network/tg-ico.png';

interface AccountTelegramNotificationsProps {
	canUseNotification: boolean;
	connected: boolean;
	username?: string | null;
	showDivider?: boolean;
	onChanged: () => void | Promise<void>;
}

const AccountTelegramNotifications: React.FC<AccountTelegramNotificationsProps> = ({ canUseNotification, connected, username, showDivider = false, onChanged }) => {
	const { t } = useTranslation();
	const { open } = useModalContext();

	const status = connected ? 'connected' : 'disconnected';
	const statusLabel = connected ? (username ? `@${username}` : t('account.messenger-connected-status')) : t('account.messenger-not-connected');

	const onClick = () => {
		if (!canUseNotification) {
			return;
		}

		openMessengerConnectModal(open, {
			config: {
				name: t('account.telegram-name'),
				iconSrc: TELEGRAM_ICON,
				host: 'telegram.org',
				accent: '#2AABEE',
			},
			connected,
			username,
			onConnect: async () => {
				const response = await basicUserTelegramLink();
				return response.bot_url;
			},
			onDisconnect: async () => {
				await basicUserTelegramDisconnect();
				await onChanged();
			},
		});
	};

	return (
		<AccountMessengerRow
			name={t('account.telegram-name')}
			iconSrc={TELEGRAM_ICON}
			iconAlt={t('account.telegram-name')}
			status={status}
			statusLabel={statusLabel}
			showDivider={showDivider}
			disabled={!canUseNotification}
			onClick={onClick}
		/>
	);
};

export default AccountTelegramNotifications;
