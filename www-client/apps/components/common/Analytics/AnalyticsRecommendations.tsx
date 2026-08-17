import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { ROUTES, buildRoute } from '@/constants/constants';
import type { AnalyticsRecommendation } from '@/rest/trackedSubscriptionAPI';

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
	overlap: '🎯',
	'family-share': '👨‍👩‍👧',
	'crowd-overpay': '📉',
	downgrade: '⬇️',
	'expensive-tariff': '💎',
};

const AnalyticsRecommendations: React.FC<AnalyticsRecommendationsProps> = ({ items }) => {
	const { t } = useTranslation();

	const interpolate = (values?: Record<string, string | number>) => {
		if (!values) {
			return undefined;
		}

		const next: Record<string, string | number> = { ...values };

		if (typeof values.tariff === 'string') {
			next.tariff = t(`subscription.tariff-${values.tariff}`);
		}

		if (typeof values.cheaper_tariff === 'string') {
			next.cheaper_tariff = t(`subscription.tariff-${values.cheaper_tariff}`);
		}

		if (typeof values.country === 'string') {
			next.country = t(`country.${values.country}`, { defaultValue: values.country });
		}

		if (typeof values.category === 'string') {
			next.category = t(`subscription.category-${values.category}`, { defaultValue: values.category });
		}

		return next;
	};

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
									{iconByType[item.type] ?? '💡'}
								</span>
								<div className="min-w-0 flex-1">
									<p className="text-[14px] font-semibold gu-text-primary">{t(item.title_key)}</p>
									<p className="mt-1 text-[13px] leading-relaxed gu-text-secondary">{t(item.desc_key, interpolate(item.desc_values))}</p>
								</div>
							</div>
						</li>
					);

					if (!item.subscription_id) {
						return content;
					}

					return (
						<Link key={item.id} to={buildRoute(ROUTES.SUBSCRIPTION_DETAIL, { id: item.subscription_id })} className="block no-underline hover:no-underline">
							{content}
						</Link>
					);
				})}
			</ul>
		</section>
	);
};

export default AnalyticsRecommendations;
