import { Link } from 'react-router-dom';
import { ROUTES } from '@/constants/constants';
import { useTranslation } from 'react-i18next';
import UserAvatar from '../common/Account/UserAvatar';
import PremiumBadgeMini from '../common/PremiumBadge/PremiumBadgeMini';
import { useLoginState } from '@/hooks/useLoginState';

const getGreetingKey = (): string => {
	const hour = new Date().getHours();

	if (hour < 12) {
		return 'home.greeting-morning';
	}

	if (hour < 18) {
		return 'home.greeting-afternoon';
	}

	return 'home.greeting-evening';
};

const Header = () => {
	const { t } = useTranslation();
	const { displayName, initials, isPremium } = useLoginState();

	return (
		<header className="mb-5 flex items-center justify-between gap-3">
			<Link to={ROUTES.ACCOUNT} className="flex min-w-0 items-center gap-3 cursor-pointer no-underline hover:no-underline">
				<UserAvatar initials={initials} isPremium={isPremium} size="sm" />

				<div className="min-w-0">
					<p className="truncate text-[13px] gu-text-muted">{t(getGreetingKey())}</p>
					<p className="truncate font-serif text-[18px] font-bold leading-tight gu-text-primary">{displayName}</p>
				</div>
			</Link>

			{!isPremium && <PremiumBadgeMini />}
		</header>
	);
};

export default Header;
