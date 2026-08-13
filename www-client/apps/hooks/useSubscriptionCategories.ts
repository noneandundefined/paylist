import { QUERY_KEYS } from '@/constants/QueryKeys.constants';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { basicSubscriptionCategoryList } from '@/rest/trackedSubscriptionAPI';

export const useSubscriptionCategories = () => {
	const { data, loading, reload } = useHandleServer([QUERY_KEYS.trackedSubscriptionCategoryList], () => basicSubscriptionCategoryList());

	return {
		categories: data ?? [],
		loading,
		reload,
	};
};
