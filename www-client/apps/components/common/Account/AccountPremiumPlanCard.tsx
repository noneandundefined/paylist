import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { QUERY_KEYS } from '@/constants/QueryKeys.constants';
import { PREMIUM_GRADIENT } from '@/constants/designTokens';
import StarFilledIcon from '@/components/@icons/star-filled';
import GUIButton from '@/components/ui/Button/GUIButton';
import { useConfirm } from '@/hooks/useConfirm';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { basicPaymentAutoRenewUpdate, basicPaymentBilling } from '@/rest/paymentAPI';
import { formatPlanAmount } from '@/utils/planDisplayUtils';
import { formatSavedPaymentMethodLabel } from '@/utils/paymentMethodDisplayUtils';
import { formatSubscriptionDate } from '@/utils/TrackedSubscriptionDisplayUtils';

interface AccountPremiumPlanCardProps {
	validTo?: string | null;
	amount?: number;
	currency?: string;
}

const AccountPremiumPlanCard: React.FC<AccountPremiumPlanCardProps> = ({ validTo, amount, currency }) => {
	const { t, i18n } = useTranslation();
	const { confirm } = useConfirm();
	const [cancelling, setCancelling] = useState(false);

	const { data: billing, loading, updateHServer } = useHandleServer([QUERY_KEYS.paymentBilling], () => basicPaymentBilling());

	const formattedDate = validTo ? formatSubscriptionDate(validTo, i18n.language) : '';
	const autoRenewEnabled = billing?.auto_renew_enabled ?? false;
	const formattedAmount = typeof amount === 'number' && currency ? formatPlanAmount(amount, currency, i18n.language) : '';
	const cardLabel = formatSavedPaymentMethodLabel(billing?.payment_method_type, billing?.payment_method_title);

	const onCancel = async () => {
		if (!(await confirm('account.premium-cancel-confirm', 'account.premium-cancel'))) {
			return;
		}

		setCancelling(true);

		try {
			const next = await basicPaymentAutoRenewUpdate({ enabled: false });
			updateHServer(next);
		} finally {
			setCancelling(false);
		}
	};

	return (
		<div className="flex flex-col md:flex-row justify-between md:items-center gap-3 relative overflow-hidden rounded-[20px] px-5 py-4 shadow-[0_8px_24px_rgba(0,133,255,0.28)]" style={{ background: PREMIUM_GRADIENT }}>
			<div>
				<p className="flex flex-wrap items-center gap-2 text-[17px] font-semibold text-white">
					<span className="inline-flex items-center gap-1 rounded-full bg-white px-2 py-0.5 text-[11px] font-bold uppercase tracking-wide text-[#0085FF]">
						<StarFilledIcon size={10} />
						{t('home.premium-badge')}
					</span>
				</p>

				{formattedDate ? (
					<p className="mt-2 text-[13px] text-white/90">
						{autoRenewEnabled
							? formattedAmount
								? t('account.premium-renews', { amount: formattedAmount, date: formattedDate })
								: t('account.premium-ends', { date: formattedDate })
							: t('account.premium-cancelled', { date: formattedDate })}
					</p>
				) : null}

				{cardLabel ? <p className="mt-2 text-[12px] tracking-[0.04em] text-white/80">{cardLabel}</p> : null}
			</div>

			{!loading && autoRenewEnabled ? (
				<GUIButton type="button" isLoading={false} disabled={cancelling} onClick={onCancel} className="rounded-full border border-white/35 bg-white/20 px-4 py-2 text-[14px] font-semibold text-white backdrop-blur-sm">
					{t('account.premium-cancel')}
				</GUIButton>
			) : null}
		</div>
	);
};

export default AccountPremiumPlanCard;
