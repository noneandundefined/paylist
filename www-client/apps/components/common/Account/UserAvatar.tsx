import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { basicUserAvatarUpdate } from '@/rest/userAPI';
import CameraPlusOutline from '@/components/@icons/camera-plus-outline';
import Pencil from '@/components/@icons/pencil';
import Fallback from '@/components/Fallback/Fallback';
import ImageSpinner from '@/components/ui/ImageSpinner/ImageSpinner';
import RemoteImage from '@/components/ui/RemoteImage/RemoteImage';

type UserAvatarSize = 'sm' | 'lg';

interface UserAvatarProps {
	initials: string;
	isPremium?: boolean;
	size?: UserAvatarSize;
	src?: string | null;
	editable?: boolean;
	onUpdated?: () => void;
}

const premiumGradientClass = 'bg-[linear-gradient(98deg,#0085FF_0%,#22A1FF_45%,#60D1FF_100%)]';
const freeRingClass = 'bg-slate-400';

const UserAvatar: React.FC<UserAvatarProps> = ({ initials, isPremium = false, size = 'lg', src, editable = false, onUpdated }) => {
	const { t } = useTranslation();

	const inputRef = useRef<HTMLInputElement>(null);
	const [broken, setBroken] = useState(false);
	const [uploading, setUploading] = useState(false);

	useEffect(() => {
		setBroken(false);
	}, [src]);

	const ringClass = isPremium ? premiumGradientClass : freeRingClass;
	const isSmall = size === 'sm';
	const innerClass = isSmall ? (isPremium ? `${premiumGradientClass} text-white` : 'bg-slate-400 text-white') : isPremium ? 'bg-[linear-gradient(180deg,#dbeafe_0%,#eff6ff_100%)] text-[#2a2867]' : 'bg-slate-200 text-slate-600';
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
		<div className={`rounded-full p-[2px] ${ringClass}`}>
			<div className={`relative rounded-full ${isSmall ? 'p-[2px]' : 'p-1'}`}>
				<div className={`relative overflow-hidden rounded-full bg-white ${isSmall ? 'h-11 w-11' : 'h-[8rem] w-[8rem]'}`}>
					{showImage ? (
						<RemoteImage src={src ?? undefined} alt="" className="h-full w-full rounded-full object-cover" spinnerSize={isSmall ? 14 : 22} onError={() => setBroken(true)} />
					) : (
						<div className={`flex h-full w-full items-center justify-center font-semibold ${isSmall ? 'text-sm' : 'text-[20px]'} ${innerClass}`}>{initials}</div>
					)}
					{(editable || uploading) && (
						<div
							className={`absolute inset-0 z-[2] flex items-center justify-center rounded-full bg-black/50 transition-opacity ${uploading ? 'opacity-100' : 'opacity-0 group-hover:opacity-100 group-active:opacity-100'}`}
						>
							{uploading ? <ImageSpinner size={isSmall ? 16 : 28} light /> : <CameraPlusOutline fill="#ffffff" size={isSmall ? 18 : 36} />}
						</div>
					)}
				</div>
				{!isSmall && (
					<div className="absolute top-3 right-2 bg-white p-1 rounded-full z-[999]">
						<Pencil fill="#000" size={13} />
					</div>
				)}
			</div>
		</div>
	);

	if (!editable) {
		return avatar;
	}

	return (
		<>
			<button type="button" className="group cursor-pointer rounded-full border-0 bg-transparent p-0" aria-label={t('account.change-avatar')} disabled={uploading} onClick={() => inputRef.current?.click()}>
				{avatar}
			</button>
			<input ref={inputRef} type="file" accept="image/jpeg,image/png,image/webp,image/gif" className="hidden" onChange={onFileChange} />
			{uploading && <Fallback text={t('account.avatar-uploading')} />}
		</>
	);
};

export default UserAvatar;
