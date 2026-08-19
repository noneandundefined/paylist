import PageLayout from '@/pages/PageLayout';
import { ROUTES } from '@/constants/constants';
import { useTranslation } from 'react-i18next';
import { QUERY_KEYS } from '@/constants/QueryKeys.constants';
import PageHeader from '@/components/common/PageHeader/PageHeader';
import PageLoadingState from '@/components/common/PageLoadingState/PageLoadingState';
import AccountReferralProgram from '@/components/common/Account/AccountReferralProgram';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { basicUserReferralGet } from '@/rest/userAPI';

const ReferralPage = () => {
	const { t } = useTranslation();
	const { data: referral, loading } = useHandleServer([QUERY_KEYS.userReferral], () => basicUserReferralGet());

	if (loading) {
		return <PageLoadingState />;
	}

	return (
		<PageLayout>
			<div className="mx-auto flex w-full flex-col space-y-5">
				<PageHeader title={t('label.page-referrals')} backTo={ROUTES.ACCOUNT} backLabel={t('action.back')} />
				{referral ? <AccountReferralProgram referral={referral} /> : null}
			</div>
		</PageLayout>
	);
};

export default ReferralPage;
