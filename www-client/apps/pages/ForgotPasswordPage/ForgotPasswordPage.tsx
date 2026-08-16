import { useForm } from 'react-hook-form';
import { ROUTES } from '@/constants/constants';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from 'react-router-dom';
import { GUInput } from '@/components/ui/Input/GUInput';
import GUIButton from '@/components/ui/Button/GUIButton';
import { notify } from '@/components/Notification/notify';
import AuthPageLayout from '@/components/common/Auth/AuthPageLayout';
import { basicAuthPasswordResetRequest } from '@/rest/authAPI';
import { ValidationEmailSchema } from '@/utils/ValidationSchema';

interface ForgotPasswordForm {
	email: string;
}

const ForgotPasswordPage = () => {
	const { t } = useTranslation();
	const navigate = useNavigate();

	const {
		register,
		handleSubmit,
		formState: { errors, isSubmitting },
	} = useForm<ForgotPasswordForm>({
		mode: 'onChange',
		defaultValues: {
			email: '',
		},
	});

	const onSubmit = async (data: ForgotPasswordForm) => {
		const message = await basicAuthPasswordResetRequest(data.email);
		notify.success(message || t('auth.forgot-sent'));
		navigate(ROUTES.SIGNIN, { replace: true });
	};

	return (
		<AuthPageLayout
			title={t('label.page-forgot-password')}
			subtitle={
				<>
					{t('auth.forgot-subtitle')}{' '}
					<Link to={ROUTES.SIGNIN} className="gu-text-primary no-underline hover:no-underline">
						{t('action.sign-in')}
					</Link>
				</>
			}
		>
			<form className="space-y-5" onSubmit={handleSubmit(onSubmit)}>
				<GUInput id="forgot-email" type="email" autoComplete="email" placeholder={t('auth.email-placeholder')} {...register('email', ValidationEmailSchema<ForgotPasswordForm, 'email'>(t))} error={errors.email?.message} />

				<GUIButton type="submit" variant="primary" isLoading={isSubmitting} loadingText={t('auth.forgot-loading')}>
					{t('auth.forgot-submit')}
				</GUIButton>
			</form>
		</AuthPageLayout>
	);
};

export default ForgotPasswordPage;
