import { useTranslation } from 'react-i18next';
import { PREMIUM_GRADIENT } from '@/constants/designTokens';
import StarFilledIcon from '@/components/@icons/star-filled';
import { formatSubscriptionDate } from '@/utils/TrackedSubscriptionDisplayUtils';

interface AccountPremiumPlanCardProps {
	validTo?: string | null;
	amount?: number;
	currency?: string;
}

const AccountPremiumPlanCard: React.FC<AccountPremiumPlanCardProps> = ({ validTo }) => {
	const { t, i18n } = useTranslation();
	const formattedDate = validTo ? formatSubscriptionDate(validTo, i18n.language) : '';

	return (
		<div className="flex flex-col md:flex-row justify-between md:items-center gap-3 relative overflow-hidden rounded-[20px] px-5 py-4 shadow-[0_8px_24px_rgba(0,133,255,0.28)]" style={{ background: PREMIUM_GRADIENT }}>
			<div>
				<p className="flex flex-wrap items-center gap-2 text-[17px] font-semibold text-white">
					<span className="inline-flex items-center gap-1 rounded-full bg-white px-2 py-0.5 text-[11px] font-bold uppercase tracking-wide text-[#0085FF]">
						<StarFilledIcon size={10} />
						{t('home.premium-badge')}
					</span>
				</p>

				{formattedDate ? <p className="mt-2 text-[13px] text-white/90">{t('account.premium-ends', { date: formattedDate })}</p> : null}
			</div>
		</div>
	);
};

export default AccountPremiumPlanCard;
