import Navigation from '@/components/Navigation/Navigation';

interface PageLayoutProps {
	children: React.ReactNode;
}

const PageLayout: React.FC<PageLayoutProps> = ({ children }) => {
	return (
		<div className="flex min-h-screen min-w-0 flex-col">
			<div className="mx-auto flex w-full min-w-0 max-w-7xl flex-1 gap-5 px-4 py-4 lg:px-8">
				<main className="flex min-h-0 min-w-0 flex-1 flex-col pb-28">{children}</main>

				<Navigation />
			</div>
		</div>
	);
};

export default PageLayout;
