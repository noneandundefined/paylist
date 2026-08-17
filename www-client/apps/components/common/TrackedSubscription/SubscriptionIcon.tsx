import { useState } from 'react';
import { getTrackedSubscriptionImageUrl } from '@/rest/trackedSubscriptionAPI';
import { getSubscriptionIconLabel } from '@/utils/TrackedSubscriptionDisplayUtils';
import RemoteImage from '@/components/ui/RemoteImage/RemoteImage';

interface SubscriptionIconProps {
	name: string;
	size?: 'sm' | 'md' | 'xs';
	className?: string;
}

const sizeClasses: Record<NonNullable<SubscriptionIconProps['size']>, string> = {
	xs: 'h-8 w-8 text-[10px]',
	sm: 'h-[3.2rem] w-[3.2rem]',
	md: 'h-[5rem] w-[5rem]',
};

const spinnerSizes: Record<NonNullable<SubscriptionIconProps['size']>, number> = {
	xs: 12,
	sm: 16,
	md: 20,
};

const SubscriptionIcon = ({ name, size = 'sm', className = '' }: SubscriptionIconProps) => {
	const [hasImageError, setHasImageError] = useState(false);
	const label = getSubscriptionIconLabel(name);
	const sizeClass = sizeClasses[size];

	if (!hasImageError) {
		return (
			<div className={`relative shrink-0 overflow-hidden rounded-xl ${sizeClass} ${className}`}>
				<RemoteImage src={getTrackedSubscriptionImageUrl(name)} alt="" className="h-full w-full object-cover" spinnerSize={spinnerSizes[size]} onError={() => setHasImageError(true)} />
			</div>
		);
	}

	return <div className={`${sizeClass} flex shrink-0 items-center justify-center rounded-xl bg-slate-400 font-bold text-white ${className}`}>{label}</div>;
};

export default SubscriptionIcon;
