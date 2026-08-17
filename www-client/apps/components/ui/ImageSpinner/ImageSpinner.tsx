interface ImageSpinnerProps {
	size?: number;
	light?: boolean;
	className?: string;
}

const ImageSpinner: React.FC<ImageSpinnerProps> = ({ size = 16, light = false, className = '' }) => {
	const borderClass = light ? 'border-white/35 border-t-white' : 'border-slate-300/80 border-t-slate-700 dark:border-slate-600 dark:border-t-slate-200';

	return <span className={`inline-block shrink-0 animate-spin rounded-full border-2 ${borderClass} ${className}`} style={{ width: size, height: size }} aria-hidden="true" />;
};

export default ImageSpinner;
