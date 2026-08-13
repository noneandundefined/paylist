import React, { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import GUIButton from '../ui/Button/GUIButton';

interface ModalConfirmProps {
	message: string;
	onConfirm: () => void;
}

const ModalConfirm: React.FC<ModalConfirmProps> = ({ message, onConfirm }) => {
	const { t } = useTranslation();

	useEffect(() => {
		const handler = (e: KeyboardEvent) => {
			if (e.key === 'Enter') {
				e.preventDefault();
				onConfirm();
			}
		};

		window.addEventListener('keydown', handler);
		return () => window.removeEventListener('keydown', handler);
	}, [onConfirm]);

	return (
		<React.Fragment>
			<div className="mt-4 sm:px-0 space-y-2">
				<p className="text-sm sm:text-base max-w-full sm:max-w-[40rem]">{message}</p>
			</div>

			<div className="mt-5 flex justify-end">
				<GUIButton variant="primary" onClick={onConfirm} isLoading={false}>
					{t('action.confirm')}
				</GUIButton>
			</div>
		</React.Fragment>
	);
};

export default ModalConfirm;
