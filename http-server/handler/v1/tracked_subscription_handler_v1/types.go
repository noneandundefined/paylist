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

type TrackedSubscriptionInvitePayload struct {
	Email        string  `json:"email" validate:"required,email"`
	SharePercent float64 `json:"share_percent" validate:"required,gt=0,lte=100"`
}

type TrackedSubscriptionInviteAcceptPayload struct {
	Token string `json:"token" validate:"required"`
}

type TrackedSubscriptionShareItemPayload struct {
	MemberID     uint64  `json:"member_id" validate:"required"`
	SharePercent float64 `json:"share_percent" validate:"required,gt=0,lte=100"`
}

type TrackedSubscriptionShareProposalPayload struct {
	Shares []TrackedSubscriptionShareItemPayload `json:"shares" validate:"required,min=1,dive"`
}

type TrackedSubscriptionShareVotePayload struct {
	Accept bool `json:"accept"`
}

type TrackedSubscriptionShareProposalResponse struct {
	ID                 uint64                                        `json:"id"`
	ProposedByUserUUID string                                        `json:"proposed_by_user_uuid"`
	Status             string                                        `json:"status"`
	Items              []models.TrackedSubscriptionShareProposalItem `json:"items"`
	Votes              []models.TrackedSubscriptionShareVote         `json:"votes"`
	MyVote             *bool                                         `json:"my_vote"`
}

type TrackedSubscriptionMembersResponse struct {
	Members         []models.TrackedSubscriptionMember        `json:"members"`
	PendingProposal *TrackedSubscriptionShareProposalResponse `json:"pending_proposal"`
}

type TrackedSubscriptionInviteAcceptResponse struct {
	Message        string `json:"message"`
	SubscriptionID uint64 `json:"subscription_id"`
}
