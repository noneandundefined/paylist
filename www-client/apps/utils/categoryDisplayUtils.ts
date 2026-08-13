import type { TFunction } from 'i18next';
import type { SubscriptionCategoryResponse } from '@/rest/trackedSubscriptionAPI';

export const getCategoryLabel = (category: SubscriptionCategoryResponse, t: TFunction): string => {
	if (category.label) {
		return category.label;
	}

	return t(`subscription.category-${category.slug}`, category.slug);
};

export const isCustomCategory = (category: SubscriptionCategoryResponse): boolean => {
	return category.is_custom === true || Boolean(category.label);
};

export const areCategoriesEqual = (left: string[], right: string[]): boolean => {
	if (left.length !== right.length) {
		return false;
	}

	const sortedLeft = [...left].sort();
	const sortedRight = [...right].sort();

	return sortedLeft.every((item, index) => item === sortedRight[index]);
};
