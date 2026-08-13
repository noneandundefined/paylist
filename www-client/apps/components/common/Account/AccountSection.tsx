interface AccountSectionProps {
	title: string;
	children: React.ReactNode;
}

const AccountSection: React.FC<AccountSectionProps> = ({ title, children }) => {
	return (
		<section className="space-y-2">
			<h2 className="px-1 text-[13px] font-medium uppercase tracking-wide gu-text-muted">{title}</h2>
			{children}
		</section>
	);
};

export default AccountSection;
