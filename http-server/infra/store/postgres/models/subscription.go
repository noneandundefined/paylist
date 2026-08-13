package models

import "time"

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
	Price              float64   `json:"price" db:"price"`
	Currency           string    `json:"currency" db:"currency"`
	Period             string    `json:"period" db:"period"`
	DatePay            time.Time `json:"date_pay" db:"date_pay"`
	AutoRenewal        bool      `json:"auto_renewal" db:"auto_renewal"`
	Notification       bool      `json:"notification" db:"notification"`
	IncludeInAnalytics bool      `json:"include_in_analytics" db:"include_in_analytics"`
	Note               *string   `json:"note,omitempty" db:"note"`
}

type TrackedSubscriptionNotifyCandidate struct {
	TrackedSubscription
	TelegramChatID   int64  `db:"telegram_chat_id"`
	TelegramLanguage string `db:"telegram_language"`
	NotifyKind       string `db:"notify_kind"`
}

type SubscriptionCategory struct {
	ID        uint64    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Slug      string    `json:"slug" db:"slug"`
	UserUUID  *string   `json:"user_uuid,omitempty" db:"user_uuid"`
	Label     *string   `json:"label,omitempty" db:"label"`
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
