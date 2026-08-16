import { Helmet } from 'react-helmet-async';
import { useLocation } from 'react-router-dom';
import { getCanonicalUrl } from '@/utils/canonicalUrlUtils';

const PageCanonical = () => {
	const { pathname } = useLocation();
	const canonicalUrl = getCanonicalUrl(pathname);

	return (
		<Helmet>
			<link rel="canonical" href={canonicalUrl} />
			<meta property="og:url" content={canonicalUrl} />
		</Helmet>
	);
};

export default PageCanonical;
