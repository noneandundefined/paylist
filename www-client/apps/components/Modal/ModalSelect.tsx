import { GUInput } from '../ui/Input/GUInput';
import { useTranslation } from 'react-i18next';
import { useEffect, useMemo, useState } from 'react';
import { useModalContext } from '@/context/useModalContext';
import CurrencyFlag from '@/components/common/CurrencyFlag';
import CountryFlag from '@/components/common/Country/CountryFlag';

import Check from '../@icons/check';

export type ModalSelectOption = {
	value: string | number;
	label: React.ReactNode;
	description?: React.ReactNode;
	disabled?: boolean;
	keywords?: string;
};

interface ModalSelectProps {
	options?: ModalSelectOption[];
	loadOptions?: () => Promise<ModalSelectOption[]>;
	value: string | number;
	onChange: (value: string | number) => void;
	searchPlaceholder?: string;
	emptyText?: string;
	loadingText?: string;
	errorText?: string;
	optionVariant?: 'default' | 'currency' | 'country';
	allowCustom?: boolean;
}

const getOptionSearchText = (option: ModalSelectOption): string => {
	if (option.keywords) {
		return option.keywords.toLowerCase();
	}

	if (typeof option.label === 'string' || typeof option.label === 'number') {
		return String(option.label).toLowerCase();
	}

	return String(option.value).toLowerCase();
};

const ModalSelect = ({ options: initialOptions, loadOptions, value, onChange, searchPlaceholder, emptyText, loadingText, errorText, optionVariant = 'default', allowCustom = false }: ModalSelectProps) => {
	const { t } = useTranslation();
	const { close } = useModalContext();

	const [search, setSearch] = useState('');
	const [options, setOptions] = useState<ModalSelectOption[]>(initialOptions ?? []);
	const [loading, setLoading] = useState(Boolean(loadOptions));
	const [error, setError] = useState(false);

	useEffect(() => {
		if (!loadOptions) {
			setOptions(initialOptions ?? []);
			return;
		}

		let cancelled = false;

		loadOptions()
			.then((items) => {
				if (!cancelled) {
					setOptions(items);
				}
			})
			.catch(() => {
				if (!cancelled) {
					setError(true);
				}
			})
			.finally(() => {
				if (!cancelled) {
					setLoading(false);
				}
			});

		return () => {
			cancelled = true;
		};
	}, [initialOptions, loadOptions]);

	const filteredOptions = useMemo(() => {
		const query = search.trim().toLowerCase();

		if (!query) {
			return options;
		}

		return options.filter((option) => getOptionSearchText(option).includes(query) || String(option.value).toLowerCase().includes(query));
	}, [options, search]);

	const visibleOptions = useMemo(() => {
		const query = search.trim();

		if (allowCustom && !query) {
			return [];
		}

		if (!allowCustom || !query || filteredOptions.length > 0) {
			return filteredOptions;
		}

		return [{ value: query, label: query }];
	}, [allowCustom, filteredOptions, search]);

	const handleSelect = (option: ModalSelectOption) => {
		if (option.disabled) {
			return;
		}

		onChange(option.value);
		close();
	};

	return (
		<div className="space-y-3">
			<GUInput type="text" value={search} onChange={(event) => setSearch(event.target.value)} placeholder={searchPlaceholder ?? t('label.search')} autoComplete="off" autoFocus />

			<ul className="max-h-[320px] space-y-3 overflow-y-auto">
				{loading && <li className="list-none py-6 text-center text-sm gu-text-muted">{loadingText ?? t('action.loading')}</li>}

				{error && !loading && <li className="list-none py-6 text-center text-sm text-red-500">{errorText ?? t('subscription.currencies-load-error')}</li>}

				{!loading && !error && visibleOptions.length === 0 && search.trim() !== '' && <li className="list-none py-6 text-center text-sm gu-text-muted">{emptyText ?? t('message.options-not-found')}</li>}

				{!loading &&
					!error &&
					visibleOptions.map((option) => {
						const isSelected = option.value === value;

						return (
							<li key={option.value} onClick={() => handleSelect(option)} className={`gu-modal-option ${option.disabled ? 'cursor-not-allowed opacity-50' : ''} ${isSelected ? 'gu-modal-option--selected' : ''}`}>
								{optionVariant === 'currency' && <CurrencyFlag code={String(option.value)} />}
								{optionVariant === 'country' && <CountryFlag code={String(option.value)} />}

								<div className="min-w-0 flex-1">
									<p className="truncate text-[15px] font-medium gu-text-primary">{option.label}</p>
									{option.description && <p className="text-[13px] gu-text-muted">{option.description}</p>}
								</div>

								{isSelected && (
									<span className="inline-flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-[#d7ff00]">
										<Check fill="#000" size={14} />
									</span>
								)}
							</li>
						);
					})}
			</ul>
		</div>
	);
};

export default ModalSelect;
