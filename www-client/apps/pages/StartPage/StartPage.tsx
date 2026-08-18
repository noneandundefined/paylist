import { Link } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { ROUTES } from '@/constants/constants';
import { getAuthState } from '@/private-route';
import { readAuthSession } from '@/utils/authSessionUtils';

const StartPage = () => {
	const { t } = useTranslation();
	const [isLoggedIn, setIsLoggedIn] = useState(readAuthSession);

	useEffect(() => {
		void getAuthState().then(setIsLoggedIn);
	}, []);

	return (
		<div className="flex min-h-screen items-center justify-center gu-page-bg px-5 py-16 sm:px-10 lg:px-[10rem]">
			<div className="flex w-full max-w-6xl space-y-8 flex-col items-center justify-center text-center">
				<h1 className="text-[36px] font-extrabold leading-[1.1] gu-text-primary sm:text-[52px]">
					{t('start.hero-title-lead')} <span className="rounded-md bg-[#d7ff00] px-1.5 text-black dark:bg-transparent dark:px-0 dark:text-[#d7ff00]">{t('start.hero-title-accent')}</span>
				</h1>

				<p className="max-w-2xl text-[16px] leading-relaxed gu-text-secondary sm:text-[18px]">{t('start.hero-subtitle')}</p>

				<div>
					<Link
						to={isLoggedIn ? ROUTES.HOME : ROUTES.SIGNIN}
						className="inline-flex items-center justify-center rounded-[16px] bg-[#d7ff00] w-[20rem] py-3.5 text-[16px] font-semibold text-black no-underline transition hover:bg-[#d7ff00]/90 hover:no-underline"
					>
						{t(isLoggedIn ? 'start.cta-home' : 'start.cta-manage')}
					</Link>
					<p className="mt-2 text-[13px] gu-text-muted">{t('start.save-hint')}</p>
				</div>
			</div>
		</div>
	);
};

export default StartPage;
