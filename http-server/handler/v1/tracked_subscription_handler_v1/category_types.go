package tracked_subscription_handler_v1

type SubscriptionCategoryCreatePayload struct {
	Label string `json:"label" validate:"required,min=2,max=64"`
}
