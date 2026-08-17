import { useEffect, useState } from 'react';
import RemoteImage from '@/components/ui/RemoteImage/RemoteImage';

interface MemberAvatarProps {
	initials: string;
	src?: string | null;
	size?: number;
	className?: string;
}

const MemberAvatar: React.FC<MemberAvatarProps> = ({ initials, src, size = 40, className = '' }) => {
	const [broken, setBroken] = useState(false);

	useEffect(() => {
		setBroken(false);
	}, [src]);

	const showImage = Boolean(src) && !broken;

	return (
		<div className={`relative shrink-0 overflow-hidden rounded-full bg-[var(--surface-muted)] ${className}`} style={{ width: size, height: size }}>
			{showImage ? (
				<RemoteImage src={src ?? undefined} alt="" className="h-full w-full object-cover" spinnerSize={Math.max(12, Math.round(size / 3))} referrerPolicy="no-referrer" onError={() => setBroken(true)} />
			) : (
				<div className="flex h-full w-full items-center justify-center text-[12px] font-semibold gu-text-primary">{initials}</div>
			)}
		</div>
	);
};

export default MemberAvatar;
