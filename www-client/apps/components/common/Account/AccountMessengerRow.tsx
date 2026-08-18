import InformationOutline from '@/components/@icons/information-outline';

type MessengerStatus = 'connected' | 'disconnected' | 'connecting';

interface AccountMessengerRowProps {
	name: string;
	iconSrc: string;
	iconAlt: string;
	iconClassName?: string;
	status: MessengerStatus;
	statusLabel: string;
	showDivider?: boolean;
	disabled?: boolean;
	onClick: () => void;
}

const statusClassName: Record<MessengerStatus, string> = {
	connected: 'text-[#16a34a] dark:text-[#4ade80]',
	disconnected: 'gu-text-muted',
	connecting: 'gu-text-muted',
};

const AccountMessengerRow: React.FC<AccountMessengerRowProps> = ({ name, iconSrc, iconAlt, iconClassName = '', status, statusLabel, showDivider = false, disabled = false, onClick }) => {
	return (
		<button type="button" onClick={onClick} disabled={disabled} className="flex w-full items-center gap-4 text-left disabled:cursor-not-allowed">
			<img src={iconSrc} alt={iconAlt} draggable={false} className={`h-11 w-11 shrink-0 rounded-[6px] object-cover ${iconClassName}`} />

			<div className={`flex min-w-0 flex-1 items-center gap-3 py-3 ${showDivider ? 'border-b border-[var(--divider)]' : ''}`}>
				<div className="min-w-0 flex-1">
					<p className="text-[16px] font-semibold leading-tight gu-text-primary">{name}</p>
					<p className={`mt-0.5 truncate text-[13px] leading-snug ${statusClassName[status]}`}>{statusLabel}</p>
				</div>

				{status === 'connecting' ? (
					<span className="h-5 w-5 shrink-0 animate-spin rounded-full border-2 border-[var(--divider)] border-t-[var(--text-muted)]" aria-hidden />
				) : (
					<InformationOutline fill="var(--text-muted)" size={22} className="shrink-0" />
				)}
			</div>
		</button>
	);
};

export default AccountMessengerRow;
