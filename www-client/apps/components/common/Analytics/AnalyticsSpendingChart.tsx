import { useEffect, useLayoutEffect, useMemo, useRef, useState } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { COLORS } from '@/constants/designTokens';
import { ROUTES } from '@/constants/constants';
import PremiumBadgeMini from '@/components/common/PremiumBadge/PremiumBadgeMini';
import { formatSubscriptionPrice } from '@/utils/TrackedSubscriptionDisplayUtils';
import type { MonthlyProjectionPoint } from '@/utils/analyticsUtils';

interface AnalyticsSpendingChartProps {
	points: MonthlyProjectionPoint[];
	inflationPoints?: MonthlyProjectionPoint[];
	currency: string;
	periodLabel: string;
	showInflation?: boolean;
	inflationRate?: number;
	isPremium?: boolean;
	hasCountry?: boolean;
}

const CHART_HEIGHT = 160;
const CHART_WIDTH = 320;
const MIN_CHART_HEIGHT = 170;
const MAX_CHART_HEIGHT = 400;
const CHART_ASPECT = CHART_HEIGHT / CHART_WIDTH;
const PADDING_X = 8;
const PADDING_Y = 16;

const buildLineGeometry = (series: MonthlyProjectionPoint[], baseline: number, maxAmount: number, range: number, chartWidth: number, chartHeight: number) => {
	const innerWidth = chartWidth - PADDING_X * 2;
	const innerHeight = chartHeight - PADDING_Y * 2;

	const coordinates = series.map((point, index) => {
		const x = PADDING_X + (index / Math.max(series.length - 1, 1)) * innerWidth;
		const normalized = (point.amount - baseline) / (maxAmount - baseline + range * 0.05);
		const y = PADDING_Y + innerHeight - normalized * innerHeight;

		return { x, y, point };
	});

	const linePath = coordinates.map((item, index) => `${index === 0 ? 'M' : 'L'} ${item.x.toFixed(2)} ${item.y.toFixed(2)}`).join(' ');
	const areaPath = `${linePath} L ${coordinates[coordinates.length - 1]?.x.toFixed(2) ?? PADDING_X} ${chartHeight - PADDING_Y} L ${coordinates[0]?.x.toFixed(2) ?? PADDING_X} ${chartHeight - PADDING_Y} Z`;

	return { coordinates, linePath, areaPath };
};

const AnalyticsSpendingChart: React.FC<AnalyticsSpendingChartProps> = ({ points, inflationPoints, currency, periodLabel, showInflation = false, inflationRate, isPremium = false, hasCountry = false }) => {
	const { t, i18n } = useTranslation();
	const chartContainerRef = useRef<HTMLDivElement>(null);
	const [chartWidth, setChartWidth] = useState(CHART_WIDTH);
	const [activeIndex, setActiveIndex] = useState(() => Math.max(points.length - 1, 0));

	const chartHeight = Math.min(Math.max(chartWidth * CHART_ASPECT, MIN_CHART_HEIGHT), MAX_CHART_HEIGHT);

	useLayoutEffect(() => {
		const node = chartContainerRef.current;

		if (!node) {
			return;
		}

		const nextWidth = node.getBoundingClientRect().width;

		if (nextWidth > 0) {
			setChartWidth(nextWidth);
		}
	}, []);

	useEffect(() => {
		const node = chartContainerRef.current;

		if (!node) {
			return;
		}

		const updateWidth = () => {
			const nextWidth = node.getBoundingClientRect().width;

			if (nextWidth > 0) {
				setChartWidth(nextWidth);
			}
		};

		updateWidth();

		const observer = new ResizeObserver(updateWidth);
		observer.observe(node);

		return () => observer.disconnect();
	}, []);

	const chart = useMemo(() => {
		if (points.length === 0) {
			return null;
		}

		const inflationSeries = showInflation && inflationPoints?.length ? inflationPoints : [];
		const allAmounts = [...points.map((point) => point.amount), ...inflationSeries.map((point) => point.amount)];
		const maxAmount = Math.max(...allAmounts, 1);
		const minAmount = Math.min(...allAmounts);
		const range = Math.max(maxAmount - minAmount, maxAmount * 0.15, 1);
		const baseline = Math.max(minAmount - range * 0.1, 0);

		const innerHeight = chartHeight - PADDING_Y * 2;

		const baseGeometry = buildLineGeometry(points, baseline, maxAmount, range, chartWidth, chartHeight);
		const inflationGeometry = inflationSeries.length > 0 ? buildLineGeometry(inflationSeries, baseline, maxAmount, range, chartWidth, chartHeight) : null;

		const average = allAmounts.reduce((sum, value) => sum + value, 0) / allAmounts.length;
		const averageY = PADDING_Y + innerHeight - ((average - baseline) / (maxAmount - baseline + range * 0.05)) * innerHeight;

		return {
			baseGeometry,
			inflationGeometry,
			averageY,
			activeAmount: points[activeIndex]?.amount ?? 0,
			activeInflationAmount: inflationSeries[activeIndex]?.amount ?? 0,
		};
	}, [activeIndex, chartHeight, chartWidth, inflationPoints, points, showInflation]);

	if (!chart || points.length === 0) {
		return (
			<section className="gu-glass-card p-5">
				<p className="text-center text-sm gu-text-muted">{t('analytics.empty-chart')}</p>
			</section>
		);
	}

	const activePoint = chart.baseGeometry.coordinates[activeIndex];
	const activeInflationPoint = chart.inflationGeometry?.coordinates[activeIndex];

	return (
		<section className="gu-glass-card w-full min-w-0 p-5">
			<div className="text-center">
				<div className="flex items-center justify-center gap-2">
					<p className="text-[11px] font-semibold uppercase tracking-[0.08em] gu-text-secondary">{t('analytics.spending-forecast')}</p>
					{showInflation ? <PremiumBadgeMini mobileView={true} /> : null}
				</div>

				<div className="mt-2 flex flex-col items-center gap-1">
					<div className="flex items-end justify-center gap-2">
						<p className="text-[32px] font-bold leading-none gu-text-primary sm:text-[36px]">{formatSubscriptionPrice(chart.activeAmount, currency, i18n.language)}</p>
						<span className="mb-1 rounded-full bg-[#ecfccb] px-2 py-0.5 text-[12px] font-medium text-[#3f6212]">{periodLabel}</span>
					</div>

					{showInflation ? (
						<p className="text-[14px] font-medium text-[#d7ff00]">
							{t('analytics.inflation-adjusted-value', {
								amount: formatSubscriptionPrice(chart.activeInflationAmount, currency, i18n.language),
								rate: inflationRate ?? 0,
							})}
						</p>
					) : null}
				</div>
			</div>

			<div ref={chartContainerRef} className="relative mt-4 w-full min-w-0" style={{ height: chartHeight }}>
				<svg viewBox={`0 0 ${chartWidth} ${chartHeight}`} width="100%" height="100%" className="block h-full w-full touch-none" onMouseLeave={() => setActiveIndex(Math.max(points.length - 1, 0))}>
					<defs>
						<linearGradient id="analytics-area-gradient" x1="0" y1="0" x2="0" y2="1">
							<stop offset="0%" stopColor={COLORS.brandBlue} stopOpacity="0.35" />
							<stop offset="100%" stopColor={COLORS.brandBlue} stopOpacity="0" />
						</linearGradient>
					</defs>

					<line x1={PADDING_X} x2={chartWidth - PADDING_X} y1={chart.averageY} y2={chart.averageY} stroke="currentColor" strokeDasharray="4 4" className="gu-text-muted" strokeOpacity={0.35} />

					{activePoint ? (
						<>
							<rect x={activePoint.x - 18} y={PADDING_Y} width={36} height={chartHeight - PADDING_Y * 2} fill={COLORS.brandBlue} opacity={0.08} rx={8} />
							<line x1={activePoint.x} x2={activePoint.x} y1={PADDING_Y} y2={chartHeight - PADDING_Y} stroke="#ffffff" strokeOpacity={0.65} />
						</>
					) : null}

					<path d={chart.baseGeometry.areaPath} fill="url(#analytics-area-gradient)" />
					<path d={chart.baseGeometry.linePath} fill="none" stroke={COLORS.brandBlue} strokeWidth={2.5} strokeLinecap="round" strokeLinejoin="round" />

					{chart.inflationGeometry ? <path d={chart.inflationGeometry.linePath} fill="none" stroke={COLORS.accentLime} strokeWidth={2} strokeDasharray="6 4" strokeLinecap="round" strokeLinejoin="round" /> : null}

					{chart.baseGeometry.coordinates.map((item, index) => (
						<g key={item.point.monthKey}>
							<circle cx={item.x} cy={item.y} r={index === activeIndex ? 4.5 : 0} fill="#ffffff" />
							{activeInflationPoint && index === activeIndex ? <circle cx={activeInflationPoint.x} cy={activeInflationPoint.y} r={4} fill={COLORS.accentLime} /> : null}
							<rect x={item.x - 24} y={0} width={48} height={chartHeight} fill="transparent" onMouseEnter={() => setActiveIndex(index)} onFocus={() => setActiveIndex(index)} onClick={() => setActiveIndex(index)} />
						</g>
					))}
				</svg>
			</div>

			{showInflation ? (
				<div className="mt-3 flex flex-wrap items-center justify-center gap-4 text-[12px] gu-text-muted">
					<span className="inline-flex items-center gap-2">
						<span className="h-0.5 w-5 rounded-full bg-[#0085FF]" />
						{t('analytics.legend-current')}
					</span>
					<span className="inline-flex items-center gap-2">
						<span className="h-0.5 w-5 rounded-full border-t-2 border-dashed border-[#d7ff00]" />
						{t('analytics.legend-inflation')}
					</span>
				</div>
			) : null}

			<div className="mt-2 flex justify-between gap-2 px-1">
				{points.map((point, index) => (
					<button key={point.monthKey} type="button" onClick={() => setActiveIndex(index)} className={`flex-1 text-center text-[11px] font-medium ${index === activeIndex ? 'gu-text-primary' : 'gu-text-muted'}`}>
						{point.label}
					</button>
				))}
			</div>

			{!isPremium ? (
				<p className="mt-4 text-center text-[13px] gu-text-muted">
					{t('analytics.inflation-premium-hint')}{' '}
					<Link to={ROUTES.PLANS} className="font-semibold text-[#0085FF] no-underline hover:no-underline">
						{t('analytics.inflation-premium-link')}
					</Link>
				</p>
			) : null}

			{isPremium && !hasCountry ? <p className="mt-4 text-center text-[13px] gu-text-muted">{t('analytics.inflation-country-hint')}</p> : null}

			{showInflation ? <p className="mt-3 text-center text-[11px] leading-relaxed gu-text-muted">{t('analytics.inflation-disclaimer')}</p> : null}
		</section>
	);
};

export default AnalyticsSpendingChart;
