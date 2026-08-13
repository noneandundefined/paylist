import { useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import Modal from '@/components/Modal/Modal';
import ModalConfirm from '@/components/Modal/ModalConfirm';
import { useModalContext } from '@/context/useModalContext';

interface ConfirmModalBodyProps {
	message: string;
	onConfirm: () => void;
	onDismiss: () => void;
}

const ConfirmModalBody = ({ message, onConfirm, onDismiss }: ConfirmModalBodyProps) => {
	useEffect(() => () => onDismiss(), [onDismiss]);

	return <ModalConfirm message={message} onConfirm={onConfirm} />;
};

export const useConfirm = () => {
	const { t } = useTranslation();
	const { open, close } = useModalContext();

	const confirm = useCallback(
		(messageKey: string, titleKey?: string): Promise<boolean> =>
			new Promise((resolve) => {
				let settled = false;

				const settle = (value: boolean) => {
					if (settled) {
						return;
					}

					settled = true;
					resolve(value);
				};

				open(
					<Modal title={titleKey ? t(titleKey) : ''}>
						<ConfirmModalBody
							message={t(messageKey)}
							onConfirm={() => {
								close();
								settle(true);
							}}
							onDismiss={() => settle(false)}
						/>
					</Modal>
				);
			}),
		[open, close, t]
	);

	return { confirm };
};
