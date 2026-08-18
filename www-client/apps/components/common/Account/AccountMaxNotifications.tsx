import { useTranslation } from 'react-i18next';
import { useModalContext } from '@/context/useModalContext';
import { basicUserMaxDisconnect, basicUserMaxLink } from '@/rest/userAPI';
import AccountMessengerRow from '@/components/common/Account/AccountMessengerRow';
import { openMessengerConnectModal } from '@/components/common/Account/AccountMessengerConnectModal';

const MAX_ICON = '/local/images/social-network/max-ico.png';

interface AccountMaxNotificationsProps {
	canUseNotification: boolean;
	connected: boolean;
	username?: string | null;
	onChanged: () => void | Promise<void>;
}

const AccountMaxNotifications: React.FC<AccountMaxNotificationsProps> = ({ canUseNotification, connected, username, onChanged }) => {
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
				name: t('account.max-name'),
				iconSrc: MAX_ICON,
				host: 'max.ru',
				accent: '#7C3AED',
			},
			connected,
			username,
			onConnect: async () => {
				const response = await basicUserMaxLink();
				return response.bot_url;
			},
			onDisconnect: async () => {
				await basicUserMaxDisconnect();
				await onChanged();
			},
		});
	};

	return <AccountMessengerRow name={t('account.max-name')} iconSrc={MAX_ICON} iconAlt={t('account.max-name')} status={status} statusLabel={statusLabel} disabled={!canUseNotification} onClick={onClick} />;
};

export default AccountMaxNotifications;
