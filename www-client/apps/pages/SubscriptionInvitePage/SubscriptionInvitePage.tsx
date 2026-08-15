import { useEffect, useState } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { ROUTES, buildRoute } from '@/constants/constants';
import { CACHEKEYs } from '@/constants/CacheKeys.constants';
import GUIButton from '@/components/ui/Button/GUIButton';
import { notify } from '@/components/Notification/notify';
import AuthPageLayout from '@/components/common/Auth/AuthPageLayout';
import PageLoadingState from '@/components/common/PageLoadingState/PageLoadingState';
import { basicTrackedSubscriptionAcceptInvite, basicTrackedSubscriptionInvitePreview, type SharedSubscriptionInvitePreview } from '@/rest/trackedSubscriptionAPI';
import { readAuthSession } from '@/utils/authSessionUtils';

const SubscriptionInvitePage = () => {
	const { t } = useTranslation();
	const navigate = useNavigate();
	const [searchParams] = useSearchParams();
	const token = searchParams.get('token') ?? '';
	const isLoggedIn = readAuthSession();

	const [preview, setPreview] = useState<SharedSubscriptionInvitePreview | null>(null);
	const [loading, setLoading] = useState(Boolean(token));
	const [accepting, setAccepting] = useState(false);

	useEffect(() => {
		if (!token) {
			setLoading(false);
			return;
		}

		if (!isLoggedIn) {
			sessionStorage.setItem(CACHEKEYs.SUBSCRIPTION_INVITE_TOKEN, token);
			setLoading(false);
			return;
		}

		let cancelled = false;

		const loadPreview = async () => {
			try {
				const data = await basicTrackedSubscriptionInvitePreview(token);
				if (!cancelled) {
					setPreview(data);
				}
			} catch {
				if (!cancelled) {
					setPreview(null);
				}
			} finally {
				if (!cancelled) {
					setLoading(false);
				}
			}
		};

		void loadPreview();

		return () => {
			cancelled = true;
		};
	}, [isLoggedIn, token]);

	const onAccept = async () => {
		if (!token || accepting) {
			return;
		}

		setAccepting(true);

		try {
			const result = await basicTrackedSubscriptionAcceptInvite(token);
			sessionStorage.removeItem(CACHEKEYs.SUBSCRIPTION_INVITE_TOKEN);
			notify.success(result.message || t('subscription.invite-accepted'));
			navigate(buildRoute(ROUTES.SUBSCRIPTION_DETAIL, { id: result.subscription_id }), { replace: true });
		} finally {
			setAccepting(false);
		}
	};

	if (loading) {
		return <PageLoadingState />;
	}

	if (!token) {
		return (
			<AuthPageLayout title={t('subscription.invite-title')} subtitle={t('subscription.invite-invalid')}>
				<Link to={ROUTES.HOME} className="font-semibold text-[#0085FF] no-underline hover:no-underline">
					{t('action.go-home')}
				</Link>
			</AuthPageLayout>
		);
	}

	if (!isLoggedIn) {
		return (
			<AuthPageLayout title={t('subscription.invite-title')} subtitle={t('subscription.invite-signin-desc')}>
				<GUIButton type="button" variant="primary" onClick={() => navigate(ROUTES.SIGNIN)}>
					{t('action.sign-in')}
				</GUIButton>
			</AuthPageLayout>
		);
	}

	if (!preview) {
		return (
			<AuthPageLayout title={t('subscription.invite-title')} subtitle={t('subscription.invite-invalid')}>
				<Link to={ROUTES.HOME} className="font-semibold text-[#0085FF] no-underline hover:no-underline">
					{t('action.go-home')}
				</Link>
			</AuthPageLayout>
		);
	}

	return (
		<AuthPageLayout
			title={t('subscription.invite-title')}
			subtitle={t('subscription.invite-desc', {
				owner: preview.owner_name,
				name: preview.subscription_name,
				percent: preview.share_percent,
			})}
		>
			<GUIButton type="button" variant="primary" isLoading={accepting} onClick={onAccept} loadingText={t('subscription.invite-accepting')}>
				{t('subscription.invite-accept')}
			</GUIButton>
		</AuthPageLayout>
	);
};

export default SubscriptionInvitePage;
