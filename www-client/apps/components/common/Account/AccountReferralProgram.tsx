import { useTranslation } from 'react-i18next';
import { notify } from '@/components/Notification/notify';
import Lock from '@/components/@icons/lock';
import ContentCopy from '@/components/@icons/content-copy';
import type { UserReferralRank, UserReferralResponse } from '@/rest/userAPI';
import { GUInput } from '@/components/ui/Input/GUInput';
import GUIButton from '@/components/ui/Button/GUIButton';
import { PREMIUM_GRADIENT } from '@/constants/designTokens';

interface AccountReferralProgramProps {
	referral: UserReferralResponse;
}

const AccountReferralProgram: React.FC<AccountReferralProgramProps> = ({ referral }) => {
	const { t } = useTranslation();
	const currentRank = referral.ranks.find((rank) => rank.current) ?? referral.ranks[0];

	const copyLink = async (url: string) => {
		try {
			await navigator.clipboard.writeText(url);
			notify.success(t('account.referral-copied'));
		} catch {
			notify.error(t('account.referral-copy-failed'));
		}
	};

	return (
		<div className="space-y-6">
			<section>
				<h2 className="text-[22px] font-bold leading-tight gu-text-primary">{t('account.referral-hero')}</h2>
				<p className="mt-2 text-[14px] leading-relaxed gu-text-muted">{t('account.referral-hero-desc')}</p>
			</section>

			<section className="flex flex-col gap-3">
				<ReferralLinkField value={referral.bot_url} label={t('account.referral-action-bot')} onCopy={() => void copyLink(referral.bot_url)} />
				<ReferralLinkField value={referral.site_url} label={t('account.referral-action-site')} onCopy={() => void copyLink(referral.site_url)} />
			</section>

			<section className="space-y-3">
				<h3 className="px-1 text-[13px] font-medium uppercase tracking-wide gu-text-muted">{t('account.referral-stats')}</h3>
				<div className="grid grid-cols-2 gap-3">
					<StatCard label={t('account.referral-count')} value={String(referral.referral_count)} />
					<StatCard label={t('account.referral-current-rank')} value={t(`account.referral-rank-${currentRank.level}-name`)} />
				</div>
				<StatCard label={t('account.referral-bonus')} value={bonusLabel(t, currentRank)} />
			</section>

			<section className="space-y-3">
				<h3 className="px-1 text-[13px] font-medium uppercase tracking-wide gu-text-muted">{t('account.referral-ranks')}</h3>
				<div className="flex flex-col gap-3">
					{referral.ranks.map((rank) => (
						<RankCard key={rank.level} rank={rank} currentLevel={currentRank.level} />
					))}
				</div>
			</section>
		</div>
	);
};

const ReferralLinkField: React.FC<{ value: string; label: string; onCopy: () => void }> = ({ value, label, onCopy }) => {
	const { t } = useTranslation();

	return (
		<div className="flex flex-col gap-1.5">
			<p className="px-1 text-[13px] font-medium uppercase tracking-wide gu-text-muted">{label}</p>
			<div className="relative flex">
				<GUInput value={value} readOnly onChange={() => undefined} aria-label={label} className="min-w-0 flex-1 cursor-default" />
				<GUIButton type="button" onClick={onCopy} className="absolute right-[2px] top-1/2 -translate-y-1/2 gu-glass-icon-btn shrink-0 !rounded-[12px]" aria-label={t('account.referral-copy')}>
					<div>
						<ContentCopy fill="currentColor" size={19} />
					</div>
				</GUIButton>
			</div>
		</div>
	);
};

const bonusLabel = (t: (key: string) => string, rank: UserReferralRank) => {
	if (rank.reward_days <= 0) {
		return t('account.referral-bonus-none');
	}

	return t(`account.referral-rank-${rank.level}-bonus`);
};

const StatCard: React.FC<{ label: string; value: string }> = ({ label, value }) => {
	return (
		<div className="gu-glass-card p-4">
			<p className="text-[12px] gu-text-muted">{label}</p>
			<p className="mt-2 text-[20px] font-semibold leading-tight gu-text-primary">{value}</p>
		</div>
	);
};

const RankCard: React.FC<{ rank: UserReferralRank; currentLevel: number }> = ({ rank, currentLevel }) => {
	const { t } = useTranslation();
	const locked = rank.level > currentLevel;

	return (
		<article
			className={`relative shrink-0 overflow-hidden rounded-[20px] p-4 ${rank.current ? 'shadow-[0_8px_24px_rgba(0,133,255,0.28)]' : 'gu-glass-card'}`}
			style={rank.current ? { background: PREMIUM_GRADIENT } : undefined}
		>
			{locked ? (
				<span className="absolute right-3 top-3 inline-flex h-8 w-8 items-center justify-center rounded-full bg-black/10 dark:bg-white/10">
					<Lock fill="currentColor" size={16} />
				</span>
			) : null}

			<p className={`text-[12px] font-medium ${rank.current ? 'text-white/85' : 'gu-text-muted'}`}>{rank.current ? `• ${t('account.referral-current')}` : t('account.referral-level', { level: rank.level })}</p>
			<p className={`mt-6 text-[20px] font-bold ${rank.current ? 'text-white' : 'gu-text-primary'}`}>{t(`account.referral-rank-${rank.level}-name`)}</p>
			<p className={`mt-1 text-[13px] ${rank.current ? 'text-white/85' : 'gu-text-muted'}`}>{t(`account.referral-rank-${rank.level}-range`)}</p>
			<p className={`mt-4 text-[13px] font-medium ${rank.current ? 'text-white' : 'gu-text-primary'}`}>{rank.reward_days > 0 ? t(`account.referral-rank-${rank.level}-bonus-short`) : t('account.referral-bonus-none')}</p>
		</article>
	);
};

export default AccountReferralProgram;
