interface AccentBadgeProps {
	children: React.ReactNode;
	className?: string;
}

const AccentBadge: React.FC<AccentBadgeProps> = ({ children, className = '' }) => {
	return <span className={`gu-accent-badge ${className}`.trim()}>{children}</span>;
};

export default AccentBadge;
