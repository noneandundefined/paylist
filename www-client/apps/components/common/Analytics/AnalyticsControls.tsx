import type { AnalyticsPeriod } from '@/utils/analyticsUtils';

interface AnalyticsControlsProps {
	period: AnalyticsPeriod;
	onPeriodChange: (period: AnalyticsPeriod) => void;
}

const periods: AnalyticsPeriod[] = ['1M', '3M', '6M', '1Y'];

const AnalyticsControls: React.FC<AnalyticsControlsProps> = ({ period, onPeriodChange }) => {
	return (
		<div className="flex justify-center gap-2">
			{periods.map((item) => (
				<button key={item} type="button" onClick={() => onPeriodChange(item)} className={`rounded-full px-4 py-1.5 text-[13px] font-semibold transition ${period === item ? 'bg-[#0085FF] text-white' : 'gu-text-muted'}`}>
					{item}
				</button>
			))}
		</div>
	);
};

export default AnalyticsControls;
