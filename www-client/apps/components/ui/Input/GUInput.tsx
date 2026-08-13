import { forwardRef } from 'react';

export const inputClassName = 'gu-field';

interface GUInputProps extends React.InputHTMLAttributes<HTMLInputElement> {
	error?: string;
}

export const GUInput = forwardRef<HTMLInputElement, GUInputProps>(({ error, className = '', ...props }, ref) => {
	return (
		<div className="w-full">
			<input ref={ref} {...props} aria-invalid={!!error} className={`${inputClassName} ${error ? 'ring-2 ring-red-400' : ''} ${className}`.trim()} />

			{error && <span className="mt-1.5 block text-sm text-red-500">{error}</span>}
		</div>
	);
});

GUInput.displayName = 'GUInput';
