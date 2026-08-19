import { ROUTES } from '@/constants/constants';
import { lazy, LazyExoticComponent } from 'react';

const HomePage = lazy(() => import('@/pages/HomePage/index'));
const StartPage = lazy(() => import('@/pages/StartPage/index'));
const NotFoundPage = lazy(() => import('@/pages/NotFoundPage/index'));

const SubscriptionDetailPage = lazy(() => import('@/pages/TrackedSubscriptionPage/index'));
const CreateSubscriptionPage = lazy(() => import('@/pages/CreateSubscriptionPage/index'));

const PlansPage = lazy(() => import('@/pages/PlansPage/index'));
const PaidPage = lazy(() => import('@/pages/PaidPage/index'));
const AccountPage = lazy(() => import('@/pages/AccountPage/index'));
const ReferralPage = lazy(() => import('@/pages/ReferralPage/index'));
const AnalyticsPage = lazy(() => import('@/pages/AnalyticsPage/index'));
const AdminPage = lazy(() => import('@/pages/AdminPage/index'));
const LegalPage = lazy(() => import('@/pages/LegalPage/index'));

/** Auth */
const SignInPage = lazy(() => import('@/pages/SignInPage/index'));
const SignUpPage = lazy(() => import('@/pages/SignUpPage/index'));
const ConfirmEmail = lazy(() => import('@/pages/ConfirmEmailPage/index'));
const ForgotPasswordPage = lazy(() => import('@/pages/ForgotPasswordPage/index'));
const ResetPasswordPage = lazy(() => import('@/pages/ResetPasswordPage/index'));
const SubscriptionInvitePage = lazy(() => import('@/pages/SubscriptionInvitePage/index'));

export interface CustomRouteConfig {
	path: string;
	title?: string;
	loginRequired?: boolean;
	redirectIfLogged?: boolean;
	component: LazyExoticComponent<() => JSX.Element>;
}

const config: CustomRouteConfig[] = [
	{
		path: ROUTES.HOME,
		loginRequired: true,
		component: HomePage,
		title: 'label.page-home',
	},
	{
		path: ROUTES.START,
		loginRequired: false,
		redirectIfLogged: false,
		component: StartPage,
		title: 'label.page-start',
	},
	{
		path: ROUTES.SUBSCRIPTION_CREATE,
		loginRequired: true,
		component: CreateSubscriptionPage,
		title: 'subscription.create-title',
	},
	{
		path: ROUTES.SUBSCRIPTION_DETAIL,
		loginRequired: true,
		component: SubscriptionDetailPage,
	},
	{
		path: ROUTES.ANALYTICS,
		loginRequired: true,
		component: AnalyticsPage,
		title: 'label.page-analytics',
	},
	{
		path: ROUTES.ADMIN,
		loginRequired: true,
		component: AdminPage,
		title: 'label.page-admin',
	},
	{
		path: ROUTES.ACCOUNT,
		loginRequired: true,
		component: AccountPage,
		title: 'label.page-account',
	},
	{
		path: ROUTES.REFERRALS,
		loginRequired: true,
		component: ReferralPage,
		title: 'label.page-referrals',
	},
	{
		path: ROUTES.PLANS,
		loginRequired: true,
		component: PlansPage,
		title: 'label.page-plans',
	},
	{
		path: ROUTES.PAID,
		loginRequired: true,
		component: PaidPage,
		title: 'label.page-paid',
	},
	{
		path: ROUTES.NOT_FOUND,
		loginRequired: false,
		redirectIfLogged: false,
		component: NotFoundPage,
		title: 'label.page-not-found',
	},
	/** AUTH */
	{
		path: ROUTES.SIGNIN,
		loginRequired: false,
		redirectIfLogged: true,
		component: SignInPage,
		title: 'label.page-signin',
	},
	{
		path: ROUTES.SIGNUP,
		loginRequired: false,
		redirectIfLogged: true,
		component: SignUpPage,
		title: 'label.page-signup',
	},
	{
		path: ROUTES.AUTH_CONFIRM_EMAIL,
		loginRequired: false,
		redirectIfLogged: true,
		component: ConfirmEmail,
		title: 'label.page-confirm-email',
	},
	{
		path: ROUTES.AUTH_FORGOT_PASSWORD,
		loginRequired: false,
		redirectIfLogged: true,
		component: ForgotPasswordPage,
		title: 'label.page-forgot-password',
	},
	{
		path: ROUTES.AUTH_RESET_PASSWORD,
		loginRequired: false,
		redirectIfLogged: false,
		component: ResetPasswordPage,
		title: 'label.page-reset-password',
	},
	{
		path: ROUTES.SUBSCRIPTION_INVITE,
		loginRequired: false,
		redirectIfLogged: false,
		component: SubscriptionInvitePage,
		title: 'subscription.invite-title',
	},
	{
		path: ROUTES.LEGAL,
		loginRequired: false,
		redirectIfLogged: false,
		component: LegalPage,
	},
];

export default config;
