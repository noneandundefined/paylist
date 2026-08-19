import { useMemo, useState } from 'react';
import { Navigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import PageLayout from '@/pages/PageLayout';
import Header from '@/components/Header/Header';
import PageLoadingState from '@/components/common/PageLoadingState/PageLoadingState';
import GUISelect from '@/components/ui/Select/GUISelect';
import { GUITextarea } from '@/components/ui/Input/GUITextarea';
import GUIButton from '@/components/ui/Button/GUIButton';
import { notify } from '@/components/Notification/notify';
import { QUERY_KEYS } from '@/constants/QueryKeys.constants';
import { ROUTES } from '@/constants/constants';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { useLoginState } from '@/hooks/useLoginState';
import { basicAdminRecipients, basicAdminSendMessage, type AdminMessageRecipient } from '@/rest/userAPI';

type AdminChannel = 'email' | 'telegram' | 'max';

const ALL_RECIPIENTS = 'all';

const recipientLabel = (user: AdminMessageRecipient): string => {
	const name = [user.first_name, user.last_name].filter(Boolean).join(' ').trim();
	return name ? `${name} (${user.email})` : user.email;
};

const AdminPage = () => {
	const { t } = useTranslation();
	const { isAdmin, loading: loginLoading } = useLoginState();
	const { data: recipients, loading: recipientsLoading } = useHandleServer([QUERY_KEYS.adminRecipients], () => basicAdminRecipients(), {
		enabled: isAdmin,
	});

	const [channel, setChannel] = useState<AdminChannel>('email');
	const [recipient, setRecipient] = useState(ALL_RECIPIENTS);
	const [text, setText] = useState('');
	const [sending, setSending] = useState(false);

	const users = useMemo(() => recipients ?? [], [recipients]);
	const selectedUser = users.find((user) => user.user_uuid === recipient);
	const channelReady = recipient === ALL_RECIPIENTS || (channel === 'email' ? true : channel === 'telegram' ? Boolean(selectedUser?.telegram_connected) : Boolean(selectedUser?.max_connected));

	if (loginLoading) {
		return <PageLoadingState />;
	}

	if (!isAdmin) {
		return <Navigate to={ROUTES.NOT_FOUND} replace />;
	}

	const onSend = async () => {
		const message = text.replace(/\r\n/g, '\n').replace(/\r/g, '\n');
		if (!message.trim() || sending) {
			return;
		}

		setSending(true);

		try {
			const result = await basicAdminSendMessage({
				channel,
				user_uuid: recipient === ALL_RECIPIENTS ? null : recipient,
				text: message,
			});

			notify.success(t('admin.send-result', { sent: result.sent, skipped: result.skipped, failed: result.failed }));
			setText('');
		} finally {
			setSending(false);
		}
	};

	return (
		<PageLayout>
			<Header />

			<section className="gu-glass-card space-y-4 px-4 py-4">
				<h1 className="text-[17px] font-semibold gu-text-primary">{t('admin.title')}</h1>
				<p className="text-[13px] gu-text-muted">{t('admin.subtitle')}</p>

				<div className="space-y-3">
					<div className="space-y-1.5">
						<p className="px-1 text-[13px] font-medium gu-text-muted">{t('admin.channel')}</p>
						<GUISelect value={channel} onChange={(event) => setChannel(event.target.value as AdminChannel)} modalTitle={t('admin.channel')} aria-label={t('admin.channel')}>
							<option value="email">{t('admin.channel-email')}</option>
							<option value="telegram">{t('admin.channel-telegram')}</option>
							<option value="max">{t('admin.channel-max')}</option>
						</GUISelect>
					</div>

					<div className="space-y-1.5">
						<p className="px-1 text-[13px] font-medium gu-text-muted">{t('admin.recipient')}</p>
						<GUISelect value={recipient} onChange={(event) => setRecipient(event.target.value)} modalTitle={t('admin.recipient')} searchPlaceholder={t('admin.recipient-search')} aria-label={t('admin.recipient')}>
							<option value={ALL_RECIPIENTS}>{t('admin.recipient-all')}</option>
							{users.map((user) => (
								<option key={user.user_uuid} value={user.user_uuid}>
									{recipientLabel(user)}
								</option>
							))}
						</GUISelect>
					</div>
				</div>

				<GUITextarea value={text} onChange={(event) => setText(event.target.value)} placeholder={t('admin.message-placeholder')} rows={8} />

				{!channelReady ? <p className="text-[13px] text-red-500">{t('admin.channel-not-connected')}</p> : null}

				<GUIButton type="button" variant="primary" disabled={!text.trim() || sending || recipientsLoading || !channelReady} isLoading={sending} loadingText={t('action.loading')} onClick={onSend}>
					{t('admin.send')}
				</GUIButton>
			</section>
		</PageLayout>
	);
};

export default AdminPage;
