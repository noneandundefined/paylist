import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { basicUserAvatarUpdate } from '@/rest/userAPI';

type UserAvatarSize = 'sm' | 'lg';

interface UserAvatarProps {
	initials: string;
	isPremium?: boolean;
	size?: UserAvatarSize;
	src?: string | null;
	editable?: boolean;
	onUpdated?: () => void;
}

const premiumRingClass = 'border-[linear-gradient(98deg,#0085FF_0%,#22A1FF_45%,#60D1FF_100%)]';
const freeRingClass = 'border-slate-400';

const UserAvatar: React.FC<UserAvatarProps> = ({ initials, isPremium = false, size = 'lg', src, editable = false, onUpdated }) => {
	const { t } = useTranslation();

	const inputRef = useRef<HTMLInputElement>(null);
	const [broken, setBroken] = useState(false);
	const [uploading, setUploading] = useState(false);

	useEffect(() => {
		setBroken(false);
	}, [src]);

	const ringClass = isPremium ? premiumRingClass : freeRingClass;
	const isSmall = size === 'sm';
	const innerClass = isSmall ? (isPremium ? `${premiumRingClass} text-white` : 'bg-slate-400 text-white') : isPremium ? 'bg-[linear-gradient(180deg,#dbeafe_0%,#eff6ff_100%)] text-[#2a2867]' : 'bg-slate-200 text-slate-600';
	const showImage = Boolean(src) && !broken;

	const onFileChange: React.ChangeEventHandler<HTMLInputElement> = async (event) => {
		const file = event.target.files?.[0];
		event.target.value = '';

		if (!file || uploading) {
			return;
		}

		setUploading(true);

		try {
			await basicUserAvatarUpdate(file);
			onUpdated?.();
		} finally {
			setUploading(false);
		}
	};

	const avatar = (
		// <div className={`rounded-full`}>
		// 	<div className={`relative ${ringClass} overflow-hidden rounded-full ${isSmall ? 'h-11 w-11' : 'h-[6rem] w-[6rem]'}`}>
		// 		<div className={`flex h-full w-full items-center justify-center overflow-hidden font-semibold ${isSmall ? 'text-sm' : 'text-[20px]'} ${innerClass} ${showImage ? 'p-[2px]' : ''}`}>
		// 			{showImage ? <img src={src ?? undefined} alt="" className="h-full w-full rounded-full object-cover" onError={() => setBroken(true)} /> : initials}
		// 		</div>
		// 		{uploading && <div className="absolute inset-0 bg-black/35" />}
		// 	</div>
		// </div>
		<div className={`rounded-full border-2 p-1 ${ringClass}`}>
			<div className="bg-white rounded-full">
				{showImage ? (
					<div className={`${isSmall ? 'h-11 w-11' : 'h-[6rem] w-[6rem]'}`}>
						<img src={src ?? undefined} alt="" className="h-full w-full rounded-full object-cover" onError={() => setBroken(true)} />
					</div>
				) : (
					<div className={`relative ${ringClass} overflow-hidden rounded-full ${isSmall ? 'h-11 w-11' : 'h-[6rem] w-[6rem]'}`}>
						<div className={`flex h-full w-full items-center justify-center overflow-hidden font-semibold ${isSmall ? 'text-sm' : 'text-[20px]'} ${innerClass} ${showImage ? 'p-[2px]' : ''}`}>{initials}</div>
					</div>
				)}
				{uploading && <div className="absolute inset-0 bg-black/35" />}
			</div>
		</div>
	);

	if (!editable) {
		return avatar;
	}

	return (
		<>
			<button type="button" className="cursor-pointer rounded-full border-0 bg-transparent p-0" aria-label={t('account.change-avatar')} disabled={uploading} onClick={() => inputRef.current?.click()}>
				{avatar}
			</button>
			<input ref={inputRef} type="file" accept="image/jpeg,image/png,image/webp,image/gif" className="hidden" onChange={onFileChange} />
		</>
	);
};

export default UserAvatar;
