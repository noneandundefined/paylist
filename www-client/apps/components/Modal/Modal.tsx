import Close from '../@icons/close';
import Tooltip from '../ui/Tooltip';
import { useTranslation } from 'react-i18next';
import { useEffect, useRef, useState } from 'react';
import { useModalContext } from '@/context/useModalContext';

interface ModalProps {
	title: string;
	argv?: number[] | string[];
	width?: string;
	panelClassName?: string;
	children: React.ReactNode;
}

/* Global component for create modal components */
const Modal: React.FC<ModalProps> = ({ title, argv, width = '600px', panelClassName = '', children }) => {
	const { t } = useTranslation();

	const { close } = useModalContext();

	const modalRef = useRef<HTMLDivElement>(null);

	const [animationClasses, setAnimationClasses] = useState('opacity-0 scale-95');

	useEffect(() => {
		const timer = setTimeout(() => {
			setAnimationClasses('opacity-100 scale-100');
		}, 1);

		return () => clearTimeout(timer);
	}, []);

	useEffect(() => {
		const handleClickOutside = (event: MouseEvent) => {
			if (modalRef.current && !modalRef.current.contains(event.target as Node)) {
				close();
			}
		};

		document.addEventListener('mousedown', handleClickOutside);
		return () => {
			document.removeEventListener('mousedown', handleClickOutside);
		};
	}, [close]);

	return (
		<div
			className={`fixed inset-0 z-[1005] flex items-center !justify-center backdrop-blur-[2px] transition-opacity duration-100 ${animationClasses.split(' ').find((cls) => cls.startsWith('opacity-'))}`}
			style={{ backgroundColor: 'var(--modal-overlay)' }}
		>
			<div
				ref={modalRef}
				style={{ width }}
				className={`gu-modal-surface relative max-w-[96%] rounded-lg p-3 shadow-lg sm:max-w-[90%] sm:p-4 transition-transform duration-100 ${animationClasses} ${panelClassName}`}
				role="dialog"
				aria-modal="true"
				aria-label={title || undefined}
			>
				{title.trim() && (
					<div className="mb-4 flex items-center justify-between">
						<div>
							<p className="text-left text-sm font-medium gu-text-primary sm:text-[15px]">
								{title} {argv}
							</p>
						</div>

						<Tooltip title={t('label.close')} position="bottom">
							<div className="cursor-pointer rounded-[8px] p-[8px] text-[1.1rem] transition hover:bg-[var(--surface-muted)] gu-text-secondary" onClick={close}>
								<Close fill="currentColor" size={17} />
							</div>
						</Tooltip>
					</div>
				)}

				{children}
			</div>
		</div>
	);
};

export default Modal;
