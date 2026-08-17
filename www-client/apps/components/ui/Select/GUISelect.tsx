import Modal from '@/components/Modal/Modal';
import ChevronDown from '@/components/@icons/chevron-down';
import { useTranslation } from 'react-i18next';
import { useModalContext } from '@/context/useModalContext';
import CurrencyFlag from '@/components/common/CurrencyFlag';
import { loadCurrencySelectOptions } from '@/rest/currencyAPI';
import { loadCountrySelectOptions } from '@/rest/countryAPI';
import { getCountryLabel } from '@/utils/countryDisplayUtils';
import React, { useCallback, useEffect, useRef, useState } from 'react';
import ModalSelect, { type ModalSelectOption } from '@/components/Modal/ModalSelect';
import CountryFlag from '@/components/common/Country/CountryFlag';

interface GUISelectProps extends Omit<React.HTMLAttributes<HTMLDivElement>, 'onChange'> {
	children?: React.ReactNode;
	value?: string | number;
	onChange?: React.ChangeEventHandler<HTMLSelectElement>;
	placeholder?: string;
	modalTitle?: string;
	searchPlaceholder?: string;
	emptyText?: string;
	variant?: 'default' | 'currency' | 'country';
}

const GUISelect: React.FC<GUISelectProps> = ({ children, value, onChange, placeholder = '', modalTitle, searchPlaceholder, emptyText, variant = 'default', className = '', ...rest }) => {
	const { t } = useTranslation();
	const { open } = useModalContext();
	const ariaLabel = rest['aria-label'];

	const [internalValue, setInternalValue] = useState<string | number>();

	const ref = useRef<HTMLDivElement>(null);
	const onChangeRef = useRef(onChange);

	useEffect(() => {
		onChangeRef.current = onChange;
	}, [onChange]);

	const options: ModalSelectOption[] = React.Children.toArray(children ?? [])
		.filter((child): child is React.ReactElement => React.isValidElement(child) && child.type === 'option')
		.map((child) => ({
			value: child.props.value,
			label: child.props.children,
			disabled: child.props.disabled,
		}));

	const selectedValue = value ?? internalValue ?? options[0]?.value;
	const selected = options.find((option) => option.value === selectedValue);

	const triggerChange = useCallback((val: string | number) => {
		const nextValue = String(val);

		setInternalValue(nextValue);

		onChangeRef.current?.({
			target: { value: nextValue },
		} as React.ChangeEvent<HTMLSelectElement>);
	}, []);
	const openSelectModal = () => {
		if (variant === 'currency') {
			open(
				<Modal title={modalTitle ?? placeholder} width="420px">
					<ModalSelect loadOptions={loadCurrencySelectOptions} optionVariant="currency" value={selectedValue ?? 'USD'} onChange={triggerChange} searchPlaceholder={searchPlaceholder} emptyText={emptyText} />
				</Modal>
			);
			return;
		}

		if (variant === 'country') {
			open(
				<Modal title={modalTitle ?? placeholder} width="420px">
					<ModalSelect loadOptions={loadCountrySelectOptions} optionVariant="country" value={selectedValue ?? 'US'} onChange={triggerChange} searchPlaceholder={searchPlaceholder} emptyText={emptyText} />
				</Modal>
			);
			return;
		}

		if (selectedValue === undefined) return;

		open(
			<Modal title={modalTitle ?? placeholder} width="420px">
				<ModalSelect options={options} value={selectedValue} onChange={triggerChange} searchPlaceholder={searchPlaceholder} emptyText={emptyText} />
			</Modal>
		);
	};

	const triggerLabel =
		variant === 'currency' ? (
			<span className="flex items-center gap-2.5 font-normal">
				<CurrencyFlag code={String(selectedValue ?? 'USD')} size="sm" />
				<span>{selectedValue ?? placeholder}</span>
			</span>
		) : variant === 'country' ? (
			<span className="flex items-center gap-2.5 font-normal">
				<CountryFlag code={String(selectedValue ?? 'US')} size="sm" />
				<span>{getCountryLabel(String(selectedValue ?? 'US'), String(selectedValue ?? 'US'), t)}</span>
			</span>
		) : (
			<span className="min-w-0 flex-1 truncate font-normal">{selected?.label || placeholder}</span>
		);

	return (
		<div
			ref={ref}
			role="button"
			tabIndex={0}
			onClick={openSelectModal}
			onKeyDown={(event) => {
				if (event.key === 'Enter' || event.key === ' ') {
					event.preventDefault();
					openSelectModal();
				}
			}}
			aria-label={typeof ariaLabel === 'string' ? ariaLabel : undefined}
			className={`gu-field cursor-pointer ${className}`}
		>
			<div className="flex min-w-0 items-center justify-between gap-2">
				{triggerLabel}

				<ChevronDown fill="currentColor" className="gu-text-muted shrink-0" size={19} />
			</div>
		</div>
	);
};

export default GUISelect;
