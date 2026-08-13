import { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import type { CustomRouteConfig } from './config';

const RoutePage: React.FC<Pick<CustomRouteConfig, 'title' | 'component'>> = ({ title, component: Component }) => {
	const { t } = useTranslation();

	useEffect(() => {
		if (title) {
			document.title = `${t(title)}`;
			return;
		}
	}, [title, t]);

	return <Component />;
};

export default RoutePage;
