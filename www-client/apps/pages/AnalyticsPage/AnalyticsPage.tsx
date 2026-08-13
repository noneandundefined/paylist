import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import PageLayout from '@/pages/PageLayout';
import Header from '@/components/Header/Header';
import PageLoadingState from '@/components/common/PageLoadingState/PageLoadingState';
import AnalyticsControls from '@/components/common/Analytics/AnalyticsControls';
import AnalyticsStatRow from '@/components/common/Analytics/AnalyticsStatRow';
import AnalyticsSpendingChart from '@/components/common/Analytics/AnalyticsSpendingChart';
import AnalyticsBreakdown from '@/components/common/Analytics/AnalyticsBreakdown';
import AnalyticsCategoryDonut from '@/components/common/Analytics/AnalyticsCategoryDonut';
import AnalyticsRecommendations from '@/components/common/Analytics/AnalyticsRecommendations';
import { useSubscriptionCategories } from '@/hooks/useSubscriptionCategories';
import { QUERY_KEYS } from '@/constants/QueryKeys.constants';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { useLoginState } from '@/hooks/useLoginState';
import { basicCurrencyRates } from '@/rest/currencyAPI';
import { fetchCountries, getCountryInflationRate } from '@/rest/countryAPI';
import { basicTrackedSubscriptionList, basicTrackedSubscriptionSummary } from '@/rest/trackedSubscriptionAPI';
import { basicUserSettingsGet } from '@/rest/userAPI';
import { buildAnalyticsSnapshot, buildInflationProjections, type AnalyticsPeriod } from '@/utils/analyticsUtils';

const AnalyticsPage = () => {
	const { t, i18n } = useTranslation();
	const { isPremium } = useLoginState();
	const { categories } = useSubscriptionCategories();

	const [period, setPeriod] = useState<AnalyticsPeriod>('6M');

	const { data: subscriptionsData, loading: subscriptionsLoading } = useHandleServer([QUERY_KEYS.trackedSubscriptionList], () => basicTrackedSubscriptionList());
	const { data: summaryData, loading: summaryLoading } = useHandleServer([QUERY_KEYS.trackedSubscriptionSummary], () => basicTrackedSubscriptionSummary());

	const { data: userSettings, loading: settingsLoading } = useHandleServer([QUERY_KEYS.userSettings], () => basicUserSettingsGet(), {
		enabled: isPremium,
	});

	const { data: countriesData } = useHandleServer([QUERY_KEYS.countryList], () => fetchCountries());

	const displayCurrency = summaryData?.display_currency ?? 'USD';
	const userCountry = userSettings?.country ?? null;

	const uniqueCurrencies = useMemo(() => {
		const codes = new Set((subscriptionsData ?? []).map((item) => item.currency.toUpperCase()));
		codes.delete(displayCurrency.toUpperCase());
		return Array.from(codes);
	}, [displayCurrency, subscriptionsData]);

	const { data: ratesData, loading: ratesLoading } = useHandleServer(
		[QUERY_KEYS.currencyRates, displayCurrency, uniqueCurrencies.join(',')],
		() => basicCurrencyRates(displayCurrency, uniqueCurrencies.length > 0 ? uniqueCurrencies : undefined),
		{ enabled: Boolean(summaryData) }
	);

	const snapshot = useMemo(() => {
		if (!summaryData) {
			return null;
		}

		const rates = ratesData?.rates ?? {};
		const ratesBase = ratesData?.base ?? displayCurrency;

		return buildAnalyticsSnapshot(subscriptionsData ?? [], displayCurrency, rates, ratesBase, period, i18n.language);
	}, [displayCurrency, i18n.language, period, ratesData, subscriptionsData, summaryData]);

	const inflationRate = useMemo(() => {
		if (!isPremium || !userCountry || !countriesData) {
			return null;
		}

		return getCountryInflationRate(countriesData, userCountry);
	}, [countriesData, isPremium, userCountry]);

	const inflationPoints = useMemo(() => {
		if (!snapshot || inflationRate === null) {
			return undefined;
		}

		return buildInflationProjections(snapshot.projections, inflationRate);
	}, [inflationRate, snapshot]);

	const showInflation = Boolean(isPremium && userCountry && inflationPoints && inflationRate !== null);

	const loading = subscriptionsLoading || summaryLoading || (isPremium && settingsLoading) || (Boolean(summaryData) && ratesLoading && uniqueCurrencies.length > 0);

	if (loading || !snapshot) {
		return <PageLoadingState />;
	}

	const activePoint = snapshot.projections[snapshot.projections.length - 1];

	return (
		<PageLayout>
			<Header />

			<AnalyticsControls period={period} onPeriodChange={setPeriod} />

			<div className="mt-4 space-y-4">
				<AnalyticsStatRow snapshot={snapshot} />

				<AnalyticsSpendingChart
					points={snapshot.projections}
					inflationPoints={inflationPoints}
					currency={snapshot.displayCurrency}
					periodLabel={activePoint ? `${activePoint.label}` : period}
					showInflation={showInflation}
					inflationRate={inflationRate ?? undefined}
					isPremium={isPremium}
					hasCountry={Boolean(userCountry)}
				/>

				<AnalyticsCategoryDonut items={snapshot.categoryBreakdown} currency={snapshot.displayCurrency} categories={categories} />

				<AnalyticsBreakdown items={snapshot.topSubscriptions} currency={snapshot.displayCurrency} />

				<AnalyticsRecommendations items={snapshot.recommendations} />

				<section className="gu-glass-card p-5">
					<h2 className="text-[15px] font-semibold gu-text-primary">{t('analytics.insights-title')}</h2>
					<div className="mt-4 grid grid-cols-2 gap-3 text-center sm:grid-cols-4">
						<div>
							<p className="text-[11px] gu-text-muted">{t('analytics.insight-monthly-plans')}</p>
							<p className="mt-1 text-[20px] font-bold gu-text-primary">{snapshot.monthlyCount}</p>
						</div>
						<div>
							<p className="text-[11px] gu-text-muted">{t('analytics.insight-yearly-plans')}</p>
							<p className="mt-1 text-[20px] font-bold gu-text-primary">{snapshot.yearlyCount}</p>
						</div>
						<div>
							<p className="text-[11px] gu-text-muted">{t('analytics.insight-excluded')}</p>
							<p className="mt-1 text-[20px] font-bold gu-text-primary">{snapshot.excludedCount}</p>
						</div>
						<div>
							<p className="text-[11px] gu-text-muted">{t('analytics.insight-per-day')}</p>
							<p className="mt-1 text-[20px] font-bold gu-text-primary">{snapshot.dailyAverage}</p>
						</div>
					</div>
				</section>
			</div>
		</PageLayout>
	);
};

export default AnalyticsPage;
