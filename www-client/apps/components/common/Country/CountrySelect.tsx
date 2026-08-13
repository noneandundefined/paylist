import { useTranslation } from 'react-i18next';
import GUISelect from '@/components/ui/Select/GUISelect';

interface CountrySelectProps {
	value: string;
	onChange: React.ChangeEventHandler<HTMLSelectElement>;
	ariaLabel?: string;
}

const CountrySelect: React.FC<CountrySelectProps> = ({ value, onChange, ariaLabel }) => {
	const { t } = useTranslation();

	return (
		<GUISelect
			variant="country"
			value={value}
			onChange={onChange}
			modalTitle={t('account.country-select-title')}
			searchPlaceholder={t('account.search-countries')}
			emptyText={t('account.countries-not-found')}
			aria-label={ariaLabel ?? t('account.country')}
		/>
	);
};

export default CountrySelect;
