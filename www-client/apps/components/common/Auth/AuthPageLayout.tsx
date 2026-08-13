import AuthTermsFooter from '@/components/common/Auth/AuthTermsFooter';

interface AuthPageLayoutProps {
	title: string;
	subtitle: React.ReactNode;
	children: React.ReactNode;
}

const AuthPageLayout: React.FC<AuthPageLayoutProps> = ({ title, subtitle, children }) => {
	return (
		<div className="flex min-h-screen justify-center gu-page-bg px-6 py-10">
			<div className="flex w-full max-w-md flex-col justify-center space-y-5">
				<div className="space-y-2">
					<h1 className="text-[30px] font-bold gu-text-primary">{title}</h1>
					<p className="text-[15px] font-normal gu-text-muted">{subtitle}</p>
				</div>

				{children}

				<AuthTermsFooter />
			</div>
		</div>
	);
};

export default AuthPageLayout;
