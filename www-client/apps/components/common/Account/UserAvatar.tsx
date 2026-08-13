type UserAvatarSize = 'sm' | 'lg';

interface UserAvatarProps {
	initials: string;
	isPremium?: boolean;
	size?: UserAvatarSize;
}

const premiumRingClass = 'bg-[linear-gradient(98deg,#0085FF_0%,#22A1FF_45%,#60D1FF_100%)]';
const freeRingClass = 'bg-slate-400';

const UserAvatar: React.FC<UserAvatarProps> = ({ initials, isPremium = false, size = 'lg' }) => {
	const ringClass = isPremium ? premiumRingClass : freeRingClass;

	if (size === 'sm') {
		const innerClass = isPremium ? `${premiumRingClass} text-white` : 'bg-slate-400 text-white';

		return (
			<div className={`rounded-full p-[3px] ${ringClass}`}>
				<div className="rounded-full bg-white p-[2px]">
					<div className="relative h-11 w-11 overflow-hidden rounded-full bg-indigo-100">
						<div className={`flex h-full w-full items-center justify-center text-sm font-semibold ${innerClass}`}>{initials}</div>
					</div>
				</div>
			</div>
		);
	}

	const innerClass = isPremium ? 'bg-[linear-gradient(180deg,#dbeafe_0%,#eff6ff_100%)] text-[#2a2867]' : 'bg-slate-200 text-slate-600';

	return (
		<div className={`rounded-full p-[3px] ${ringClass}`}>
			<div className="rounded-full bg-white p-[2px]">
				<div className="relative h-[6rem] w-[6rem] overflow-hidden rounded-full bg-indigo-100">
					<div className={`flex h-full w-full items-center justify-center text-[20px] font-semibold ${innerClass}`}>{initials}</div>
				</div>
			</div>
		</div>
	);
};

export default UserAvatar;
