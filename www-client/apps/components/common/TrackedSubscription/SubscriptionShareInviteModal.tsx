import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import GUIButton from '@/components/ui/Button/GUIButton';
import MemberAvatar from '@/components/common/TrackedSubscription/MemberAvatar';
import ShareWedgeSlider from '@/components/common/TrackedSubscription/ShareWedgeSlider';
import { formatSubscriptionPrice } from '@/utils/TrackedSubscriptionDisplayUtils';

export interface ShareInvitePerson {
	email: string;
	name: string;
	initials: string;
	avatars?: string | null;
}

interface SubscriptionShareInviteModalProps {
	people: ShareInvitePerson[];
	subscriptionName: string;
	price: number;
	currency: string;
	ownerShare: number;
	alreadyShared: number;
	onConfirm: (sharePercent: number) => Promise<void>;
}

const MIN_SHARE = 0;

const roundShare = (value: number) => Math.round(value * 10) / 10;

const formatShare = (value: number) => {
	const rounded = roundShare(value);
	return Number.isInteger(rounded) ? String(rounded) : rounded.toFixed(1);
};

const shareAmount = (price: number, percent: number) => Math.round(((price * percent) / 100) * 1000) / 1000;

const SubscriptionShareInviteModal: React.FC<SubscriptionShareInviteModalProps> = ({ people, subscriptionName, price, currency, ownerShare, alreadyShared, onConfirm }) => {
	const { t, i18n } = useTranslation();

	const count = Math.max(people.length, 1);
	const maxShare = Math.max(MIN_SHARE, roundShare(ownerShare / count));
	const equalShare = Math.min(maxShare, roundShare(ownerShare / (count + 1)));

	const [shareEach, setShareEach] = useState(equalShare);
	const [submitting, setSubmitting] = useState(false);

	const ownerRemaining = Math.max(0, roundShare(ownerShare - shareEach * count));
	const inviteTotal = roundShare(shareEach * count);
	const valid = shareEach >= 0 && inviteTotal <= ownerShare + 0.001;

	const levelKey = useMemo(() => {
		if (shareEach < equalShare * 0.85) {
			return 'subscription.share-level-light';
		}

		if (shareEach > equalShare * 1.15) {
			return 'subscription.share-level-heavy';
		}

		return 'subscription.share-level-fair';
	}, [equalShare, shareEach]);

	const money = (percent: number) => formatSubscriptionPrice(shareAmount(price, percent), currency, i18n.language);

	const onSubmit = async () => {
		if (!valid || submitting) {
			return;
		}

		setSubmitting(true);

		try {
			await onConfirm(shareEach);
		} finally {
			setSubmitting(false);
		}
	};

	return (
		<div className="space-y-4">
			<div className="gu-glass-card flex items-center gap-3 px-3 py-3">
				<div className="flex items-center">
					{people.slice(0, 3).map((person, index) => (
						<div key={person.email} className={`rounded-full ring-2 ring-[var(--surface)] ${index > 0 ? '-ml-2' : ''}`} style={{ zIndex: people.length - index }}>
							<MemberAvatar initials={person.initials} src={person.avatars} size={40} />
						</div>
					))}
				</div>
				<div className="min-w-0 flex-1">
					<p className="truncate text-[15px] font-semibold gu-text-primary">{people.length === 1 ? people[0].name : t('subscription.members-count', { count: people.length })}</p>
					<p className="truncate text-[13px] gu-text-muted">{subscriptionName}</p>
				</div>
			</div>

			<section className="space-y-2">
				<div className="flex items-start justify-between gap-3">
					<div>
						<p className="text-[15px] font-semibold gu-text-primary">{t('subscription.share-modal-question')}</p>
						<p className="mt-1 text-[11px] font-semibold tracking-[0.14em] text-[#f97316]">
							{t(levelKey)} · {formatShare(shareEach)}%
						</p>
					</div>
					<p className="shrink-0 text-[13px] font-semibold gu-text-muted">{money(shareEach)}</p>
				</div>

				<div className="relative">
					{shareEach > equalShare * 1.02 ? <p className="absolute -top-1 right-0 text-[11px] font-semibold text-[#ef4444]">{t('subscription.share-above-equal')}</p> : null}

					<ShareWedgeSlider
						value={shareEach}
						min={MIN_SHARE}
						max={maxShare}
						marker={equalShare}
						markerLabel={t('subscription.share-equal-mark')}
						minLabel={t('subscription.share-slider-min')}
						maxLabel={t('subscription.share-slider-max')}
						onChange={setShareEach}
						ariaLabel={t('subscription.share-percent')}
					/>
				</div>
			</section>

			<div className="space-y-2 text-[13px]">
				{people.map((person) => (
					<div key={person.email} className="flex items-center justify-between gap-3">
						<p className="min-w-0 truncate gu-text-muted">{person.name}</p>
						<p className="shrink-0 font-semibold gu-text-primary">
							{formatShare(shareEach)}% · {money(shareEach)}
						</p>
					</div>
				))}

				<div className="flex items-center justify-between gap-3">
					<p className="gu-text-muted">{t('subscription.share-you-keep')}</p>
					<p className="shrink-0 font-semibold gu-text-primary">
						{formatShare(ownerRemaining)}% · {money(ownerRemaining)}
					</p>
				</div>

				{alreadyShared > 0.01 ? (
					<div className="flex items-center justify-between gap-3">
						<p className="gu-text-muted">{t('subscription.share-already')}</p>
						<p className="shrink-0 font-semibold gu-text-primary">
							{formatShare(alreadyShared)}% · {money(alreadyShared)}
						</p>
					</div>
				) : null}

				<div className="flex items-center justify-between gap-3 border-t border-dashed border-[var(--divider)] pt-2">
					<p className="font-semibold gu-text-primary">{t('subscription.share-total')}</p>
					<p className="shrink-0 font-semibold gu-text-primary">{formatSubscriptionPrice(price, currency, i18n.language)}</p>
				</div>
			</div>

			<GUIButton type="button" variant="primary" disabled={!valid} isLoading={submitting} loadingText={t('subscription.share-sending')} onClick={onSubmit}>
				{t('subscription.share-send')}
			</GUIButton>
		</div>
	);
};

export default SubscriptionShareInviteModal;
