import { useEffect, useState } from 'react';
import { getTrackedSubscriptionImageUrl } from '@/rest/trackedSubscriptionAPI';
import { isImageFailed } from '@/utils/imageCacheUtils';
import { useTheme } from '@/context/ThemeContext';
import RemoteImage from '@/components/ui/RemoteImage/RemoteImage';
import SubscriptionFallbackGlyph from './SubscriptionFallbackGlyph';

interface SubscriptionIconProps {
	name: string;
	categories?: string[];
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

const glyphSizes: Record<NonNullable<SubscriptionIconProps['size']>, number> = {
	xs: 18,
	sm: 26,
	md: 36,
};

const SubscriptionIcon = ({ name, categories, size = 'sm', className = '' }: SubscriptionIconProps) => {
	const { isDark } = useTheme();
	const imageUrl = getTrackedSubscriptionImageUrl(name);
	const [hasImageError, setHasImageError] = useState(() => isImageFailed(imageUrl));
	const sizeClass = sizeClasses[size];
	const fill = isDark ? '#f1f5f9' : '#0f172a';

	useEffect(() => {
		setHasImageError(isImageFailed(imageUrl));
	}, [imageUrl]);

	if (!hasImageError) {
		return (
			<div className={`relative shrink-0 overflow-hidden rounded-xl ${sizeClass} ${className}`}>
				<RemoteImage src={imageUrl} alt="" className="h-full w-full object-cover" spinnerSize={spinnerSizes[size]} onError={() => setHasImageError(true)} />
			</div>
		);
	}

	return (
		<div className={`${sizeClass} flex shrink-0 items-center justify-center rounded-xl bg-[var(--surface-muted)] ${className}`}>
			<SubscriptionFallbackGlyph name={name} categories={categories} size={glyphSizes[size]} fill={fill} />
		</div>
	);
};

export default SubscriptionIcon;
