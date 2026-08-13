import { useTranslation } from 'react-i18next';
import type { SubscriptionCategoryResponse } from '@/rest/trackedSubscriptionAPI';
import { getCategoryLabel } from '@/utils/categoryDisplayUtils';

import Close from '@/components/@icons/close';

interface CategoryChipGroupProps {
	categories: SubscriptionCategoryResponse[];
	selectedSlugs?: string[];
	onToggle?: (slug: string) => void;
	onDelete?: (categoryId: number) => void;
	deleteAriaLabel?: string;
}

const CategoryChipGroup: React.FC<CategoryChipGroupProps> = ({ categories, selectedSlugs = [], onToggle, onDelete, deleteAriaLabel }) => {
	const { t } = useTranslation();

	return (
		<div className="flex flex-wrap gap-2">
			{categories.map((category) => {
				if (onDelete) {
					return (
						<span key={category.id} className="gu-chip gu-chip--active inline-flex items-center gap-1.5 pl-3 pr-1.5">
							{getCategoryLabel(category, t)}

							<button
								type="button"
								onClick={() => onDelete(category.id)}
								className="inline-flex h-5 w-5 shrink-0 items-center justify-center rounded-full transition hover:bg-black/10"
								aria-label={deleteAriaLabel ?? t('account.category-delete')}
							>
								<Close fill="currentColor" size={14} />
							</button>
						</span>
					);
				}

				const isSelected = selectedSlugs.includes(category.slug);

				return (
					<button key={category.slug} type="button" onClick={() => onToggle?.(category.slug)} className={`gu-chip ${isSelected ? 'gu-chip--active' : ''}`}>
						{getCategoryLabel(category, t)}
					</button>
				);
			})}
		</div>
	);
};

export default CategoryChipGroup;
