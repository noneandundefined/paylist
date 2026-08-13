import { useTranslation } from 'react-i18next';
import GUISelect from '@/components/ui/Select/GUISelect';

interface CurrencySelectProps {
	value: string;
	onChange: React.ChangeEventHandler<HTMLSelectElement>;
	ariaLabel?: string;
}

const CurrencySelect: React.FC<CurrencySelectProps> = ({ value, onChange, ariaLabel }) => {
	const { t } = useTranslation();

	return (
		<GUISelect
			variant="currency"
			value={value}
			onChange={onChange}
			modalTitle={t('subscription.currency-select-title')}
			searchPlaceholder={t('subscription.search-currencies')}
			emptyText={t('subscription.currencies-not-found')}
			aria-label={ariaLabel ?? t('account.currency')}
		/>
	);
};

export default CurrencySelect;
