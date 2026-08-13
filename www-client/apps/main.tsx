import '@/utils/i18n';
import '@/utils/theme';

import { Suspense } from 'react';
import { createRoot } from 'react-dom/client';
import 'react-toastify/dist/ReactToastify.css';
import { BrowserRouter } from 'react-router-dom';
import { HelmetProvider } from 'react-helmet-async';
import { ErrorBoundary } from 'react-error-boundary';
import NotificationProvider from './components/Notification/NotificationProvider';

import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import Router from '@/router';
import Fallback from '@/components/Fallback/Fallback';
import { ThemeProvider } from './context/ThemeContext';
import { ModalProvider } from './context/useModalContext';
import ErrorFallback from '@/components/Fallback/ErrorFallback';

const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			staleTime: 1000 * 30,
			refetchOnWindowFocus: false,
			refetchOnReconnect: true,
			refetchOnMount: true,
		},
	},
});

const App = () => {
	return (
		<ErrorBoundary fallback={<ErrorFallback />}>
			<Suspense fallback={<Fallback />}>
				<Router />
			</Suspense>
		</ErrorBoundary>
	);
};

const Root = () => {
	return (
		<>
			<NotificationProvider />

			<App />
		</>
	);
};

createRoot(document.getElementById('root')!).render(
	<BrowserRouter>
		<HelmetProvider>
			<QueryClientProvider client={queryClient}>
				<ThemeProvider>
					<ModalProvider>
						<Root />

						<ReactQueryDevtools initialIsOpen={true} />
					</ModalProvider>
				</ThemeProvider>
			</QueryClientProvider>
		</HelmetProvider>
	</BrowserRouter>
);
