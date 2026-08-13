import PageLayout from '@/pages/PageLayout';
import Fallback from '@/components/Fallback/Fallback';
import { useTranslation } from 'react-i18next';

const PageLoadingState: React.FC = () => {
	const { t } = useTranslation();

	return (
		<PageLayout>
			<Fallback text={t('message.page-loading')} />
		</PageLayout>
	);
};

export default PageLoadingState;
