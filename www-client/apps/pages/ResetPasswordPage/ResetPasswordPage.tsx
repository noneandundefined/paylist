import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { ROUTES } from '@/constants/constants';
import { useTranslation } from 'react-i18next';
import { useNavigate, useSearchParams } from 'react-router-dom';
import GUIButton from '@/components/ui/Button/GUIButton';
import { notify } from '@/components/Notification/notify';
import AuthPageLayout from '@/components/common/Auth/AuthPageLayout';
import InputPassword from '@/components/ui/Input/InputPassword';
import { basicAuthPasswordResetConfirm } from '@/rest/authAPI';
import { ValidationPasswordRequiredSchema } from '@/utils/ValidationSchema';

interface ResetPasswordForm {
	password: string;
	confirm_password: string;
}

const ResetPasswordPage = () => {
	const { t } = useTranslation();

	const navigate = useNavigate();
	const [searchParams] = useSearchParams();

	const uuid = searchParams.get('uuid') ?? '';
	const exp = searchParams.get('exp') ?? '';
	const sig = searchParams.get('sig') ?? '';
	const hasValidLink = Boolean(uuid && exp && sig);

	const {
		register,
		handleSubmit,
		watch,
		formState: { errors, isSubmitting },
	} = useForm<ResetPasswordForm>({
		mode: 'onChange',
		defaultValues: {
			password: '',
			confirm_password: '',
		},
	});

	const password = watch('password');

	useEffect(() => {
		if (hasValidLink) {
			return;
		}

		notify.error(t('auth.reset-invalid-link'));
		navigate(ROUTES.AUTH_FORGOT_PASSWORD, { replace: true });
	}, [hasValidLink, navigate, t]);

	const onSubmit = async (data: ResetPasswordForm) => {
		const message = await basicAuthPasswordResetConfirm(exp, sig, uuid, data.password);
		notify.success(message || t('auth.reset-success'));
		navigate(ROUTES.SIGNIN, { replace: true });
	};

	if (!hasValidLink) {
		return <main className="min-h-screen" />;
	}

	return (
		<AuthPageLayout title={t('label.page-reset-password')} subtitle={t('auth.reset-subtitle')}>
			<form className="space-y-5" onSubmit={handleSubmit(onSubmit)}>
				<InputPassword
					id="reset-password"
					autoComplete="new-password"
					placeholder={t('auth.new-password-placeholder')}
					{...register('password', ValidationPasswordRequiredSchema<ResetPasswordForm, 'password'>(t))}
					error={errors.password?.message}
				/>

				<InputPassword
					id="reset-confirm-password"
					autoComplete="new-password"
					placeholder={t('auth.repeat-new-password-placeholder')}
					{...register('confirm_password', {
						required: t('auth.error.confirm-password-required'),
						validate: (value) => value === password || t('auth.error.confirm-password-mismatch'),
					})}
					error={errors.confirm_password?.message}
				/>

				<GUIButton type="submit" variant="primary" isLoading={isSubmitting} loadingText={t('auth.reset-loading')}>
					{t('auth.reset-submit')}
				</GUIButton>
			</form>
		</AuthPageLayout>
	);
};

export default ResetPasswordPage;
