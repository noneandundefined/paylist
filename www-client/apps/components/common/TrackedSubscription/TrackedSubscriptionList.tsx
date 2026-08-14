import { useTranslation } from 'react-i18next';
import TrackedSubscriptionCard from './TrackedSubscriptionCard';
import SectionHeader from '@/components/common/SectionHeader/SectionHeader';
import type { TrackedSubscriptionResponse } from '@/rest/trackedSubscriptionAPI';

interface TrackedSubscriptionListProps {
	subscriptions: TrackedSubscriptionResponse[];
}

const TrackedSubscriptionList: React.FC<TrackedSubscriptionListProps> = ({ subscriptions }) => {
	const { t } = useTranslation();

	return (
		<section className="min-w-0 space-y-6">
			<SectionHeader title={t('home.active-subscriptions')} />

			<div className="flex min-w-0 flex-col gap-2.5">
				{subscriptions.length === 0 ? (
					<p className="px-1 text-[14px] text-center gu-text-muted">{t('home.no-active-subscriptions')}</p>
				) : (
					subscriptions.map((item) => <TrackedSubscriptionCard key={item.id} subscription={item} />)
				)}
			</div>
		</section>
	);
};

export default TrackedSubscriptionList;
