import { useState } from 'react';
import Plus from '@/components/@icons/plus';
import { useTranslation } from 'react-i18next';
import { useConfirm } from '@/hooks/useConfirm';
import { useQueryClient } from '@tanstack/react-query';
import { GUInput } from '@/components/ui/Input/GUInput';
import GUIButton from '@/components/ui/Button/GUIButton';
import { notify } from '@/components/Notification/notify';
import { notifyPremiumRequired } from '@/utils/premiumUtils';
import { QUERY_KEYS } from '@/constants/QueryKeys.constants';
import { isCustomCategory } from '@/utils/categoryDisplayUtils';
import { useSubscriptionCategories } from '@/hooks/useSubscriptionCategories';
import CategoryChipGroup from '@/components/common/Category/CategoryChipGroup';
import PremiumBadgeMini from '@/components/common/PremiumBadge/PremiumBadgeMini';
import { basicSubscriptionCategoryCreate, basicSubscriptionCategoryDelete } from '@/rest/trackedSubscriptionAPI';

interface AccountCategoriesManagerProps {
	isPremium: boolean;
}

const AccountCategoriesManager: React.FC<AccountCategoriesManagerProps> = ({ isPremium }) => {
	const { t } = useTranslation();

	const { confirm } = useConfirm();

	const queryClient = useQueryClient();

	const { categories: categoriesData } = useSubscriptionCategories();

	const [label, setLabel] = useState('');
	const [saving, setSaving] = useState(false);

	const customCategories = categoriesData.filter(isCustomCategory);

	const onAddCategory = async () => {
		if (!isPremium) {
			notifyPremiumRequired(t);

			return;
		}

		const trimmed = label.trim();

		if (trimmed.length < 2) {
			return;
		}

		setSaving(true);

		try {
			await basicSubscriptionCategoryCreate(trimmed);

			setLabel('');

			notify.success(t('account.category-added'));

			await queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.trackedSubscriptionCategoryList] });
		} finally {
			setSaving(false);
		}
	};

	const onDeleteCategory = async (categoryId: number) => {
		if (!(await confirm('account.category-delete-confirm', 'account.categories'))) {
			return;
		}

		await basicSubscriptionCategoryDelete(categoryId);

		await queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.trackedSubscriptionCategoryList] });
	};

	return (
		<div className="p-3 space-y-3">
			<div className="flex items-center justify-between gap-3">
				<h2 className="text-[15px] font-semibold gu-text-primary">{t('account.categories')}</h2>

				{!isPremium && <PremiumBadgeMini mobileView={true} />}
			</div>

			{isPremium && (
				<>
					<div className="flex gap-2">
						<GUInput
							value={label}
							onChange={(event) => setLabel(event.target.value)}
							placeholder={t('account.categories-add-placeholder')}
							aria-label={t('account.categories-add-placeholder')}
							className="min-w-0 flex-1"
						/>

						<GUIButton type="button" onClick={onAddCategory} disabled={saving || label.trim().length < 2} className="gu-glass-icon-btn shrink-0 px-4">
							<span className="gu-glass-icon-btn">
								<Plus fill="#000" size={18} />
							</span>
						</GUIButton>
					</div>

					{customCategories.length === 0 ? (
						<p className="text-[13px] gu-text-muted">{t('account.categories-empty')}</p>
					) : (
						<CategoryChipGroup categories={customCategories} onDelete={onDeleteCategory} deleteAriaLabel={t('account.category-delete')} />
					)}
				</>
			)}
		</div>
	);
};

export default AccountCategoriesManager;
