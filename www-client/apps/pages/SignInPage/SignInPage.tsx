import { useForm } from 'react-hook-form';
import { ROUTES } from '@/constants/constants';
import { useTranslation } from 'react-i18next';
import { basicAuthSignIn } from '@/rest/authAPI';
import { Link, useNavigate } from 'react-router-dom';
import { GUInput } from '@/components/ui/Input/GUInput';
import GUIButton from '@/components/ui/Button/GUIButton';
import { notify } from '@/components/Notification/notify';
import AuthPageLayout from '@/components/common/Auth/AuthPageLayout';
import InputPassword from '@/components/ui/Input/InputPassword';
import { AuthSigninRequest } from '@/interface/auth/authSigninRequest.interface';
import { ValidationEmailSchema, ValidationPasswordRequiredSchema } from '@/utils/ValidationSchema';

const SignInPage = () => {
	const { t } = useTranslation();

	const navigate = useNavigate();

	const {
		register,
		handleSubmit,
		formState: { errors, isSubmitting },
	} = useForm<AuthSigninRequest>({
		mode: 'onChange',
		defaultValues: {
			email: '',
			password: '',
		},
	});

	const onSubmit = async (data: AuthSigninRequest) => {
		const result = await basicAuthSignIn(data);

		if (result.status === 'signed_in') {
			notify.success(t('auth.signin-success'));
			navigate(ROUTES.HOME, { replace: true });
			return;
		}

		if (result.status === 'sent') {
			notify.success(result.message || t('auth.confirm-email-sent'));
			return;
		}

		if (result.status === 'password') {
			notify.error(result.message);
		}
	};

	return (
		<AuthPageLayout
			title={t('label.page-signin')}
			subtitle={
				<>
					{t('auth.new-user-prompt')}{' '}
					<Link to={ROUTES.SIGNUP} className="gu-text-primary no-underline hover:no-underline">
						{t('auth.create-account')}
					</Link>
				</>
			}
		>
			<form className="space-y-5" onSubmit={handleSubmit(onSubmit)}>
				<GUInput id="signin-email" type="email" autoComplete="email" placeholder={t('auth.email-placeholder')} {...register('email', ValidationEmailSchema<AuthSigninRequest, 'email'>(t))} error={errors.email?.message} />

				<InputPassword
					id="signin-password"
					autoComplete="current-password"
					placeholder={t('auth.password-placeholder')}
					{...register('password', ValidationPasswordRequiredSchema<AuthSigninRequest, 'password'>(t))}
					error={errors.password?.message}
				/>

				<Link to="#" className="inline-block text-[14px] font-bold gu-text-primary no-underline hover:no-underline">
					{t('auth.forgot-password')}
				</Link>

				<GUIButton type="submit" variant="primary" isLoading={isSubmitting} loadingText={t('auth.signin-loading')}>
					{t('auth.login-button')}
				</GUIButton>
			</form>
		</AuthPageLayout>
	);
};

export default SignInPage;
