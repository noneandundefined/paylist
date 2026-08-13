import { Link } from 'react-router-dom';
import ChevronLeft from '@/components/@icons/chevron-left';
import Close from '@/components/@icons/close';

interface PageHeaderProps {
	title?: string;
	backTo?: string;
	backLabel?: string;
	onClose?: () => void;
	variant?: 'back' | 'close';
}

const PageHeader: React.FC<PageHeaderProps> = ({ title, backTo, backLabel, onClose, variant = 'back' }) => {
	if (variant === 'close' && onClose) {
		return (
			<div className="flex justify-end">
				<button type="button" onClick={onClose} className="gu-glass-icon-btn" aria-label={backLabel}>
					<Close fill="currentColor" size={22} />
				</button>
			</div>
		);
	}

	return (
		<header className="flex items-center">
			{backTo ? (
				<Link to={backTo} className="gu-glass-icon-btn no-underline hover:no-underline" aria-label={backLabel}>
					<ChevronLeft fill="currentColor" size={25} />
				</Link>
			) : (
				<span className="h-10 w-10" aria-hidden />
			)}

			<h1 className="flex-1 pr-10 text-center text-[17px] font-semibold gu-text-primary">{title}</h1>
		</header>
	);
};

export default PageHeader;
