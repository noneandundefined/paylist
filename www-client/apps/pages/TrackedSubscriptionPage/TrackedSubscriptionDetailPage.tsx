import PageLayout from '../PageLayout';
import { ROUTES } from '@/constants/constants';
import { Helmet } from 'react-helmet-async';
import { useTranslation } from 'react-i18next';
import { useEffect, useMemo, useState } from 'react';
import { useNavigate, Navigate, useParams } from 'react-router-dom';
import { QUERY_KEYS } from '@/constants/QueryKeys.constants';
import GUIButton from '@/components/ui/Button/GUIButton';
import { notify } from '@/components/Notification/notify';
import PageHeader from '@/components/common/PageHeader/PageHeader';
import PageLoadingState from '@/components/common/PageLoadingState/PageLoadingState';
import CategoryChipGroup from '@/components/common/Category/CategoryChipGroup';
import SubscriptionSettingsPanel from '@/components/common/TrackedSubscription/SubscriptionSettingsPanel';
import SubscriptionSummaryCard from '@/components/common/TrackedSubscription/SubscriptionSummaryCard';
import { useConfirm } from '@/hooks/useConfirm';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { usePremiumFeatureGate } from '@/hooks/usePremiumFeatureGate';
import { useSubscriptionCategories } from '@/hooks/useSubscriptionCategories';
import { useInvalidateSubscriptions } from '@/hooks/useInvalidateSubscriptions';
import { areCategoriesEqual, isCustomCategory } from '@/utils/categoryDisplayUtils';
import { notifyPremiumRequired } from '@/utils/premiumUtils';
import { clampPremiumSubscriptionFlags, resolvePremiumSubscriptionFlags } from '@/utils/subscriptionPremiumUtils';
import { getNextBillingLabel } from '@/utils/TrackedSubscriptionDisplayUtils';
import { isSubscriptionOverdue } from '@/utils/SubscriptionRenewalUtils';
import SubscriptionSharingPanel from '@/components/common/TrackedSubscription/SubscriptionSharingPanel';
import { basicTrackedSubscriptionDelete, basicTrackedSubscriptionGetById, basicTrackedSubscriptionLeave, basicTrackedSubscriptionMembers, basicTrackedSubscriptionUpdate } from '@/rest/trackedSubscriptionAPI';
import { GUITextarea } from '@/components/ui/Input/GUITextarea';

import Close from '@/components/@icons/close';

const normalizeDatePay = (value: string) => value.slice(0, 10);

const normalizeNote = (value: string | null | undefined) => (value ?? '').trim();

const TrackedSubscriptionDetailPageDetailPage = () => {
	const { t, i18n } = useTranslation();
	const { confirm } = useConfirm();

	const navigate = useNavigate();
	const { invalidateAfterUpdate, invalidateAfterDelete } = useInvalidateSubscriptions();
	const { canUseNotification, isPremium } = usePremiumFeatureGate();
	const { categories: availableCategories } = useSubscriptionCategories();

	const { id = '' } = useParams<{ id: string }>();
	const subscriptionId = Number(id);

	const { data: subscription, loading } = useHandleServer([QUERY_KEYS.trackedSubscriptionDetail, subscriptionId], () => basicTrackedSubscriptionGetById(subscriptionId), {
		enabled: Number.isFinite(subscriptionId) && subscriptionId > 0,
	});

	const { data: membersData, reload: reloadMembers } = useHandleServer([QUERY_KEYS.trackedSubscriptionMembers, subscriptionId], () => basicTrackedSubscriptionMembers(subscriptionId), {
		enabled: Number.isFinite(subscriptionId) && subscriptionId > 0,
	});

	const isOwner = Boolean(subscription?.is_owner);

	const [categories, setCategories] = useState<string[]>([]);
	const [autoRenewal, setAutoRenewal] = useState(false);
	const [notification, setNotification] = useState(false);
	const [includeInAnalytics, setIncludeInAnalytics] = useState(false);
	const [note, setNote] = useState('');
	const [saving, setSaving] = useState(false);
	const [deleting, setDeleting] = useState(false);

	useEffect(() => {
		if (!subscription?.name) {
			return;
		}

		document.title = subscription.name;
	}, [subscription?.name]);

	useEffect(() => {
		if (!subscription) {
			return;
		}

		setCategories(subscription.categories ?? []);
		const flags = clampPremiumSubscriptionFlags(subscription.auto_renewal, subscription.notification, canUseNotification);
		setAutoRenewal(flags.autoRenewal);
		setNotification(flags.notification);
		setIncludeInAnalytics(subscription.include_in_analytics);
		setNote(subscription.note ?? '');
	}, [subscription, canUseNotification]);

	const isDirty = useMemo(() => {
		if (!subscription) {
			return false;
		}

		const baseline = clampPremiumSubscriptionFlags(subscription.auto_renewal, subscription.notification, canUseNotification);

		const ownerFieldsDirty = isOwner && autoRenewal !== baseline.autoRenewal;

		return (
			ownerFieldsDirty ||
			notification !== baseline.notification ||
			includeInAnalytics !== subscription.include_in_analytics ||
			normalizeNote(note) !== normalizeNote(subscription.note) ||
			!areCategoriesEqual(categories, subscription.categories ?? [])
		);
	}, [subscription, autoRenewal, notification, includeInAnalytics, note, categories, canUseNotification, isOwner]);

	if (!Number.isFinite(subscriptionId) || subscriptionId <= 0) {
		return <Navigate to={ROUTES.NOT_FOUND} replace />;
	}

	if (loading) {
		return <PageLoadingState />;
	}

	if (!subscription) {
		return <Navigate to={ROUTES.NOT_FOUND} replace />;
	}

	const periodLabel = t(`home.period-${subscription.period}`).toLowerCase();
	const overdue = isSubscriptionOverdue(subscription.date_pay);

	const toggleCategory = (slug: string) => {
		const category = availableCategories.find((item) => item.slug === slug);

		if (!isPremium && category && isCustomCategory(category)) {
			notifyPremiumRequired(t);
			return;
		}

		setCategories((prev) => (prev.includes(slug) ? prev.filter((item) => item !== slug) : [...prev, slug]));
	};

	const onUpdate = async () => {
		if (!isDirty || saving) {
			return;
		}

		setSaving(true);

		try {
			const premiumFlags = resolvePremiumSubscriptionFlags({
				autoRenewal,
				notification,
				canUseNotification,
			});

			await basicTrackedSubscriptionUpdate(subscription.id, {
				name: subscription.name,
				price: subscription.price,
				currency: subscription.currency,
				period: subscription.period,
				date_pay: normalizeDatePay(subscription.date_pay),
				auto_renewal: premiumFlags.auto_renewal,
				notification: premiumFlags.notification,
				include_in_analytics: includeInAnalytics,
				categories,
				note: normalizeNote(note) === '' ? null : normalizeNote(note),
			});

			await invalidateAfterUpdate(subscriptionId);

			notify.success(t('subscription.update-success'));
		} finally {
			setSaving(false);
		}
	};

	const onCancelSubscription = async () => {
		const confirmKey = isOwner ? 'subscription.cancel-confirm-desc' : 'subscription.leave-confirm-desc';
		const titleKey = isOwner ? 'subscription.cancel' : 'subscription.leave';

		if (deleting || !(await confirm(confirmKey, titleKey))) {
			return;
		}

		setDeleting(true);

		try {
			const message = isOwner ? await basicTrackedSubscriptionDelete(subscription.id) : await basicTrackedSubscriptionLeave(subscription.id);

			await invalidateAfterDelete(subscriptionId);

			notify.success(message || (isOwner ? t('subscription.cancel-success') : t('subscription.leave-success')));
			navigate(ROUTES.HOME, { replace: true });
		} finally {
			setDeleting(false);
		}
	};

	return (
		<PageLayout>
			<Helmet>
				<title>{subscription.name}</title>
			</Helmet>
			<div className="flex flex-col space-y-3">
				<PageHeader title={t('subscription.detail-title')} backTo={ROUTES.HOME} backLabel={t('action.back')} />

				<SubscriptionSummaryCard subscription={subscription} periodLabel={periodLabel} overdue={overdue} />

				{overdue ? (
					<div className="gu-glass-card gu-overdue-surface px-4 py-3.5">
						<p className="text-[15px] font-semibold gu-overdue-title">{t('subscription.overdue-title')}</p>
						<p className="mt-1 text-[13px] leading-relaxed gu-overdue-muted">{t('subscription.overdue-desc')}</p>
					</div>
				) : null}

				<p className={`text-center text-[13px] ${overdue ? 'gu-overdue-muted' : 'gu-text-muted'}`}>{getNextBillingLabel(subscription.date_pay, t, i18n.language)}</p>

				<GUIButton
					type="button"
					onClick={onCancelSubscription}
					disabled={saving || deleting}
					isLoading={deleting}
					loadingText={t('subscription.cancel-loading')}
					className="gu-glass-card flex w-full items-center gap-2.5 px-3 py-3.5 text-left transition hover:bg-[var(--surface-muted)]"
				>
					<span className="gu-glass-icon-btn">
						<Close fill="currentColor" size={18} />
					</span>

					<span className="font-medium leading-tight gu-text-primary">{isOwner ? t('subscription.cancel') : t('subscription.leave')}</span>
				</GUIButton>

				<section className="gu-glass-card space-y-3 px-3 py-3.5">
					<h2 className="text-[15px] font-semibold gu-text-primary">{t('subscription.note')}</h2>
					<p className="text-[12px] gu-text-muted">{t('subscription.note-personal')}</p>

					<GUITextarea value={note} onChange={(event) => setNote(event.target.value)} placeholder={t('subscription.note-placeholder')} aria-label={t('subscription.note')} maxLength={2000} />
				</section>

				<section className="gu-glass-card space-y-3 px-3 py-3.5">
					<h2 className="text-[15px] font-semibold gu-text-primary">{t('subscription.category')}</h2>
					<p className="text-[12px] gu-text-muted">{t('subscription.category-personal')}</p>

					<CategoryChipGroup categories={availableCategories} selectedSlugs={categories} onToggle={toggleCategory} />
				</section>

				<SubscriptionSharingPanel
					subscription={subscription}
					membersData={membersData}
					onChanged={async () => {
						await invalidateAfterUpdate(subscriptionId);
						await reloadMembers();
					}}
				/>

				<SubscriptionSettingsPanel
					autoRenewal={autoRenewal}
					notification={notification}
					includeInAnalytics={includeInAnalytics}
					onAutoRenewalChange={setAutoRenewal}
					onNotificationChange={setNotification}
					onIncludeInAnalyticsChange={setIncludeInAnalytics}
					canUseNotification={canUseNotification}
					onPremiumRequired={() => notifyPremiumRequired(t)}
					canChangeAutoRenewal={isOwner}
				/>

				{isDirty && (
					<GUIButton type="submit" variant="primary" onClick={onUpdate} disabled={saving} loadingText={t('subscription.update-loading')}>
						{t('subscription.update-action')}
					</GUIButton>
				)}
			</div>
		</PageLayout>
	);
};

export default TrackedSubscriptionDetailPageDetailPage;
