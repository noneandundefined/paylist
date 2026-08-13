import { useTranslation } from 'react-i18next';

const ErrorFallback = () => {
	const { t } = useTranslation();

	return (
		<div className="flex min-h-screen flex-col items-center justify-center bg-gray-50 px-4 text-center">
			<h1 className="text-xl font-semibold text-gray-900">{t('message.error-fallback-title')}</h1>
			<p className="mt-3 max-w-md text-sm text-gray-600">{t('message.error-fallback-description')}</p>
		</div>
	);
};

export default ErrorFallback;
