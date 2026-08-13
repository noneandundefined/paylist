import Eye from '@/components/@icons/eye';
import EyeOff from '@/components/@icons/eye-off';
import { inputClassName } from '@/components/ui/Input/GUInput';
import { forwardRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

interface InputPasswordProps extends React.InputHTMLAttributes<HTMLInputElement> {
	label?: string;
	error?: string;
}

const InputPassword = forwardRef<HTMLInputElement, InputPasswordProps>(({ label, error, className = '', ...props }, ref) => {
	const { t } = useTranslation();
	const [showPassword, setShowPassword] = useState(false);

	return (
		<div className="w-full">
			{label && <label className="mb-1.5 block text-[14px] font-medium text-slate-700">{label}</label>}

			<div className="relative">
				<input ref={ref} {...props} aria-invalid={!!error} type={showPassword ? 'text' : 'password'} className={`${inputClassName} pr-11 ${error ? 'ring-2 ring-red-400' : ''} ${className}`.trim()} />

				<button
					type="button"
					aria-label={showPassword ? t('label.hide-password') : t('label.show-password')}
					onClick={() => setShowPassword((prev) => !prev)}
					className="absolute right-4 top-1/2 -translate-y-1/2 text-[#808080]"
				>
					{showPassword ? <EyeOff size={20} fill="currentColor" /> : <Eye size={20} fill="currentColor" />}
				</button>
			</div>

			{error && <span className="mt-1.5 block text-sm text-red-500">{error}</span>}
		</div>
	);
});

InputPassword.displayName = 'InputPassword';

export default InputPassword;
