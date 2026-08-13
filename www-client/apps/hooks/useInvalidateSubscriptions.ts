import { useQueryClient } from '@tanstack/react-query';
import { QUERY_KEYS } from '@/constants/QueryKeys.constants';

export const useInvalidateSubscriptions = () => {
	const queryClient = useQueryClient();

	const invalidateListAndSummary = async () => {
		await queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.trackedSubscriptionList] });
		await queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.trackedSubscriptionSummary] });
	};

	const invalidateAfterUpdate = async (subscriptionId: number) => {
		await queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.trackedSubscriptionDetail, subscriptionId] });
		await invalidateListAndSummary();
	};

	const invalidateAfterDelete = async (subscriptionId: number) => {
		await invalidateListAndSummary();
		queryClient.removeQueries({ queryKey: [QUERY_KEYS.trackedSubscriptionDetail, subscriptionId] });
	};

	return {
		invalidateListAndSummary,
		invalidateAfterUpdate,
		invalidateAfterDelete,
	};
};
