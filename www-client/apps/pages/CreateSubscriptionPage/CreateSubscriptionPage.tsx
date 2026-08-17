import PageLayout from '../PageLayout';
import { useForm } from 'react-hook-form';
import { ROUTES } from '@/constants/constants';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { GUInput } from '@/components/ui/Input/GUInput';
import GUIButton from '@/components/ui/Button/GUIButton';
import { notify } from '@/components/Notification/notify';
import GUISelect from '@/components/ui/Select/GUISelect';
import PageHeader from '@/components/common/PageHeader/PageHeader';
import CurrencySelect from '@/components/common/Currency/CurrencySelect';
import CategoryChipGroup from '@/components/common/Category/CategoryChipGroup';
import SubscriptionSettingsPanel from '@/components/common/TrackedSubscription/SubscriptionSettingsPanel';
import TariffSelect from '@/components/common/TrackedSubscription/TariffSelect';
import ServiceSelect from '@/components/common/TrackedSubscription/ServiceSelect';
import { SUBSCRIPTION_TARIFF_NONE, type SubscriptionTariff } from '@/constants/subscriptionTariffs';
import { usePremiumFeatureGate } from '@/hooks/usePremiumFeatureGate';
import { useSubscriptionCategories } from '@/hooks/useSubscriptionCategories';
import { useInvalidateSubscriptions } from '@/hooks/useInvalidateSubscriptions';
import { notifyPremiumRequired } from '@/utils/premiumUtils';
import { resolvePremiumSubscriptionFlags } from '@/utils/subscriptionPremiumUtils';
import { isCustomCategory } from '@/utils/categoryDisplayUtils';
import { validateFutureDate } from '@/utils/subscriptionFormUtils';
import { basicTrackedSubscriptionCreate } from '@/rest/trackedSubscriptionAPI';
import type { TrackedSubscriptionPeriod } from '@/interface/trackedSubscription/trackedSubscriptionCreateRequest.interface';

interface CreateSubscriptionForm {
	name: string;
	tariff: SubscriptionTariff;
	priceInput: string;
	currency: string;
	period: TrackedSubscriptionPeriod;
	date_pay: string;
	auto_renewal: boolean;
	notification: boolean;
	include_in_analytics: boolean;
	categories: string[];
}

const CreateSubscriptionPage = () => {
	const { t } = useTranslation();

	const navigate = useNavigate();
	const { invalidateListAndSummary } = useInvalidateSubscriptions();
	const { canUseNotification, isPremium } = usePremiumFeatureGate();
	const { categories } = useSubscriptionCategories();

	const {
		register,
		handleSubmit,
		watch,
		setValue,
		formState: { errors, isSubmitting },
	} = useForm<CreateSubscriptionForm>({
		mode: 'onChange',
		defaultValues: {
			name: '',
			tariff: SUBSCRIPTION_TARIFF_NONE,
			priceInput: '',
			currency: 'USD',
			period: 'monthly',
			date_pay: '',
			auto_renewal: true,
			notification: false,
			include_in_analytics: true,
			categories: [],
		},
	});

	const selectedCategories = watch('categories') ?? [];
	const autoRenewal = watch('auto_renewal');
	const notificationEnabled = watch('notification');
	const includeInAnalytics = watch('include_in_analytics');

	const toggleCategory = (slug: string) => {
		const category = categories.find((item) => item.slug === slug);

		if (!isPremium && category && isCustomCategory(category)) {
			notifyPremiumRequired(t);
			return;
		}

		const next = selectedCategories.includes(slug) ? selectedCategories.filter((item) => item !== slug) : [...selectedCategories, slug];

		setValue('categories', next, { shouldDirty: true });
	};

	const onSubmit = async (data: CreateSubscriptionForm) => {
		const price = Number(data.priceInput);
		const premiumFlags = resolvePremiumSubscriptionFlags({
			autoRenewal: data.auto_renewal,
			notification: data.notification,
			canUseNotification,
		});

		await basicTrackedSubscriptionCreate({
			name: data.name.trim(),
			tariff: data.tariff,
			price,
			currency: data.currency,
			period: data.period,
			date_pay: data.date_pay,
			auto_renewal: premiumFlags.auto_renewal,
			notification: premiumFlags.notification,
			include_in_analytics: data.include_in_analytics,
			categories: data.categories,
		});

		await invalidateListAndSummary();

		notify.success(t('subscription.create-success'));
		navigate(ROUTES.HOME, { replace: true });
	};

	return (
		<PageLayout>
			<div className="mx-auto flex w-full flex-col space-y-3">
				<PageHeader title={t('subscription.create-title')} backTo={ROUTES.HOME} backLabel={t('action.back')} />

				<form className="space-y-3" onSubmit={handleSubmit(onSubmit)}>
					<section className="gu-glass-card space-y-3 px-4 py-4">
						<div className="grid grid-cols-2 gap-3">
							<div className="min-w-0 w-full">
								<input type="hidden" {...register('name', { required: t('message.validation-required-field'), minLength: { value: 3, message: t('subscription.name-min-length') } })} />

								<ServiceSelect
									value={watch('name')}
									onChange={(name, category) => {
										setValue('name', name, { shouldDirty: true, shouldValidate: true });

										if (!category || selectedCategories.includes(category)) {
											return;
										}

										setValue('categories', [...selectedCategories, category], { shouldDirty: true });
									}}
									error={errors.name?.message}
								/>
							</div>

							<div className="min-w-0 w-full">
								<input type="hidden" {...register('tariff')} />

								<TariffSelect value={watch('tariff')} onChange={(tariff) => setValue('tariff', tariff, { shouldDirty: true, shouldValidate: true })} />
							</div>
						</div>

						<div className="grid grid-cols-2 gap-3">
							<GUInput
								id="subscription-price"
								type="number"
								min="0"
								step="0.01"
								inputMode="decimal"
								placeholder={t('subscription.price-placeholder')}
								{...register('priceInput', {
									required: t('message.validation-required-field'),
									validate: (value) => Number(value) > 0 || t('subscription.price-invalid'),
								})}
								error={errors.priceInput?.message}
							/>

							<div className="w-full">
								<input type="hidden" {...register('currency', { required: t('message.validation-required-field') })} />

								<CurrencySelect value={watch('currency')} onChange={(event) => setValue('currency', event.target.value, { shouldDirty: true, shouldValidate: true })} ariaLabel={t('subscription.currency-label')} />
							</div>
						</div>

						<div className="w-full">
							<input type="hidden" {...register('period')} />

							<GUISelect
								value={watch('period')}
								onChange={(event) => setValue('period', event.target.value as TrackedSubscriptionPeriod, { shouldDirty: true, shouldValidate: true })}
								modalTitle={t('subscription.period-label')}
								aria-label={t('subscription.period-label')}
							>
								<option value="monthly">{t('home.period-monthly')}</option>
								<option value="yearly">{t('home.period-yearly')}</option>
							</GUISelect>
						</div>

						<GUInput
							id="subscription-date-pay"
							type="date"
							{...register('date_pay', {
								required: t('message.validation-required-field'),
								validate: (value) => validateFutureDate(value, t),
							})}
							error={errors.date_pay?.message}
						/>
					</section>

					<section className="gu-glass-card space-y-3 px-4 py-4">
						<h2 className="text-[15px] font-semibold gu-text-primary">{t('subscription.category')}</h2>

						<CategoryChipGroup categories={categories} selectedSlugs={selectedCategories} onToggle={toggleCategory} />
					</section>

					<SubscriptionSettingsPanel
						autoRenewal={autoRenewal}
						notification={notificationEnabled}
						includeInAnalytics={includeInAnalytics}
						onAutoRenewalChange={(checked) => setValue('auto_renewal', checked, { shouldDirty: true })}
						onNotificationChange={(checked) => setValue('notification', checked, { shouldDirty: true })}
						onIncludeInAnalyticsChange={(checked) => setValue('include_in_analytics', checked, { shouldDirty: true })}
						canUseNotification={canUseNotification}
						onPremiumRequired={() => notifyPremiumRequired(t)}
					/>

					<GUIButton type="submit" variant="primary" isLoading={isSubmitting} loadingText={t('subscription.create-loading')}>
						{t('subscription.create-action')}
					</GUIButton>
				</form>
			</div>
		</PageLayout>
	);
};

export default CreateSubscriptionPage;
