import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { ROUTES } from '@/constants/constants';

const NotFoundPage = () => {
	const { t } = useTranslation();

	return (
		<div className="flex flex-col min-h-screen items-center justify-center bg-gray-50">
			<h1 className="text-[20vw] font-bold text-black">404</h1>
			<Link to={ROUTES.HOME} className="inline-block font-medium text-blue-600">
				{t('action.go-home')}
			</Link>
		</div>
	);
};

export default NotFoundPage;
