import PremiumBadgeMini from '@/components/common/PremiumBadge/PremiumBadgeMini';
import GUISwitch from '@/components/ui/Checkbox/GUISwitch';

interface SubscriptionSettingRowProps {
	icon: React.ReactNode;
	label: string;
	checked: boolean;
	onChange: (checked: boolean) => void;
	canUse?: boolean;
	locked?: boolean;
	showSwitch?: boolean;
	hint?: React.ReactNode;
}

const SubscriptionSettingRow: React.FC<SubscriptionSettingRowProps> = ({ icon, label, checked, onChange, canUse = true, locked = false, showSwitch = true, hint }) => {
	const enabled = canUse && !locked;

	return (
		<div>
			<div className="flex items-center justify-between gap-2 px-4 py-3.5">
				<div className="flex items-center gap-3">
					<span className="inline-flex h-9 w-9 shrink-0 items-center justify-center">{icon}</span>
					<span className="text-[15px] font-medium gu-text-primary">{label}</span>
					{!canUse && !locked && <PremiumBadgeMini mobileView={true} />}
				</div>

				{showSwitch && <GUISwitch checked={enabled ? checked : locked ? checked : false} disabled={!enabled} onChange={(event) => onChange(event.target.checked)} />}
			</div>

			{hint ? <p className="px-4 pb-3.5 text-[13px] leading-relaxed gu-text-muted">{hint}</p> : null}
		</div>
	);
};

export default SubscriptionSettingRow;
