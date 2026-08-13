import { useState } from 'react';
import { getCurrencyFlagUrl } from '@/utils/currencyUtils';

interface CurrencyFlagProps {
	code: string;
	size?: 'sm' | 'md';
}

const CurrencyFlag = ({ code, size = 'md' }: CurrencyFlagProps) => {
	const flagUrl = getCurrencyFlagUrl(code);
	const [hasError, setHasError] = useState(false);

	const sizeClass = size === 'sm' ? 'h-6 w-6 text-[10px]' : 'h-8 w-8 text-[11px]';

	if (!flagUrl || hasError) {
		return <span className={`inline-flex shrink-0 items-center justify-center rounded-full bg-[#ececec] font-semibold text-[#555] ${sizeClass}`}>{code.slice(0, 2)}</span>;
	}

	return <img src={flagUrl} alt="" className={`shrink-0 rounded-full object-cover ${sizeClass}`} onError={() => setHasError(true)} />;
};

export default CurrencyFlag;
