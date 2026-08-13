export const config = {
	links: {
		URL_BACKEND_DEV: 'http://192.168.0.4:8080/api/v1',
		URL_FRONTEND_DEV: 'http://192.168.0.4:5173',
		URL_BACKEND_PROD: 'https://paylist.site/api/v1',
		URL_FRONTEND_PROD: 'https://paylist.site',
	},
	type: {
		release: 'dev' as 'dev' | 'prod',
	},
};
