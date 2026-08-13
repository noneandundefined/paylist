interface SectionHeaderProps {
	title: string;
	titleClassName?: string;
}

const SectionHeader: React.FC<SectionHeaderProps> = ({ title, titleClassName = 'text-[14px] text-[#555] font-medium' }) => {
	return (
		<div className="flex items-center justify-between gap-3">
			<h2 className={titleClassName}>{title}</h2>
		</div>
	);
};

export default SectionHeader;
