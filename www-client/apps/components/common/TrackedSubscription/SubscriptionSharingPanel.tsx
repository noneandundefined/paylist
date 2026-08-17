import { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { ROUTES } from '@/constants/constants';
import Check from '@/components/@icons/check';
import Plus from '@/components/@icons/plus';
import AccentBadge from '@/components/common/AccentBadge/AccentBadge';
import MemberAvatar from '@/components/common/TrackedSubscription/MemberAvatar';
import GUIButton from '@/components/ui/Button/GUIButton';
import { GUInput } from '@/components/ui/Input/GUInput';
import { notify } from '@/components/Notification/notify';
import Modal from '@/components/Modal/Modal';
import SubscriptionShareInviteModal from '@/components/common/TrackedSubscription/SubscriptionShareInviteModal';
import { useConfirm } from '@/hooks/useConfirm';
import { useLoginState } from '@/hooks/useLoginState';
import { useModalContext } from '@/context/useModalContext';
import {
	basicTrackedSubscriptionInvite,
	basicTrackedSubscriptionRemoveMember,
	basicTrackedSubscriptionVoteShares,
	type SharedSubscriptionMember,
	type SharedSubscriptionMembersResponse,
	type TrackedSubscriptionResponse,
} from '@/rest/trackedSubscriptionAPI';
import { basicUserSearchByEmail, type UserPublicProfile } from '@/rest/userAPI';
import { getInitialsFromName } from '@/utils/stringUtils';
import { formatSubscriptionName } from '@/utils/TrackedSubscriptionDisplayUtils';

interface SubscriptionSharingPanelProps {
	subscription: TrackedSubscriptionResponse;
	membersData?: SharedSubscriptionMembersResponse | null;
	onChanged: () => Promise<void> | void;
}

interface InviteCandidate {
	email: string;
	first_name?: string | null;
	last_name?: string | null;
	avatars?: string | null;
}

const memberName = (person: { first_name?: string | null; last_name?: string | null; email: string }) => {
	const name = [person.first_name, person.last_name].filter((part) => part && part.trim()).join(' ');
	return name || person.email.split('@')[0] || person.email;
};

const memberInitials = (person: { first_name?: string | null; last_name?: string | null; email: string }) => getInitialsFromName(memberName(person));

const isEmailQuery = (value: string) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value);

const SubscriptionSharingPanel: React.FC<SubscriptionSharingPanelProps> = ({ subscription, membersData, onChanged }) => {
	const { t } = useTranslation();
	const { confirm } = useConfirm();
	const { open, close } = useModalContext();
	const { loginState, isPremium } = useLoginState();

	const members = membersData?.members ?? [];
	const proposal = membersData?.pending_proposal ?? null;
	const isOwner = Boolean(subscription.is_owner);
	const sharePercent = subscription.share_percent ?? 100;

	const [searchOpen, setSearchOpen] = useState(false);
	const [search, setSearch] = useState('');
	const [results, setResults] = useState<UserPublicProfile[]>([]);
	const [searching, setSearching] = useState(false);
	const [selected, setSelected] = useState<InviteCandidate[]>([]);
	const [inviting, setInviting] = useState(false);
	const [voting, setVoting] = useState(false);

	const acceptedMembers = useMemo(() => members.filter((member) => member.status === 'accepted'), [members]);
	const payingMembers = useMemo(() => members.filter((member) => member.role !== 'observer'), [members]);
	const memberEmails = useMemo(() => new Set(members.map((member) => member.email.toLowerCase())), [members]);
	const currentEmail = loginState?.email?.toLowerCase() ?? '';
	const myMember = useMemo(() => members.find((member) => currentEmail && member.email.toLowerCase() === currentEmail), [currentEmail, members]);
	const canVoteShares = myMember?.role !== 'observer';
	const ownerShare = acceptedMembers.find((member) => member.role === 'owner')?.share_percent ?? sharePercent;
	const alreadyShared = payingMembers.filter((member) => member.role !== 'owner').reduce((sum, member) => sum + member.share_percent, 0);
	const visibleAvatars = members.slice(0, 4);
	const extraAvatars = Math.max(members.length - visibleAvatars.length, 0);

	const memberAvatar = (person: { email: string; avatars?: string | null }) => {
		if (person.avatars) {
			return person.avatars;
		}

		if (currentEmail && person.email.toLowerCase() === currentEmail) {
			return loginState?.avatars ?? null;
		}

		return null;
	};

	useEffect(() => {
		const query = search.trim().toLowerCase();

		if (!searchOpen || !isEmailQuery(query)) {
			setResults([]);
			setSearching(false);
			return;
		}

		let cancelled = false;
		setSearching(true);

		const timer = window.setTimeout(async () => {
			try {
				const users = await basicUserSearchByEmail(query);
				if (!cancelled) {
					setResults(users.filter((user) => !memberEmails.has(user.email.toLowerCase())));
				}
			} catch {
				if (!cancelled) {
					setResults([]);
				}
			} finally {
				if (!cancelled) {
					setSearching(false);
				}
			}
		}, 300);

		return () => {
			cancelled = true;
			window.clearTimeout(timer);
		};
	}, [memberEmails, search, searchOpen]);

	const candidates = useMemo(() => {
		const query = search.trim().toLowerCase();
		const items: InviteCandidate[] = [...results];

		if (isEmailQuery(query) && !memberEmails.has(query) && !items.some((item) => item.email.toLowerCase() === query)) {
			items.push({ email: query });
		}

		return items;
	}, [memberEmails, results, search]);

	const toggleCandidate = (candidate: InviteCandidate) => {
		setSelected((prev) => {
			if (prev.some((item) => item.email.toLowerCase() === candidate.email.toLowerCase())) {
				return prev.filter((item) => item.email.toLowerCase() !== candidate.email.toLowerCase());
			}

			return [...prev, candidate];
		});
	};

	const onInviteSelected = () => {
		if (!selected.length || inviting) {
			return;
		}

		const people = selected.map((candidate) => ({
			email: candidate.email,
			name: memberName(candidate),
			initials: memberInitials(candidate),
			avatars: candidate.avatars,
		}));

		open(
			<Modal title={t('subscription.share-modal-title')} width="420px">
				<SubscriptionShareInviteModal
					people={people}
					subscriptionName={formatSubscriptionName(subscription.name, subscription.tariff, t)}
					price={subscription.price}
					currency={subscription.currency}
					ownerShare={ownerShare}
					alreadyShared={Math.max(0, alreadyShared)}
					onConfirm={async (shareEach, role) => {
						if (role !== 'observer' && (shareEach < 0 || shareEach * people.length > ownerShare + 0.001)) {
							notify.error(t('subscription.share-invalid'));
							return;
						}

						setInviting(true);

						try {
							for (const person of people) {
								await basicTrackedSubscriptionInvite(subscription.id, person.email, role === 'observer' ? 0 : shareEach, role);
							}

							notify.success(t('subscription.invite-sent'));
							setSelected([]);
							setSearch('');
							setSearchOpen(false);
							close();
							await onChanged();
						} finally {
							setInviting(false);
						}
					}}
				/>
			</Modal>
		);
	};

	const onRemove = async (member: SharedSubscriptionMember) => {
		if (!(await confirm('subscription.remove-member-confirm', 'subscription.remove-member'))) {
			return;
		}

		const message = await basicTrackedSubscriptionRemoveMember(subscription.id, member.id);
		notify.success(message || t('subscription.member-removed'));
		await onChanged();
	};

	const onVote = async (accept: boolean) => {
		if (!proposal) {
			return;
		}

		setVoting(true);

		try {
			const message = await basicTrackedSubscriptionVoteShares(subscription.id, proposal.id, accept);
			notify.success(message || t(accept ? 'subscription.shares-voted' : 'subscription.shares-rejected'));
			await onChanged();
		} finally {
			setVoting(false);
		}
	};

	const hasMembers = members.length > 1;
	const findButton = isOwner ? (
		<button type="button" onClick={() => setSearchOpen((open) => !open)} className="gu-glass-pill shrink-0 p-3 text-[13px] font-semibold gu-text-primary">
			{/* {t('subscription.find-members')} */}
			<div>
				<Plus fill="currentColor" size={19} />
			</div>
		</button>
	) : null;

	return (
		<section className="gu-glass-card relative space-y-4 overflow-hidden px-4 py-4">
			<div className="flex items-center justify-between gap-3">
				<h2 className="text-[15px] font-semibold gu-text-primary">{t('subscription.members-title')}</h2>
				{hasMembers ? <p className="text-[13px] gu-text-muted">{t('subscription.members-count', { count: members.length })}</p> : findButton}
			</div>

			{hasMembers ? (
				<div className="flex items-center gap-3">
					<div className="flex min-w-0 flex-1 items-center">
						{visibleAvatars.map((member, index) => (
							<div key={member.id} className={`rounded-full ring-2 ring-[var(--surface)] ${index > 0 ? '-ml-3' : ''}`} style={{ zIndex: visibleAvatars.length - index }}>
								<MemberAvatar initials={memberInitials(member)} src={memberAvatar(member)} size={40} />
							</div>
						))}
						{extraAvatars > 0 ? <span className="ml-2 text-[13px] font-semibold gu-text-muted">+{extraAvatars}</span> : null}
					</div>

					{findButton}
				</div>
			) : null}

			{searchOpen ? <GUInput type="email" autoFocus value={search} onChange={(event) => setSearch(event.target.value)} placeholder={t('subscription.search-email-placeholder')} /> : null}

			{searchOpen || hasMembers ? <div className="border-t border-dashed border-[var(--divider)]" /> : null}

			{searchOpen ? (
				<div className="space-y-3">
					<div className="flex items-center justify-between gap-3">
						<h3 className="text-[15px] font-semibold gu-text-primary">{t('subscription.search-results')}</h3>
						<p className="text-[13px] gu-text-muted">{t('subscription.members-count', { count: candidates.length })}</p>
					</div>

					{searching ? <p className="text-[13px] gu-text-muted">{t('action.loading')}</p> : null}

					{!searching && isEmailQuery(search.trim()) && candidates.length === 0 ? <p className="text-[13px] gu-text-muted">{t('subscription.search-empty')}</p> : null}

					<div className="space-y-1">
						{candidates.map((candidate) => {
							const checked = selected.some((item) => item.email.toLowerCase() === candidate.email.toLowerCase());

							return (
								<button
									key={candidate.email}
									type="button"
									onClick={() => toggleCandidate(candidate)}
									className="flex w-full items-center gap-3 rounded-2xl px-1 py-2 text-left transition hover:bg-[var(--surface-muted)]"
								>
									<MemberAvatar initials={memberInitials(candidate)} src={memberAvatar(candidate)} size={44} />
									<div className="min-w-0 flex-1">
										<p className="truncate text-[15px] font-semibold gu-text-primary">{memberName(candidate)}</p>
										<p className="truncate text-[13px] gu-text-muted">{candidate.email}</p>
									</div>
									<span className={`inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-full ${checked ? 'bg-[var(--text-primary)]' : 'bg-[var(--surface-muted)]'}`}>
										{checked ? <Check fill="var(--surface)" size={16} /> : <Plus fill="currentColor" size={16} />}
									</span>
								</button>
							);
						})}
					</div>
				</div>
			) : hasMembers ? (
				<div className="space-y-1">
					{members.map((member) => {
						const proposed = proposal?.items.find((item) => item.member_id === member.id)?.share_percent;
						// const isCurrentUser = Boolean(currentEmail && member.email.toLowerCase() === currentEmail);

						return (
							<div key={member.id} className="flex items-center gap-3 rounded-2xl px-1 py-2">
								<MemberAvatar initials={memberInitials(member)} src={memberAvatar(member)} size={44} />
								<div className="min-w-0 flex-1">
									<div className="flex min-w-0 items-center gap-2">
										<p className="truncate text-[15px] font-semibold gu-text-primary">{memberName(member)}</p>
										{member.role === 'owner' ? <AccentBadge className="shrink-0 px-2 py-0.5 rounded-full text-[13px]">{t('subscription.member-owner')}</AccentBadge> : null}
										{member.role === 'observer' ? <AccentBadge className="shrink-0 px-2 py-0.5 rounded-full text-[13px]">{t('subscription.member-observer')}</AccentBadge> : null}
										{/* {isCurrentUser ? <span className="shrink-0 text-[11px] gu-text-muted">{t('subscription.member-you')}</span> : null} */}
									</div>
									<p className="truncate text-[13px] gu-text-muted">
										{member.status === 'pending' ? `${t('subscription.member-pending')} · ${member.email}` : member.email}
										{proposed != null && proposed !== member.share_percent ? ` · ${t('subscription.proposed-share', { percent: proposed })}` : ''}
									</p>
								</div>

								<span className="shrink-0 text-[14px] font-semibold gu-text-primary">{member.role === 'observer' ? t('subscription.observer-share') : `${member.share_percent}%`}</span>

								{isOwner && member.role !== 'owner' ? (
									<button type="button" className="shrink-0 text-[12px] font-semibold text-[#0085FF]" onClick={() => onRemove(member)}>
										{t('subscription.remove-member')}
									</button>
								) : null}
							</div>
						);
					})}
				</div>
			) : null}

			{proposal && canVoteShares ? (
				<div className="space-y-2 rounded-2xl bg-[var(--surface-muted)] px-3 py-3">
					<p className="text-[13px] gu-text-muted">{t('subscription.share-proposal-pending')}</p>
					{proposal.my_vote == null ? (
						<div className="flex gap-2">
							<GUIButton type="button" variant="primary" disabled={voting} onClick={() => onVote(true)}>
								{t('subscription.approve-shares')}
							</GUIButton>
							<GUIButton type="button" disabled={voting} onClick={() => onVote(false)}>
								{t('action.cancel')}
							</GUIButton>
						</div>
					) : (
						<p className="text-[13px] gu-text-muted">{t('subscription.waiting-for-votes')}</p>
					)}
				</div>
			) : null}

			{selected.length > 0 ? (
				<button
					type="button"
					disabled={inviting}
					onClick={onInviteSelected}
					className="sticky bottom-2 flex w-full items-center justify-between gap-3 rounded-full bg-[var(--text-primary)] px-4 py-3 text-left text-[var(--surface)]"
				>
					<span className="text-[14px] font-semibold">{inviting ? t('action.loading') : t('subscription.send-invites', { count: selected.length })}</span>
					<span className="flex items-center">
						{selected.slice(0, 3).map((candidate, index) => (
							<span key={candidate.email} className={`overflow-hidden rounded-full ring-2 ring-[var(--text-primary)] ${index > 0 ? '-ml-2' : ''}`}>
								<MemberAvatar initials={memberInitials(candidate)} src={memberAvatar(candidate)} size={28} />
							</span>
						))}
					</span>
				</button>
			) : null}

			{isOwner && hasMembers ? (
				isPremium ? (
					<p className="text-[12px] gu-text-muted">{t('subscription.invite-unlimited')}</p>
				) : (
					<p className="text-[13px] text-center gu-text-muted">
						{t('subscription.invite-free-limit')}{' '}
						<Link to={ROUTES.PLANS} className="font-semibold text-[#0085FF] no-underline hover:no-underline">
							{t('subscription.invite-premium-link')}
						</Link>
					</p>
				)
			) : null}
		</section>
	);
};

export default SubscriptionSharingPanel;
