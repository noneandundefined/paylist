import { useEffect, useState } from 'react';
import Modal from '@/components/Modal/Modal';
import { useTranslation } from 'react-i18next';
import { useLoginState } from '@/hooks/useLoginState';
import { useModalContext } from '@/context/useModalContext';
import RemoteImage from '@/components/ui/RemoteImage/RemoteImage';
import { isImageFailed } from '@/utils/imageCacheUtils';
import InformationOutline from '@/components/@icons/information-outline';

export interface MessengerConnectConfig {
	name: string;
	iconSrc: string;
	iconClassName?: string;
	host: string;
	accent: string;
}

interface AccountMessengerConnectModalProps {
	config: MessengerConnectConfig;
	connected: boolean;
	username?: string | null;
	onConnect: () => Promise<string>;
	onDisconnect: () => Promise<void>;
}

const AccountMessengerConnectModal: React.FC<AccountMessengerConnectModalProps> = ({ config, connected, username, onConnect, onDisconnect }) => {
	const { t } = useTranslation();
	const { close } = useModalContext();
	const { displayName, initials, avatar } = useLoginState();
	const [busy, setBusy] = useState(false);
	const [avatarBroken, setAvatarBroken] = useState(() => isImageFailed(avatar));

	useEffect(() => {
		setAvatarBroken(isImageFailed(avatar));
	}, [avatar]);

	const bannerLabel = connected ? t('account.messenger-connected-banner') : t('account.messenger-connection-request');
	const primaryLabel = connected ? t('account.messenger-disconnect') : t('account.connect');
	const accentSoft = `color-mix(in srgb, ${config.accent} 16%, transparent)`;
	const showAvatar = Boolean(avatar) && !avatarBroken;

	const onPrimary = async () => {
		if (busy) {
			return;
		}

		setBusy(true);

		try {
			if (connected) {
				await onDisconnect();
				close();
				return;
			}

			const url = await onConnect();
			if (!url) {
				return;
			}

			window.open(url, '_blank', 'noopener,noreferrer');
			close();
		} finally {
			setBusy(false);
		}
	};

	return (
		<div className="flex flex-col">
			<div className="flex flex-col items-center text-center">
				<img src={config.iconSrc} alt="" draggable={false} className={`h-16 w-16 rounded-[18px] object-cover ${config.iconClassName ?? ''}`} />

				<div className="mt-3 flex items-center gap-1.5">
					<h2 className="text-[22px] font-bold leading-none gu-text-primary">{config.name}</h2>
				</div>

				<p className="mt-1.5 text-[14px] gu-text-muted">{config.host}</p>

				<div className="mt-5 inline-flex items-center gap-2 rounded-full px-3.5 py-2" style={{ backgroundColor: accentSoft, color: config.accent }}>
					<span className="text-[15px] font-medium leading-none" aria-hidden>
						∞
					</span>
					<span className="text-[14px] font-medium">{bannerLabel}</span>
					<InformationOutline fill="currentColor" size={16} />
				</div>
			</div>

			<div className="mt-8 flex items-start justify-between gap-4">
				<div className="min-w-0">
					<p className="text-[12px] gu-text-muted">{t('account.messenger-account')}</p>
					<div className="mt-2 flex min-w-0 items-center gap-2">
						<div className="relative h-8 w-8 shrink-0 overflow-hidden rounded-full bg-[var(--surface-muted)]">
							{showAvatar ? (
								<RemoteImage src={avatar ?? undefined} alt="" className="h-full w-full object-cover" spinnerSize={12} onError={() => setAvatarBroken(true)} />
							) : (
								<span className="flex h-full w-full items-center justify-center text-[11px] font-semibold gu-text-primary">{initials}</span>
							)}
						</div>
						<p className="truncate text-[15px] font-semibold gu-text-primary">{displayName}</p>
					</div>
				</div>
			</div>

			{connected && username ? <p className="mt-4 truncate text-center text-[13px] gu-text-muted">@{username}</p> : null}

			<div className="mt-8 flex flex-1">
				<button
					type="button"
					disabled={busy}
					onClick={() => void onPrimary()}
					className={`flex flex-1 items-center justify-center text-center rounded-[16px] py-3.5 text-[15px] font-semibold text-white transition hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-50 ${connected ? 'bg-red-500' : ''}`}
					style={connected ? undefined : { backgroundColor: config.accent }}
				>
					{busy ? t('action.loading') : primaryLabel}
				</button>
			</div>
		</div>
	);
};

export const openMessengerConnectModal = (open: (node: React.ReactNode) => void, props: AccountMessengerConnectModalProps): void => {
	open(
		<Modal title="Подключение">
			<AccountMessengerConnectModal {...props} />
		</Modal>
	);
};

export default AccountMessengerConnectModal;
