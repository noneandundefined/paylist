import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { formatSubscriptionPrice } from '@/utils/TrackedSubscriptionDisplayUtils';
import { getCategoryLabel } from '@/utils/categoryDisplayUtils';
import type { SubscriptionCategoryResponse } from '@/rest/trackedSubscriptionAPI';
import type { CategoryShare } from '@/utils/analyticsUtils';

interface AnalyticsCategoryDonutProps {
	items: CategoryShare[];
	currency: string;
	categories: SubscriptionCategoryResponse[];
}

const CHART_SIZE = 176;
const SEGMENT_GAP = 3;

const polarToCartesian = (center: number, radius: number, angleInDegrees: number) => {
	const angleInRadians = ((angleInDegrees - 90) * Math.PI) / 180;

	return {
		x: center + radius * Math.cos(angleInRadians),
		y: center + radius * Math.sin(angleInRadians),
	};
};

const describeDonutSegment = (center: number, innerRadius: number, outerRadius: number, startAngle: number, endAngle: number) => {
	const startOuter = polarToCartesian(center, outerRadius, startAngle);
	const endOuter = polarToCartesian(center, outerRadius, endAngle);
	const startInner = polarToCartesian(center, innerRadius, endAngle);
	const endInner = polarToCartesian(center, innerRadius, startAngle);
	const largeArcFlag = endAngle - startAngle <= 180 ? 0 : 1;

	return [
		`M ${startOuter.x} ${startOuter.y}`,
		`A ${outerRadius} ${outerRadius} 0 ${largeArcFlag} 1 ${endOuter.x} ${endOuter.y}`,
		`L ${startInner.x} ${startInner.y}`,
		`A ${innerRadius} ${innerRadius} 0 ${largeArcFlag} 0 ${endInner.x} ${endInner.y}`,
		'Z',
	].join(' ');
};

const buildDonutSegments = (items: CategoryShare[]) => {
	const center = CHART_SIZE / 2;
	const outerRadius = center - 6;
	const innerRadius = outerRadius * 0.58;
	const totalShare = items.reduce((sum, item) => sum + item.sharePercent, 0) || 100;
	const availableAngle = 360 - items.length * SEGMENT_GAP;

	let currentAngle = 0;

	return items.map((item) => {
		const sweep = (item.sharePercent / totalShare) * availableAngle;
		const startAngle = currentAngle;
		const endAngle = currentAngle + sweep;
		currentAngle = endAngle + SEGMENT_GAP;

		return {
			...item,
			path: describeDonutSegment(center, innerRadius, outerRadius, startAngle, endAngle),
		};
	});
};

const AnalyticsCategoryDonut: React.FC<AnalyticsCategoryDonutProps> = ({ items, currency, categories }) => {
	const { t, i18n } = useTranslation();

	const segments = useMemo(() => buildDonutSegments(items), [items]);

	const getLabel = (slug: string) => {
		if (slug === 'other') {
			return t('analytics.category-other');
		}

		if (slug === 'uncategorized') {
			return t('analytics.category-uncategorized');
		}

		const category = categories.find((item) => item.slug === slug);

		if (category) {
			return getCategoryLabel(category, t);
		}

		return t(`subscription.category-${slug}`, slug);
	};

	if (items.length === 0) {
		return null;
	}

	return (
		<section className="gu-glass-card p-5">
			<h2 className="text-[15px] font-semibold gu-text-primary">{t('analytics.category-breakdown-title')}</h2>
			<p className="mt-1 text-[13px] gu-text-muted">{t('analytics.category-breakdown-subtitle')}</p>

			<div className="mt-5 flex flex-col items-center gap-6 sm:flex-row sm:items-center">
				<div className="relative shrink-0" style={{ width: CHART_SIZE, height: CHART_SIZE }}>
					<svg width={CHART_SIZE} height={CHART_SIZE} viewBox={`0 0 ${CHART_SIZE} ${CHART_SIZE}`} role="img" aria-label={t('analytics.category-breakdown-title')}>
						{segments.map((segment) => (
							<path key={segment.slug} d={segment.path} fill={segment.color} />
						))}
					</svg>
				</div>

				<ul className="w-full min-w-0 flex-1 space-y-4">
					{items.map((item) => (
						<li key={item.slug} className="flex items-center gap-3">
							<span className="h-2.5 w-2.5 shrink-0 rounded-full" style={{ backgroundColor: item.color }} aria-hidden="true" />
							<span className="min-w-0 flex-1 truncate text-[14px] font-medium gu-text-primary">{getLabel(item.slug)}</span>
							<span className="shrink-0 text-[13px] font-semibold gu-text-primary">{formatSubscriptionPrice(item.monthlyAmount, currency, i18n.language)}</span>
							<span className="shrink-0 text-[12px] gu-text-muted">· {item.sharePercent}%</span>
						</li>
					))}
				</ul>
			</div>
		</section>
	);
};

export default AnalyticsCategoryDonut;
