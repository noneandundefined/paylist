package models

import (
	"time"
)

const (
	PaymentStatusPending           = "pending"
	PaymentStatusWaitingForCapture = "waiting_for_capture"
	PaymentStatusSucceeded         = "succeeded"
	PaymentStatusCanceled          = "canceled"
	PaymentStatusFailed            = "failed"

	PaymentKindInitial = "initial"
	PaymentKindRenewal = "renewal"
	PaymentKindManual  = "manual"
)

type PaymentHistory struct {
	ID                      uint64     `json:"id" db:"id"`
	CreatedAt               time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at" db:"updated_at"`
	UserUUID                string     `json:"user_uuid" db:"user_uuid"`
	PlanName                string     `json:"plan_name" db:"plan_name"`
	YookassaPaymentID       string     `json:"yookassa_payment_id" db:"yookassa_payment_id"`
	YookassaPaymentMethodID *string    `json:"yookassa_payment_method_id,omitempty" db:"yookassa_payment_method_id"`
	Amount                  float64    `json:"amount" db:"amount"`
	Currency                string     `json:"currency" db:"currency"`
	Status                  string     `json:"status" db:"status"`
	PaymentKind             string     `json:"payment_kind" db:"payment_kind"`
	Description             *string    `json:"description,omitempty" db:"description"`
	PaidAt                  *time.Time `json:"paid_at,omitempty" db:"paid_at"`
	Metadata                []byte     `json:"metadata,omitempty" db:"metadata"`
}

type UserSubscriptionRenewalDue struct {
	UserUUID                string  `db:"user_uuid"`
	PlanName                string  `db:"plan_name"`
	YookassaPaymentMethodID string  `db:"yookassa_payment_method_id"`
	Amount                  float64 `db:"amount"`
	Currency                string  `db:"currency"`
	DurationDays            int     `db:"duration_days"`
}

type UserSubscriptionBilling struct {
	AutoRenewEnabled        bool       `json:"auto_renew_enabled" db:"auto_renew_enabled"`
	YookassaPaymentMethodID *string    `json:"yookassa_payment_method_id,omitempty" db:"yookassa_payment_method_id"`
	PaymentMethodType       *string    `json:"payment_method_type,omitempty" db:"payment_method_type"`
	PaymentMethodTitle      *string    `json:"payment_method_title,omitempty" db:"payment_method_title"`
	PaymentMethodSavedAt    *time.Time `json:"payment_method_saved_at,omitempty" db:"payment_method_saved_at"`
}

type UserBillingState struct {
	AutoRenewEnabled     bool       `json:"auto_renew_enabled"`
	HasPaymentMethod     bool       `json:"has_payment_method"`
	PaymentMethodType    *string    `json:"payment_method_type,omitempty"`
	PaymentMethodTitle   *string    `json:"payment_method_title,omitempty"`
	PaymentMethodSavedAt *time.Time `json:"payment_method_saved_at,omitempty"`
}

func BillingStateFromSubscription(b *UserSubscriptionBilling) UserBillingState {
	if b == nil {
		return UserBillingState{}
	}

	hasMethod := b.YookassaPaymentMethodID != nil && *b.YookassaPaymentMethodID != ""

	return UserBillingState{
		AutoRenewEnabled:     b.AutoRenewEnabled,
		HasPaymentMethod:     hasMethod,
		PaymentMethodType:    b.PaymentMethodType,
		PaymentMethodTitle:   b.PaymentMethodTitle,
		PaymentMethodSavedAt: b.PaymentMethodSavedAt,
	}
}
