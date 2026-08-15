import { useEffect } from 'react';
import { ROUTES } from '@/constants/constants';
import { basicAuthConfirmEmail } from '@/rest/authAPI';
import { CACHEKEYs } from '@/constants/CacheKeys.constants';
import { useNavigate, useSearchParams } from 'react-router-dom';

const ConfirmEmailPage = () => {
	const navigate = useNavigate();

	const [searchParams] = useSearchParams();

	useEffect(() => {
		const confirmEmail = async () => {
			try {
				const exp = searchParams.get('exp');
				const sig = searchParams.get('sig');
				const uuid = searchParams.get('uuid');

				if (!exp || !sig || !uuid) {
					navigate(ROUTES.SIGNIN, { replace: true });
					return;
				}

				const response = await basicAuthConfirmEmail(exp, sig, uuid);
				if (response.status === 'success') {
					localStorage.setItem(CACHEKEYs.L_SESSION, response.message);

					const inviteToken = sessionStorage.getItem(CACHEKEYs.SUBSCRIPTION_INVITE_TOKEN);
					if (inviteToken) {
						navigate(`${ROUTES.SUBSCRIPTION_INVITE}?token=${encodeURIComponent(inviteToken)}`, { replace: true });
						return;
					}

					navigate(ROUTES.HOME, { replace: true });
					return;
				}

				navigate(ROUTES.SIGNIN, { replace: true });
			} catch {
				navigate(ROUTES.SIGNIN, { replace: true });
			}
		};

		confirmEmail();
	}, [navigate, searchParams]);

	return <main className="min-h-screen flex items-center justify-center px-4"></main>;
};

export default ConfirmEmailPage;
