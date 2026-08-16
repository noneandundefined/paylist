export const formatSavedPaymentMethodLabel = (type?: string | null, title?: string | null): string | null => {
	const rawTitle = title?.trim() ?? '';
	const rawType = type?.trim() ?? '';

	if (!rawTitle && !rawType) {
		return null;
	}

	const last4 = rawTitle.match(/(\d{4})\s*$/)?.[1];
	if (last4) {
		if (/\*{2,}/.test(rawTitle) && !/^bank card/i.test(rawTitle)) {
			return rawTitle.replace(/\s+/g, ' ');
		}

		return `${brandFromPaymentMethod(rawType, rawTitle)} **** ${last4}`;
	}

	if ((rawType === 'bank_card' && !rawTitle) || rawTitle === rawType) {
		return null;
	}

	return rawTitle || rawType;
};

const brandFromPaymentMethod = (type: string, title: string): string => {
	const haystack = `${type} ${title}`.toLowerCase();

	if (haystack.includes('mir')) {
		return 'MIR';
	}

	if (haystack.includes('visa')) {
		return 'Visa';
	}

	if (haystack.includes('master')) {
		return 'Mastercard';
	}

	if (haystack.includes('union')) {
		return 'UnionPay';
	}

	if (type === 'sbp' || haystack.includes('sbp')) {
		return 'СБП';
	}

	if (type === 'yoo_money' || haystack.includes('yoo')) {
		return 'ЮMoney';
	}

	return 'Card';
};
