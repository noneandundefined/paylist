import { useEffect } from 'react';
import { Helmet } from 'react-helmet-async';
import { useLocation } from 'react-router-dom';
import { getCanonicalUrl } from '@/utils/canonicalUrlUtils';
import { captureReferralCodeFromSearch } from '@/utils/referralCaptureUtils';

const PageCanonical = () => {
	const { pathname, search } = useLocation();
	const canonicalUrl = getCanonicalUrl(pathname);

	useEffect(() => {
		captureReferralCodeFromSearch(search);
	}, [search]);

	return (
		<Helmet>
			<link rel="canonical" href={canonicalUrl} />
			<meta property="og:url" content={canonicalUrl} />
		</Helmet>
	);
};

export default PageCanonical;
