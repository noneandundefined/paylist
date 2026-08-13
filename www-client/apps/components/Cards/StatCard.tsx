import type { ReactNode } from 'react';

interface StatCardProps {
	icon: ReactNode;
	title: string;
	desc: string;
}

const StatCard: React.FC<StatCardProps> = ({ icon, title, desc }) => {
	return (
		<div className="gu-glass-card flex min-h-[110px] flex-1 flex-col justify-between p-4 space-y-2">
			{icon}

			<div>
				<p className="text-[17px] font-semibold">{title}</p>
				<p className="text-[13px]">{desc}</p>
			</div>
		</div>
	);
};

export default StatCard;
