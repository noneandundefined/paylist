import { useState } from 'react';
import Fallback from '@/components/Fallback/Fallback';

interface GUIButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
	onClick?: () => Promise<void> | void;
	type?: 'button' | 'submit' | 'reset';
	loadingText?: string;
	isLoading?: boolean;
	variant?: 'default' | 'primary';
}

const GUIButton: React.FC<GUIButtonProps> = ({ disabled, children, className = '', onClick, loadingText, isLoading, type = 'button', variant = 'default', ...props }) => {
	const [internalLoading, setInternalLoading] = useState(false);
	const loading = isLoading ?? internalLoading;
	const variantClassName = variant === 'primary' ? 'gu-btn-primary' : '';

	const runHandler = async () => {
		if (!onClick || loading) {
			return;
		}

		if (isLoading !== undefined) {
			await onClick();
			return;
		}

		try {
			setInternalLoading(true);
			await onClick();
		} finally {
			setInternalLoading(false);
		}
	};

	return (
		<>
			<button
				{...props}
				type={type}
				className={`${variantClassName} ${className} ${(loading || disabled) && '!cursor-not-allowed !opacity-40'}`}
				disabled={disabled || loading}
				onClick={
					onClick
						? (event) => {
								if (type === 'submit') {
									event.preventDefault();
								}

								void runHandler();
							}
						: undefined
				}
			>
				{children}
			</button>

			{loading && <Fallback text={loadingText} />}
		</>
	);
};

export default GUIButton;
