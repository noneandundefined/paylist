import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { ROUTES } from '@/constants/constants';
import { CACHEKEYs } from '@/constants/CacheKeys.constants';
import { basicAuthSignUp } from '@/rest/authAPI';
import { Link, useNavigate } from 'react-router-dom';
import { GUInput } from '@/components/ui/Input/GUInput';
import GUIButton from '@/components/ui/Button/GUIButton';
import { notify } from '@/components/Notification/notify';
import AuthPageLayout from '@/components/common/Auth/AuthPageLayout';
import InputPassword from '@/components/ui/Input/InputPassword';
import { AuthSignupRequest } from '@/interface/auth/authSignupRequest.interface';
import { ValidationEmailSchema, ValidationPasswordSchema } from '@/utils/ValidationSchema';

const SignUpPage = () => {
	const { t } = useTranslation();

	const navigate = useNavigate();

	const {
		register,
		handleSubmit,
		formState: { errors, isSubmitting },
	} = useForm<AuthSignupRequest>({
		mode: 'onChange',
		defaultValues: {
			first_name: null,
			last_name: null,
			email: '',
			password: '',
		},
	});

	const onSubmit = async (data: AuthSignupRequest) => {
		const message = await basicAuthSignUp(data);

		notify.success(message || t('auth.confirm-email-sent'));

		const inviteToken = sessionStorage.getItem(CACHEKEYs.SUBSCRIPTION_INVITE_TOKEN);
		if (inviteToken) {
			navigate(`${ROUTES.SUBSCRIPTION_INVITE}?token=${encodeURIComponent(inviteToken)}`, { replace: true });
			return;
		}

		navigate(ROUTES.SIGNIN, { replace: true });
	};

	return (
		<AuthPageLayout
			title={t('label.page-signup')}
			subtitle={
				<>
					{t('auth.have-account')}{' '}
					<Link to={ROUTES.SIGNIN} className="gu-text-primary no-underline hover:no-underline">
						{t('action.sign-in')}
					</Link>
				</>
			}
		>
			<form className="space-y-5" onSubmit={handleSubmit(onSubmit)}>
				<div className="grid grid-cols-2 gap-3">
					<GUInput id="signup-first-name" autoComplete="given-name" {...register('first_name')} error={errors.first_name?.message} placeholder={t('auth.first-name-placeholder')} />
					<GUInput id="signup-last-name" autoComplete="family-name" {...register('last_name')} error={errors.last_name?.message} placeholder={t('auth.last-name-placeholder')} />
				</div>

				<GUInput id="signup-email" type="email" autoComplete="email" {...register('email', ValidationEmailSchema<AuthSignupRequest, 'email'>(t))} error={errors.email?.message} placeholder={t('auth.email-placeholder')} />

				<InputPassword
					id="signup-password"
					autoComplete="new-password"
					{...register('password', ValidationPasswordSchema<AuthSignupRequest, 'password'>(t))}
					error={errors.password?.message}
					placeholder={t('auth.password-placeholder')}
				/>

				<GUIButton type="submit" variant="primary" isLoading={isSubmitting} loadingText={t('auth.signup-loading')}>
					{t('auth.signup-action')}
				</GUIButton>
			</form>
		</AuthPageLayout>
	);
};

export default SignUpPage;
