package tracked_subscription_handler_v1

import (
	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/date"
)

type TrackedSubscriptionCreatePayload struct {
	Name               string    `json:"name" validate:"required,min=3"`
	Price              float64   `json:"price" validate:"required"`
	Currency           string    `json:"currency" validate:"omitempty,len=3"`
	Period             string    `json:"period" validate:"omitempty,oneof=monthly yearly"`
	DatePay            date.Date `json:"date_pay" validate:"required"`
	AutoRenewal        bool      `json:"auto_renewal"`
	Notification       bool      `json:"notification"`
	IncludeInAnalytics *bool     `json:"include_in_analytics"`
	Categories         []string  `json:"categories" validate:"omitempty,dive,min=1,max=64"`
}

type TrackedSubscriptionEditPayload struct {
	Name               string    `json:"name" validate:"required,min=3"`
	Price              float64   `json:"price" validate:"required"`
	Currency           string    `json:"currency" validate:"omitempty,len=3"`
	Period             string    `json:"period" validate:"omitempty,oneof=monthly yearly"`
	DatePay            date.Date `json:"date_pay" validate:"required"`
	AutoRenewal        bool      `json:"auto_renewal"`
	Notification       bool      `json:"notification"`
	IncludeInAnalytics *bool     `json:"include_in_analytics"`
	Categories         []string  `json:"categories" validate:"omitempty,dive,min=1,max=64"`
	Note               *string   `json:"note" validate:"omitempty,max=2000"`
}

type SubscriptionSummary struct {
	DisplayCurrency        string   `json:"display_currency"`
	TotalAmount            float64  `json:"total_amount"`
	ActiveCount            int      `json:"active_count"`
	PreviewSubscriptionIDs []uint64 `json:"preview_subscription_ids"`
	ComparisonPercent      float64  `json:"comparison_percent"`
	ComparisonDirection    string   `json:"comparison_direction"`
}

type TrackedSubscriptionDetailResponse struct {
	models.TrackedSubscription
	Categories []string `json:"categories"`
}
