import ChevronRight from '@/components/@icons/chevron-right';
import MessageArrowRightOutline from '@/components/@icons/message-arrow-right-outline';
import GUIButton from '@/components/ui/Button/GUIButton';
import { basicUserMaxDisconnect, basicUserMaxLink } from '@/rest/userAPI';
import { useTranslation } from 'react-i18next';

interface AccountMaxNotificationsProps {
	canUseNotification: boolean;
	connected: boolean;
	username?: string | null;
	onChanged: () => void | Promise<void>;
}

const AccountMaxNotifications: React.FC<AccountMaxNotificationsProps> = ({ canUseNotification, connected, username, onChanged }) => {
	const { t } = useTranslation();

	const onConnect = async () => {
		if (!canUseNotification) {
			return;
		}

		const response = await basicUserMaxLink();
		window.open(response.bot_url, '_blank', 'noopener,noreferrer');
	};

	const onDisconnect = async () => {
		await basicUserMaxDisconnect();
		await onChanged();
	};

	if (connected) {
		return (
			<div className="flex items-center justify-between gap-3 py-3">
				<div className="flex min-w-0 items-center gap-3">
					<span className="inline-flex h-9 w-9 shrink-0 items-center justify-center gu-text-primary">
						<MessageArrowRightOutline fill="currentColor" size={21} />
					</span>
					<div className="min-w-0">
						<p className="text-[15px] font-medium gu-text-primary">{t('account.max-connected')}</p>
						<p className="truncate text-[13px] gu-text-muted">{username ? `@${username}` : t('account.max-connected-desc')}</p>
					</div>
				</div>
				<GUIButton type="button" onClick={onDisconnect} className="shrink-0 text-[13px]">
					{t('account.max-disconnect')}
				</GUIButton>
			</div>
		);
	}

	return (
		<div className="space-y-2">
			{/* <p className="text-[13px] gu-text-muted">{t('account.max-connect-desc')}</p> */}
			<div className="flex items-center justify-between">
				<p className="text-[15px] gu-text-primary">{t('account.max-title')}</p>

				<button type="button" className="flex cursor-pointer items-center gap-2 rounded-md bg-[#d7ff00] px-3 py-1" onClick={onConnect}>
					<p className="text-black">{t('account.connect')}</p>
					<ChevronRight fill="#000" size={19} />
				</button>
			</div>
		</div>
	);
};

export default AccountMaxNotifications;
