import PremiumBadgeMini from '@/components/common/PremiumBadge/PremiumBadgeMini';

interface PremiumGatedSectionProps {
	title: string;
	isPremium: boolean;
	hint?: string;
	children: React.ReactNode;
}

const PremiumGatedSection: React.FC<PremiumGatedSectionProps> = ({ title, isPremium, hint, children }) => {
	return (
		<div className="p-3 space-y-2">
			<div className="flex items-center justify-between gap-3">
				<h2 className="text-[15px] font-semibold gu-text-primary">{title}</h2>
				{!isPremium && <PremiumBadgeMini mobileView={true} />}
			</div>

			{isPremium ? children : hint ? <p className="text-[13px] gu-text-muted">{hint}</p> : null}
		</div>
	);
};

export default PremiumGatedSection;
