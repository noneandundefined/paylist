export const SITE_ORIGIN = 'https://paylist.site';

export const getCanonicalUrl = (pathname: string): string => {
	const path = pathname.replace(/\/+$/, '') || '/';

	if (path === '/') {
		return `${SITE_ORIGIN}/`;
	}

	return `${SITE_ORIGIN}${path}`;
};
