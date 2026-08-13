import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { ROUTES } from '@/constants/constants';
import { PREMIUM_GRADIENT } from '@/constants/designTokens';
import StarFilledIcon from '@/components/@icons/star-filled';

const PremiumBadge = () => {
	const { t } = useTranslation();

	return (
		<div className="relative overflow-hidden rounded-[20px] px-5 py-4 shadow-[0_8px_24px_rgba(0,133,255,0.28)]" style={{ background: PREMIUM_GRADIENT }}>
			<div className="flex items-center justify-between gap-3">
				<div className="min-w-0">
					<p className="flex flex-wrap items-center gap-2 text-[17px] font-semibold text-white">
						<span>{t('account.become-pro-prefix')}</span>
						<span className="inline-flex items-center gap-1 rounded-full bg-white px-2 py-0.5 text-[11px] font-bold uppercase tracking-wide text-[#0085FF]">
							<StarFilledIcon size={10} />
							{t('home.premium-badge')}
						</span>
					</p>
					<p className="mt-1 text-[13px] text-white/85">{t('account.become-pro-subtitle')}</p>
				</div>

				<Link
					to={ROUTES.PLANS}
					className="shrink-0 rounded-full border border-white/35 bg-white/20 px-4 py-2 text-[14px] font-semibold text-white no-underline backdrop-blur-sm transition hover:bg-white/30 hover:no-underline"
				>
					{t('account.choose-plan')}
				</Link>
			</div>
		</div>
	);
};

export default PremiumBadge;
