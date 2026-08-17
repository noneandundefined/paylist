package models

import (
	"time"

	"github.com/lib/pq"
)

type Subscription struct {
	ID                        uint64              `json:"id" db:"id"`
	CreatedAt                 time.Time           `json:"created_at" db:"created_at"`
	PlanName                  string              `json:"plan_name" db:"plan_name"`
	Amount                    float64             `json:"amount" db:"amount"`
	Currency                  string              `json:"currency" db:"currency"`
	DurationDays              int                 `json:"duration_days" db:"duration_days"`
	MaxTotalSubscriptions     *int                `json:"max_total_subscriptions,omitempty" db:"max_total_subscriptions"`
	NotificationSubscriptions bool                `json:"notification_subscriptions" db:"notification_subscriptions"`
	AutoFindSubscriptions     bool                `json:"auto_find_subscriptions" db:"auto_find_subscriptions"`
	Description               map[string]string   `json:"description" db:"description"`
	Features                  map[string][]string `json:"features" db:"features"`
}

type TrackedSubscription struct {
	ID                 uint64    `json:"id" db:"id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	UserUUID           string    `json:"user_uuid" db:"user_uuid"`
	Name               string    `json:"name" db:"name"`
	Tariff             string    `json:"tariff" db:"tariff"`
	Price              float64   `json:"price" db:"price"`
	Currency           string    `json:"currency" db:"currency"`
	Period             string    `json:"period" db:"period"`
	DatePay            time.Time `json:"date_pay" db:"date_pay"`
	AutoRenewal        bool      `json:"auto_renewal" db:"auto_renewal"`
	Notification       bool      `json:"notification" db:"notification"`
	IncludeInAnalytics bool      `json:"include_in_analytics" db:"include_in_analytics"`
	Note               *string   `json:"note,omitempty" db:"note"`
	SharePercent       float64   `json:"share_percent" db:"share_percent"`
	SharePrice         float64   `json:"share_price" db:"share_price"`
	IsOwner            bool      `json:"is_owner" db:"is_owner"`
}

type TrackedSubscriptionMember struct {
	ID                    uint64     `json:"id" db:"id"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
	TrackedSubscriptionID uint64     `json:"tracked_subscription_id" db:"tracked_subscription_id"`
	UserUUID              *string    `json:"user_uuid,omitempty" db:"user_uuid"`
	Email                 string     `json:"email" db:"email"`
	Role                  string     `json:"role" db:"role"`
	SharePercent          float64    `json:"share_percent" db:"share_percent"`
	Notification          bool       `json:"notification" db:"notification"`
	IncludeInAnalytics    bool       `json:"include_in_analytics" db:"include_in_analytics"`
	Status                string     `json:"status" db:"status"`
	InviteToken           *string    `json:"-" db:"invite_token"`
	InviteExpiresAt       *time.Time `json:"invite_expires_at,omitempty" db:"invite_expires_at"`
	FirstName             *string    `json:"first_name,omitempty" db:"first_name"`
	LastName              *string    `json:"last_name,omitempty" db:"last_name"`
	Avatars               *string    `json:"avatars,omitempty" db:"avatars"`
}

type TrackedSubscriptionShareProposal struct {
	ID                    uint64    `json:"id" db:"id"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	TrackedSubscriptionID uint64    `json:"tracked_subscription_id" db:"tracked_subscription_id"`
	ProposedByUserUUID    string    `json:"proposed_by_user_uuid" db:"proposed_by_user_uuid"`
	Status                string    `json:"status" db:"status"`
}

type TrackedSubscriptionShareProposalItem struct {
	ProposalID   uint64  `json:"proposal_id" db:"proposal_id"`
	MemberID     uint64  `json:"member_id" db:"member_id"`
	SharePercent float64 `json:"share_percent" db:"share_percent"`
}

type TrackedSubscriptionShareVote struct {
	ProposalID uint64    `json:"proposal_id" db:"proposal_id"`
	UserUUID   string    `json:"user_uuid" db:"user_uuid"`
	Accepted   bool      `json:"accepted" db:"accepted"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type TrackedSubscriptionInvitePreview struct {
	SubscriptionID   uint64     `json:"subscription_id" db:"subscription_id"`
	SubscriptionName string     `json:"subscription_name" db:"subscription_name"`
	OwnerName        string     `json:"owner_name" db:"owner_name"`
	Email            string     `json:"email" db:"email"`
	SharePercent     float64    `json:"share_percent" db:"share_percent"`
	Status           string     `json:"status" db:"status"`
	Role             string     `json:"role" db:"role"`
	InviteExpiresAt  *time.Time `json:"invite_expires_at,omitempty" db:"invite_expires_at"`
}

type TrackedSubscriptionNotifyCandidate struct {
	TrackedSubscription
	MemberUserUUID   string `db:"member_user_uuid"`
	TelegramChatID   int64  `db:"telegram_chat_id"`
	TelegramLanguage string `db:"telegram_language"`
	MaxUserID        int64  `db:"max_user_id"`
	MaxLanguage      string `db:"max_language"`
	NotifyKind       string `db:"notify_kind"`
}

type SubscriptionCategory struct {
	ID        uint64    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Slug      string    `json:"slug" db:"slug"`
	UserUUID  *string   `json:"user_uuid,omitempty" db:"user_uuid"`
	Label     *string   `json:"label,omitempty" db:"label"`
}

type Service struct {
	ID        uint64         `json:"id" db:"id"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	Slug      string         `json:"slug" db:"slug"`
	Name      string         `json:"name" db:"name"`
	Category  string         `json:"category" db:"category"`
	Aliases   pq.StringArray `json:"aliases" db:"aliases"`
}

type CrowdSubscriptionPrice struct {
	Name     string  `db:"name"`
	Tariff   string  `db:"tariff"`
	Price    float64 `db:"price"`
	Currency string  `db:"currency"`
	Period   string  `db:"period"`
	Country  *string `db:"country"`
}

type TrackedSubscriptionHistory struct {
	ID                    uint64     `json:"id" db:"id"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	TrackedSubscriptionID uint64     `json:"tracked_subscription_id" db:"tracked_subscription_id"`
	UserUUID              string     `json:"user_uuid" db:"user_uuid"`
	EventType             string     `json:"event_type" db:"event_type"`
	PreviousDatePay       *time.Time `json:"previous_date_pay,omitempty" db:"previous_date_pay"`
	NewDatePay            time.Time  `json:"new_date_pay" db:"new_date_pay"`
	Price                 float64    `json:"price" db:"price"`
	Currency              string     `json:"currency" db:"currency"`
}
