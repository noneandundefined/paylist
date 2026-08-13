import PageLayout from '../PageLayout';
import Header from '@/components/Header/Header';
import PageLoadingState from '@/components/common/PageLoadingState/PageLoadingState';
import TotalSpendingCard from '@/components/Cards/TotalSpendingCard';
import TrackedSubscriptionList from '@/components/common/TrackedSubscription/TrackedSubscriptionList';
import { QUERY_KEYS } from '@/constants/QueryKeys.constants';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { basicTrackedSubscriptionList, basicTrackedSubscriptionSummary } from '@/rest/trackedSubscriptionAPI';

const HomePage = () => {
	const { data: subscriptionsData, loading: subscriptionsLoading } = useHandleServer([QUERY_KEYS.trackedSubscriptionList], () => basicTrackedSubscriptionList());
	const { data: summaryData, loading: summaryLoading } = useHandleServer([QUERY_KEYS.trackedSubscriptionSummary], () => basicTrackedSubscriptionSummary());

	const subscriptions = subscriptionsData ?? [];
	const loading = subscriptionsLoading || summaryLoading;

	if (loading || !summaryData) {
		return <PageLoadingState />;
	}

	return (
		<PageLayout>
			<Header />

			<TotalSpendingCard summary={summaryData} subscriptions={subscriptions} />

			<TrackedSubscriptionList subscriptions={subscriptions} />
		</PageLayout>
	);
};

export default HomePage;
