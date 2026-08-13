import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { ROUTES, buildRoute } from '@/constants/constants';
import type { AnalyticsRecommendation } from '@/utils/analyticsUtils';

interface AnalyticsRecommendationsProps {
	items: AnalyticsRecommendation[];
}

const iconByType: Record<AnalyticsRecommendation['type'], string> = {
	'yearly-save': '💡',
	concentration: '⚠️',
	excluded: '📊',
	cluster: '📅',
	'small-subs': '🔍',
	'upcoming-heavy': '📈',
};

const AnalyticsRecommendations: React.FC<AnalyticsRecommendationsProps> = ({ items }) => {
	const { t } = useTranslation();

	if (items.length === 0) {
		return (
			<section className="gu-glass-card p-5">
				<h2 className="text-[15px] font-semibold gu-text-primary">{t('analytics.recommendations-title')}</h2>
				<p className="mt-2 text-sm gu-text-muted">{t('analytics.recommendations-empty')}</p>
			</section>
		);
	}

	return (
		<section className="gu-glass-card p-5">
			<h2 className="text-[15px] font-semibold gu-text-primary">{t('analytics.recommendations-title')}</h2>
			<p className="mt-1 text-[13px] gu-text-muted">{t('analytics.recommendations-subtitle')}</p>

			<ul className="mt-4 space-y-8">
				{items.map((item) => {
					const content = (
						<li key={item.id} className="">
							<div className="flex items-start gap-3">
								<span className="text-xl" aria-hidden>
									{iconByType[item.type]}
								</span>
								<div className="min-w-0 flex-1">
									<p className="text-[14px] font-semibold gu-text-primary">{t(item.titleKey)}</p>
									<p className="mt-1 text-[13px] leading-relaxed gu-text-secondary">{t(item.descKey, item.descValues)}</p>
								</div>
							</div>
						</li>
					);

					if (!item.subscriptionId) {
						return content;
					}

					return (
						<Link key={item.id} to={buildRoute(ROUTES.SUBSCRIPTION_DETAIL, { id: item.subscriptionId })} className="block no-underline hover:no-underline">
							{content}
						</Link>
					);
				})}
			</ul>
		</section>
	);
};

export default AnalyticsRecommendations;
