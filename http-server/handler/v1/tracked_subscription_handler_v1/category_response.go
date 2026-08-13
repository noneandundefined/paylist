package tracked_subscription_handler_v1

import (
	"paylist.server/infra/store/postgres/models"
)

type SubscriptionCategoryResponse struct {
	ID        uint64  `json:"id"`
	CreatedAt string  `json:"created_at"`
	Slug      string  `json:"slug"`
	Label     *string `json:"label,omitempty"`
	IsCustom  bool    `json:"is_custom"`
}

func mapSubscriptionCategoryResponse(category models.SubscriptionCategory) SubscriptionCategoryResponse {
	return SubscriptionCategoryResponse{
		ID:        category.ID,
		CreatedAt: category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Slug:      category.Slug,
		Label:     category.Label,
		IsCustom:  category.UserUUID != nil,
	}
}

func mapSubscriptionCategoriesResponse(categories []models.SubscriptionCategory) []SubscriptionCategoryResponse {
	response := make([]SubscriptionCategoryResponse, 0, len(categories))

	for _, category := range categories {
		response = append(response, mapSubscriptionCategoryResponse(category))
	}

	return response
}
