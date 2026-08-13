import useIsMobile from '@/hooks/useIsMobile';
import { forwardRef, useEffect, useRef } from 'react';
import { inputClassName } from '@/components/ui/Input/GUInput';

const MAX_HEIGHT_PC = 300;
const MAX_HEIGHT_MOBILE = 200;

interface GUITextareaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
	error?: string;
}

export const GUITextarea = forwardRef<HTMLTextAreaElement, GUITextareaProps>(({ error, className, value, onChange, ...props }, ref) => {
	const isMobile = useIsMobile();
	const innerRef = useRef<HTMLTextAreaElement | null>(null);

	const MAX_HEIGHT = isMobile ? MAX_HEIGHT_MOBILE : MAX_HEIGHT_PC;

	const setRefs = (el: HTMLTextAreaElement) => {
		innerRef.current = el;

		if (typeof ref === 'function') {
			ref(el);
		} else if (ref) {
			ref.current = el;
		}
	};

	useEffect(() => {
		if (innerRef.current) {
			innerRef.current.style.height = 'auto';

			const newHeight = Math.min(innerRef.current.scrollHeight, MAX_HEIGHT);
			innerRef.current.style.height = newHeight + 'px';
		}
	}, [value, MAX_HEIGHT]);

	const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
		e.target.style.height = 'auto';

		const newHeight = Math.min(e.target.scrollHeight, MAX_HEIGHT);
		e.target.style.height = newHeight + 'px';

		onChange?.(e);
	};

	return (
		<div className="w-full">
			<textarea
				{...props}
				ref={setRefs}
				value={value}
				onChange={handleChange}
				aria-invalid={!!error}
				className={`${inputClassName} min-h-[96px] resize-none ${error ? 'ring-2 ring-red-400' : ''} ${className ?? ''} custom-scrollbar`.trim()}
			></textarea>

			{error && <span className="text-sm text-red-500">{error}</span>}
		</div>
	);
});
