import type { TrackedSubscriptionCreateRequest } from './trackedSubscriptionCreateRequest.interface';

export type TrackedSubscriptionEditRequest = TrackedSubscriptionCreateRequest & {
	note?: string | null;
};
