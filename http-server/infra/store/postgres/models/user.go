package models

import "time"

type UserCore struct {
	ID             uint64    `json:"id" db:"id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	UserUUID       string    `json:"user_uuid" db:"user_uuid"`
	Email          string    `json:"email" db:"email"`
	EmailConfirmed bool      `json:"email_confirmed" db:"email_confirmed"`
	FirstName      *string   `json:"first_name,omitempty" db:"first_name"`
	LastName       *string   `json:"last_name,omitempty" db:"last_name"`
	Password       string    `json:"-" db:"password"`
}

type UserSubscription struct {
	ID                      uint64     `json:"id" db:"id"`
	CreatedAt               time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at" db:"updated_at"`
	UserUUID                string     `json:"user_uuid" db:"user_uuid"`
	PlanName                string     `json:"plan_name" db:"plan_name"`
	ValidFrom               time.Time  `json:"valid_from" db:"valid_from"`
	ValidTo                 *time.Time `json:"valid_to,omitempty" db:"valid_to"`
	IsActive                bool       `json:"is_active" db:"is_active"`
	AutoRenewEnabled        bool       `json:"auto_renew_enabled" db:"auto_renew_enabled"`
	YookassaPaymentMethodID *string    `json:"yookassa_payment_method_id,omitempty" db:"yookassa_payment_method_id"`
	PaymentMethodType       *string    `json:"payment_method_type,omitempty" db:"payment_method_type"`
	PaymentMethodTitle      *string    `json:"payment_method_title,omitempty" db:"payment_method_title"`
	PaymentMethodSavedAt    *time.Time `json:"payment_method_saved_at,omitempty" db:"payment_method_saved_at"`

	// Used when activating a plan from user_transactions; not a DB column.
	SubscriptionPlanID int64 `json:"-" db:"-"`
}

type UserTransaction struct {
	ID              uint64    `json:"id" db:"id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	UserUUID        string    `json:"user_uuid" db:"user_uuid"`
	SubscriptionID  int64     `json:"subscription_id" db:"subscription_id"`
	TransactionType string    `json:"transaction_type" db:"transaction_type"`
	Amount          float64   `json:"amount" db:"amount"`
	Currency        string    `json:"currency" db:"currency"`
	Description     *string   `json:"description,omitempty" db:"description"`
}

type UserLoginState struct {
	/* User cores */
	ID             uint64    `json:"id" db:"id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	Email          string    `json:"email" db:"email"`
	EmailConfirmed bool      `json:"email_confirmed" db:"email_confirmed"`
	FirstName      *string   `json:"first_name" db:"first_name"`
	LastName       *string   `json:"last_name" db:"last_name"`

	/* User subscriptions */
	PlanName string     `json:"plan_name" db:"plan_name"`
	ValidTo  *time.Time `json:"valid_to,omitempty" db:"valid_to"`
	Amount   float64    `json:"amount" db:"amount"`
	Currency string     `json:"currency" db:"currency"`

	NotificationSubscriptions bool `json:"notification_subscriptions" db:"notification_subscriptions"`
	MaxTotalSubscriptions     *int `json:"max_total_subscriptions,omitempty" db:"max_total_subscriptions"`
	AutoFindSubscriptions     bool `json:"auto_find_subscriptions" db:"auto_find_subscriptions"`
}

type UserPlanPermissions struct {
	PlanName                  string
	NotificationSubscriptions bool
	MaxTotalSubscriptions     *int
	AutoFindSubscriptions     bool
}

type UserSettings struct {
	UserUUID         string    `json:"user_uuid" db:"user_uuid"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	DisplayCurrency  *string   `json:"display_currency,omitempty" db:"display_currency"`
	Country          *string   `json:"country,omitempty" db:"country"`
	TelegramChatID   *int64    `json:"telegram_chat_id,omitempty" db:"telegram_chat_id"`
	TelegramUsername *string   `json:"telegram_username,omitempty" db:"telegram_username"`
	TelegramLanguage *string   `json:"telegram_language,omitempty" db:"telegram_language"`
}
