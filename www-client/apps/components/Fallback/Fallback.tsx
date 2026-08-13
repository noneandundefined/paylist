import { useEffect } from 'react';
import { createPortal } from 'react-dom';
import { useTranslation } from 'react-i18next';

export interface FallbackProps {
	text?: string;
}

const Fallback = ({ text }: FallbackProps) => {
	const { t } = useTranslation();

	const label = text ?? t('message.page-loading');

	useEffect(() => {
		const previousOverflow = document.body.style.overflow;
		document.body.style.overflow = 'hidden';

		return () => {
			document.body.style.overflow = previousOverflow;
		};
	}, []);

	return createPortal(
		<div
			className="fixed inset-0 z-[10001] flex touch-none items-center justify-center backdrop-blur-[3px]"
			style={{ backgroundColor: 'color-mix(in srgb, var(--surface) 55%, transparent)' }}
			role="status"
			aria-live="polite"
			aria-busy="true"
			onPointerDown={(event) => event.preventDefault()}
		>
			<div className="pointer-events-none flex items-center gap-3">
				<span className="h-[22px] w-[22px] shrink-0 animate-spin rounded-full border-2 border-slate-300/80 border-t-slate-700 dark:border-slate-600 dark:border-t-slate-200" aria-hidden="true" />
				<p className="text-[15px] font-medium gu-text-primary">{label}</p>
			</div>
		</div>,
		document.body
	);
};

export default Fallback;
