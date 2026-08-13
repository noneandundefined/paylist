import GUISwitch from '@/components/ui/Checkbox/GUISwitch';
import ChevronRight from '@/components/@icons/chevron-right';

interface AccountSettingsRowProps {
	icon: React.ReactNode;
	label: string;
	value?: string;
	onClick?: () => void;
	disabled?: boolean;
	trailing?: 'chevron' | 'switch';
	switchChecked?: boolean;
	onSwitchChange?: (checked: boolean) => void;
	danger?: boolean;
}

const AccountSettingsRow: React.FC<AccountSettingsRowProps> = ({ icon, label, value, onClick, disabled = false, trailing = 'chevron', switchChecked, onSwitchChange, danger = false }) => {
	const isInteractive = trailing === 'chevron' && onClick;

	const content = (
		<>
			<span className="inline-flex h-9 w-9 shrink-0 items-center justify-center gu-text-primary">{icon}</span>
			<span className={`flex-1 text-left text-[15px] font-medium ${danger ? 'text-red-600 dark:text-red-400' : 'gu-text-primary'}`}>{label}</span>
			{value && <span className="text-[15px] gu-text-muted">{value}</span>}
			{trailing === 'chevron' && <ChevronRight fill="#94a3b8" size={20} />}
			{trailing === 'switch' && <GUISwitch checked={switchChecked ?? false} onChange={(event) => onSwitchChange?.(event.target.checked)} />}
		</>
	);

	if (isInteractive) {
		return (
			<button type="button" onClick={onClick} disabled={disabled} className="flex w-full items-center gap-3 px-4 py-3.5 text-left transition hover:bg-[var(--surface-muted)] disabled:cursor-not-allowed disabled:opacity-50">
				{content}
			</button>
		);
	}

	return <div className="flex items-center gap-3 px-4 py-3.5">{content}</div>;
};

export default AccountSettingsRow;
