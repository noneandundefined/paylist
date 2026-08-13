import { ROUTES } from '@/constants/constants';
import { useTranslation } from 'react-i18next';
import { Link, useLocation } from 'react-router-dom';
import { useTheme } from '@/context/ThemeContext';
import ElevationRise from '../@icons/elevation-rise';

import Home from '../@icons/home';
import Plus from '../@icons/plus';

const NAV_ACTIVE_COLOR = '#0085FF';

const navItems = [
	{ path: ROUTES.HOME, labelKey: 'label.page-home', Icon: Home },
	{ path: ROUTES.ANALYTICS, labelKey: 'home.analytics', Icon: ElevationRise },
] as const;

const Navigation = () => {
	const { t } = useTranslation();
	const { pathname } = useLocation();
	const { isDark } = useTheme();

	const inactiveColor = isDark ? '#f1f5f9' : '#000000';

	return (
		<aside className="pointer-events-none fixed inset-x-0 bottom-0 z-20 flex justify-center px-4 pb-5">
			<div className="pointer-events-auto flex flex-1 max-w-[90%] items-center gap-5 md:max-w-[70%] xl:max-w-[50%]">
				<nav className="gu-glass-pill flex flex-1 items-center gap-1 px-2 py-2" aria-label={t('home.navigation')}>
					{navItems.map(({ path, labelKey, Icon }) => {
						const isActive = pathname === path;
						const color = isActive ? NAV_ACTIVE_COLOR : inactiveColor;

						return (
							<Link key={path} to={path} aria-current={isActive ? 'page' : undefined} className={`inline-flex h-11 flex-1 items-center justify-center rounded-full no-underline hover:no-underline`}>
								<div className="flex flex-col items-center">
									<Icon fill={color} size={27} />
									<p className={`text-[13px] font-medium ${isActive ? 'text-[#0085FF]' : isDark ? 'text-slate-100' : 'text-black'}`}>{t(labelKey)}</p>
								</div>
							</Link>
						);
					})}
				</nav>

				<Link to={ROUTES.SUBSCRIPTION_CREATE} className="gu-glass-fab no-underline hover:no-underline" aria-label={t('home.add-subscription')}>
					<Plus fill="currentColor" size={27} />
				</Link>
			</div>
		</aside>
	);
};

export default Navigation;
