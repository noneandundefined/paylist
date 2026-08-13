import { Helmet } from 'react-helmet-async';
import { useTranslation } from 'react-i18next';

interface PageMetaProps {
	titleKey: string;
	descriptionKey?: string;
}

const PageMeta: React.FC<PageMetaProps> = ({ titleKey, descriptionKey }) => {
	const { t } = useTranslation();

	return (
		<Helmet>
			<title>{`${t(titleKey)}`}</title>
			{descriptionKey && <meta name="description" content={t(descriptionKey)} />}
		</Helmet>
	);
};

export default PageMeta;
