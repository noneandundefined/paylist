import { useTranslation } from 'react-i18next';
import GUISelect from '@/components/ui/Select/GUISelect';
import { SUBSCRIPTION_TARIFF_NONE, SUBSCRIPTION_TARIFFS, type SubscriptionTariff } from '@/constants/subscriptionTariffs';

interface TariffSelectProps {
	value?: string;
	onChange: (tariff: SubscriptionTariff) => void;
	className?: string;
}

const TariffSelect: React.FC<TariffSelectProps> = ({ value, onChange, className }) => {
	const { t } = useTranslation();

	return (
		<GUISelect
			className={className}
			value={value || SUBSCRIPTION_TARIFF_NONE}
			onChange={(event) => onChange(event.target.value as SubscriptionTariff)}
			modalTitle={t('subscription.tariff-label')}
			aria-label={t('subscription.tariff-label')}
		>
			{SUBSCRIPTION_TARIFFS.map((tariff) => (
				<option key={tariff} value={tariff}>
					{t(`subscription.tariff-${tariff}`)}
				</option>
			))}
		</GUISelect>
	);
};

export default TariffSelect;
