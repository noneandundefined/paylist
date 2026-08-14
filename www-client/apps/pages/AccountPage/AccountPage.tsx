import { useEffect, useState } from 'react';
import PageLayout from '../PageLayout';
import { ROUTES } from '@/constants/constants';
import { useTranslation } from 'react-i18next';
import { useTheme } from '@/context/ThemeContext';
import { QUERY_KEYS } from '@/constants/QueryKeys.constants';
import { basicAuthSignOut } from '@/rest/authAPI';
import GUIButton from '@/components/ui/Button/GUIButton';
import GUISelect from '@/components/ui/Select/GUISelect';
import UserAvatar from '@/components/common/Account/UserAvatar';
import AccountSection from '@/components/common/Account/AccountSection';
import PremiumGatedSection from '@/components/common/Account/PremiumGatedSection';
import AccountSettingsRow from '@/components/common/Account/AccountSettingsRow';
import AccountCategoriesManager from '@/components/common/Account/AccountCategoriesManager';
import AccountTelegramNotifications from '@/components/common/Account/AccountTelegramNotifications';
import PageLoadingState from '@/components/common/PageLoadingState/PageLoadingState';
import CurrencySelect from '@/components/common/Currency/CurrencySelect';
import CountrySelect from '@/components/common/Country/CountrySelect';
import PremiumBadge from '@/components/common/PremiumBadge/PremiumBadge';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { useLoginState } from '@/hooks/useLoginState';
import { useConfirm } from '@/hooks/useConfirm';
import { useInvalidateSubscriptions } from '@/hooks/useInvalidateSubscriptions';
import { getAppLanguage, SUPPORTED_LANGUAGES } from '@/constants/Language.constant';
import { notifyPremiumRequired, notifySoon } from '@/utils/premiumUtils';
import { clearAuthSession } from '@/utils/authSessionUtils';
import { basicUserAccountDelete, basicUserSettingsGet, basicUserSettingsUpdate } from '@/rest/userAPI';

import Delete from '@/components/@icons/delete';
import FaceAgent from '@/components/@icons/face-agent';
import LogoutVariant from '@/components/@icons/logout-variant';
import ThemeLightDark from '@/components/@icons/theme-light-dark';

const AccountPage = () => {
	const { t, i18n } = useTranslation();
	const { isDark, setTheme } = useTheme();
	const { confirm } = useConfirm();
	const { invalidateListAndSummary } = useInvalidateSubscriptions();
	const { loginState, loading, reload, displayName, initials, avatar, isPremium, canUseNotification } = useLoginState();

	const [displayCurrency, setDisplayCurrency] = useState('USD');
	const [country, setCountry] = useState('US');
	const [deletingAccount, setDeletingAccount] = useState(false);

	const currentLanguage = getAppLanguage(i18n.language);

	const {
		data: userSettings,
		loading: settingsLoading,
		reload: reloadSettings,
	} = useHandleServer([QUERY_KEYS.userSettings], () => basicUserSettingsGet(), {
		enabled: Boolean(loginState),
	});

	useEffect(() => {
		if (isPremium && userSettings?.display_currency) {
			setDisplayCurrency(userSettings.display_currency);
		}

		if (isPremium && userSettings?.country) {
			setCountry(userSettings.country);
		}
	}, [isPremium, userSettings?.display_currency, userSettings?.country]);

	const onLanguageChange: React.ChangeEventHandler<HTMLSelectElement> = (event) => {
		const next = event.target.value;
		localStorage.setItem('lang', next);
		void i18n.changeLanguage(next);
	};

	const onCurrencyChange: React.ChangeEventHandler<HTMLSelectElement> = async (event) => {
		if (!isPremium) {
			notifyPremiumRequired(t);
			return;
		}

		const next = event.target.value;
		setDisplayCurrency(next);

		await basicUserSettingsUpdate({ display_currency: next });
		await invalidateListAndSummary();
	};

	const onCountryChange: React.ChangeEventHandler<HTMLSelectElement> = async (event) => {
		if (!isPremium) {
			notifyPremiumRequired(t);
			return;
		}

		const next = event.target.value;
		setCountry(next);

		await basicUserSettingsUpdate({ country: next });
	};

	if (loading || !loginState || settingsLoading) {
		return <PageLoadingState />;
	}

	const onLogout = async () => {
		if (!(await confirm('account.logout-confirm-desc', 'account.logout'))) {
			return;
		}

		await basicAuthSignOut();
	};

	const onDeleteAccount = async () => {
		if (deletingAccount || !(await confirm('account.delete-confirm-desc', 'account.delete-account'))) {
			return;
		}

		setDeletingAccount(true);

		try {
			await basicUserAccountDelete();
			clearAuthSession();
			window.location.replace(ROUTES.SIGNIN);
		} finally {
			setDeletingAccount(false);
		}
	};

	return (
		<PageLayout>
			<div className="mx-auto flex w-full flex-col space-y-5">
				<header className="flex justify-end">
					<GUIButton type="button" onClick={onLogout} isLoading={false} className="gu-glass-icon-btn shrink-0" aria-label={t('account.logout')}>
						<LogoutVariant fill="currentColor" size={22} />
					</GUIButton>
				</header>

				<section className="flex flex-col items-center px-2 text-center">
					<UserAvatar initials={initials} isPremium={isPremium} size="lg" src={avatar} editable onUpdated={() => void reload()} />

					<p className="mt-1 font-serif text-[22px] font-bold leading-tight gu-text-primary">{displayName}</p>
					<p className="mt-1 text-[14px] gu-text-muted">
						{t('account.email-label')}: {loginState.email}
					</p>
				</section>

				{!isPremium && <PremiumBadge />}

				<AccountSection title={t('account.settings-section')}>
					<div className="gu-glass-card divide-y gu-divide overflow-hidden">
						<div className="p-3 space-y-2">
							<div className="flex items-center gap-3">
								<h2 className="text-[15px] font-semibold gu-text-primary">{t('account.translate')}</h2>
							</div>

							<GUISelect value={currentLanguage} onChange={onLanguageChange} modalTitle={t('account.translate')} aria-label={t('account.translate')}>
								{SUPPORTED_LANGUAGES.map(({ code, labelKey }) => (
									<option key={code} value={code}>
										{t(labelKey)}
									</option>
								))}
							</GUISelect>
						</div>

						<AccountSettingsRow
							icon={<ThemeLightDark fill="currentColor" size={21} />}
							label={t('account.mode')}
							value={isDark ? t('account.mode-dark') : t('account.mode-light')}
							trailing="switch"
							switchChecked={isDark}
							onSwitchChange={(checked) => setTheme(checked ? 'dark' : 'light')}
						/>
					</div>
				</AccountSection>

				<AccountSection title={t('account.subscriptions-section')}>
					<div className="gu-glass-card divide-y gu-divide overflow-hidden">
						<AccountTelegramNotifications
							isPremium={isPremium}
							canUseNotification={canUseNotification}
							connected={Boolean(userSettings?.telegram_connected)}
							username={userSettings?.telegram_username}
							onChanged={() => {
								void reloadSettings();
							}}
						/>

						<PremiumGatedSection title={t('account.currency')} isPremium={isPremium}>
							<CurrencySelect value={displayCurrency} onChange={onCurrencyChange} />
						</PremiumGatedSection>

						<AccountCategoriesManager isPremium={isPremium} />
					</div>
				</AccountSection>

				<AccountSection title={t('account.analytics-section')}>
					<div className="gu-glass-card divide-y gu-divide overflow-hidden">
						<PremiumGatedSection title={t('account.country')} isPremium={isPremium}>
							<CountrySelect value={country} onChange={onCountryChange} />
						</PremiumGatedSection>
					</div>
				</AccountSection>

				<AccountSection title={t('account.content-section')}>
					<div className="gu-glass-card divide-y gu-divide overflow-hidden">
						<AccountSettingsRow icon={<FaceAgent fill="currentColor" size={21} />} label={t('account.feedback')} onClick={() => notifySoon(t, 'account.feedback-soon')} />
						<AccountSettingsRow
							icon={<Delete fill="#dc2626" size={21} />}
							label={deletingAccount ? t('account.delete-account-progress') : t('account.delete-account')}
							onClick={onDeleteAccount}
							disabled={deletingAccount}
							danger
						/>
					</div>
				</AccountSection>
			</div>
		</PageLayout>
	);
};

export default AccountPage;
