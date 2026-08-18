import { useEffect, useState } from 'react';
import { getCountryFlagUrl } from '@/utils/countryUtils';
import { isImageFailed } from '@/utils/imageCacheUtils';
import RemoteImage from '@/components/ui/RemoteImage/RemoteImage';

interface CountryFlagProps {
	code: string;
	size?: 'sm' | 'md';
}

const CountryFlag = ({ code, size = 'md' }: CountryFlagProps) => {
	const flagUrl = getCountryFlagUrl(code);
	const [hasError, setHasError] = useState(() => isImageFailed(flagUrl));

	useEffect(() => {
		setHasError(isImageFailed(flagUrl));
	}, [flagUrl]);

	const sizeClass = size === 'sm' ? 'h-6 w-6 text-[10px]' : 'h-8 w-8 text-[11px]';

	if (!flagUrl || hasError) {
		return <span className={`inline-flex shrink-0 items-center justify-center rounded-full bg-[#ececec] font-semibold text-[#555] ${sizeClass}`}>{code.slice(0, 2)}</span>;
	}

	return (
		<span className={`relative inline-flex shrink-0 overflow-hidden rounded-full ${sizeClass}`}>
			<RemoteImage src={flagUrl} alt="" className="h-full w-full rounded-full object-cover" spinnerSize={size === 'sm' ? 10 : 12} onError={() => setHasError(true)} />
		</span>
	);
};

export default CountryFlag;
