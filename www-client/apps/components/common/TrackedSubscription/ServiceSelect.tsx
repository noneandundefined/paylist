import { useTranslation } from 'react-i18next';
import ChevronDown from '@/components/@icons/chevron-down';
import Modal from '@/components/Modal/Modal';
import ModalSelect from '@/components/Modal/ModalSelect';
import { useModalContext } from '@/context/useModalContext';
import { getServiceByName, loadServiceSelectOptions } from '@/rest/trackedSubscriptionAPI';

interface ServiceSelectProps {
	value: string;
	onChange: (name: string, category?: string) => void;
	error?: string;
}

const ServiceSelect: React.FC<ServiceSelectProps> = ({ value, onChange, error }) => {
	const { t } = useTranslation();
	const { open } = useModalContext();

	const openSelectModal = () => {
		open(
			<Modal title={t('subscription.service-select-title')} width="420px">
				<ModalSelect
					loadOptions={loadServiceSelectOptions}
					allowCustom
					value={value}
					onChange={(nextValue) => {
						const name = String(nextValue).trim();
						onChange(name, getServiceByName(name)?.category);
					}}
					searchPlaceholder={t('subscription.search-services')}
					emptyText={t('subscription.services-not-found')}
					errorText={t('subscription.services-load-error')}
				/>
			</Modal>
		);
	};

	return (
		<div className="w-full min-w-0">
			<div
				role="button"
				tabIndex={0}
				onClick={openSelectModal}
				onKeyDown={(event) => {
					if (event.key === 'Enter' || event.key === ' ') {
						event.preventDefault();
						openSelectModal();
					}
				}}
				aria-label={t('subscription.service-select-title')}
				aria-invalid={Boolean(error)}
				className={`gu-field cursor-pointer ${error ? 'ring-2 ring-red-400' : ''}`.trim()}
			>
				<div className="flex min-w-0 items-center justify-between gap-2">
					<span className={`min-w-0 flex-1 truncate font-normal ${value ? 'gu-text-primary' : 'gu-text-muted'}`}>{value || t('subscription.name-placeholder')}</span>
					<ChevronDown fill="currentColor" className="gu-text-muted shrink-0" size={19} />
				</div>
			</div>

			{error && <span className="mt-1.5 block text-sm text-red-500">{error}</span>}
		</div>
	);
};

export default ServiceSelect;
