import useIsMobile from '@/hooks/useIsMobile';
import { useTranslation } from 'react-i18next';
import { PREMIUM_GRADIENT_MINI } from '@/constants/designTokens';
import StarFilledIcon from '@/components/@icons/star-filled';

interface PremiumBadgeMiniProps {
	mobileView?: boolean;
}

const PremiumBadgeMini: React.FC<PremiumBadgeMiniProps> = ({ mobileView = false }) => {
	const { t } = useTranslation();

	const isMobile = useIsMobile();

	return (
		<div className="shrink-0">
			<span className="inline-flex items-center gap-1 rounded-full px-3 py-1.5 text-xs font-semibold uppercase tracking-wide text-white shadow-[0_4px_14px_rgba(0,133,255,0.35)]" style={{ background: PREMIUM_GRADIENT_MINI }}>
				<StarFilledIcon size={12} />
				{(!mobileView || !isMobile) && t('home.premium-badge')}
			</span>
		</div>
	);
};

export default PremiumBadgeMini;
