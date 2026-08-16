import Navigation from '@/components/Navigation/Navigation';
import CookiePreferencesLink from '@/components/common/CookieConsent/CookiePreferencesLink';

interface PageLayoutProps {
	children: React.ReactNode;
}

const PageLayout: React.FC<PageLayoutProps> = ({ children }) => {
	return (
		<div className="flex min-h-screen min-w-0 flex-col">
			<div className="mx-auto flex w-full min-w-0 max-w-7xl flex-1 gap-5 px-4 py-4 lg:px-8">
				<main className="flex min-h-0 min-w-0 flex-1 flex-col pb-28">
					{children}
					<CookiePreferencesLink className="mt-8 self-center" />
				</main>

				<Navigation />
			</div>
		</div>
	);
};

export default PageLayout;
