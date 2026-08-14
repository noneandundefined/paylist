package store

import (
	"context"
	"database/sql"
	"time"

	"paylist.server/infra/store/postgres/models"
)

type Storage struct {
	Users interface { //nolint
		Create_UserCore(ctx context.Context, tx *sql.Tx, user *models.UserCore) error
		Create_UserSubscription(ctx context.Context, tx *sql.Tx, userUuid string) error
		Create_UserTransaction(ctx context.Context, tx *sql.Tx, user *models.UserTransaction) error

		Get_UserCoreByEmail(ctx context.Context, email string) (*models.UserCore, error)
		Get_UserCoreByUserUuid(ctx context.Context, userUuid string) (*models.UserCore, error)
		Get_UserSubscriptionsByUserUuid(ctx context.Context, userUuid string) (*models.UserSubscription, error)
		Get_UserLoginStateByUserUuid(ctx context.Context, userUuid string) (*models.UserLoginState, error)
		Get_UserPermissionsByUserUuid(ctx context.Context, userUuid string) (*models.UserPlanPermissions, error)

		Update_UserSubscriptionResetExpired(ctx context.Context) error
		Update_UserEmailConfirmedByUid(ctx context.Context, userUuid string, confirmed bool) error
		Update_UserProfile(ctx context.Context, userUuid string, firstName, lastName *string) error
		Update_UserAvatar(ctx context.Context, userUuid string, avatarURL string) error
		Update_UserEmail(ctx context.Context, userUuid, email string) error

		Delete_UserByUuid(ctx context.Context, userUuid string) error

		Get_UserSettingsByUserUuid(ctx context.Context, userUuid string) (*models.UserSettings, error)
		Upsert_UserDisplayCurrency(ctx context.Context, userUuid, currency string) error
		Upsert_UserCountry(ctx context.Context, userUuid, country string) error
		Upsert_UserTelegram(ctx context.Context, userUuid string, chatID int64, username, language string) error
		Clear_UserTelegram(ctx context.Context, userUuid string) error
		Get_UserUuidByTelegramChatID(ctx context.Context, chatID int64) (string, error)
	}
	TrackedSubscriptions interface { //nolint
		Create_Subscription(ctx context.Context, sub *models.TrackedSubscription) error
		Create_SubscriptionHistory(ctx context.Context, entry *models.TrackedSubscriptionHistory) error
		Count_SubscriptionsByUuid(ctx context.Context, uuid string) (int, error)

		Get_SubscriptionsByUuid(ctx context.Context, uuid string, search string) (*[]models.TrackedSubscription, error)
		Get_AllSubscriptionsByUuid(ctx context.Context, uuid string) ([]models.TrackedSubscription, error)
		Get_SubscriptionById(ctx context.Context, id uint64, uuid string) (*models.TrackedSubscription, error)
		Get_CategorySlugsBySubscriptionID(ctx context.Context, id uint64) ([]string, error)
		Get_CategorySlugsMapByUserUUID(ctx context.Context, userUUID string) (map[uint64][]string, error)
		Get_AllSubscriptionCategories(ctx context.Context) ([]models.SubscriptionCategory, error)
		Get_SubscriptionCategoriesForUser(ctx context.Context, userUuid string) ([]models.SubscriptionCategory, error)
		Create_UserSubscriptionCategory(ctx context.Context, userUuid, slug, label string) (*models.SubscriptionCategory, error)
		Delete_UserSubscriptionCategory(ctx context.Context, userUuid string, categoryID uint64) error
		Get_SubscriptionsForTelegramNotify(ctx context.Context) (*[]models.TrackedSubscriptionNotifyCandidate, error)
		Create_SubscriptionNotificationLog(ctx context.Context, subscriptionID uint64, channel string, notifyDate time.Time) error
		Has_SubscriptionNotificationLog(ctx context.Context, subscriptionID uint64, channel string, notifyDate time.Time) (bool, error)
		Delete_OldSubscriptionNotificationLogs(ctx context.Context, olderThanDays int) error

		Update_SubscriptionsMounth(ctx context.Context) error
		Update_SubscriptionById(ctx context.Context, sub *models.TrackedSubscription, id int) error
		Replace_SubscriptionCategories(ctx context.Context, id uint64, userUUID string, slugs []string) error
		Delete_SubscriptionById(ctx context.Context, id int, uuid string) error
	}
	Payments interface { //nolint
		Create_PaymentHistory(ctx context.Context, payment *models.PaymentHistory) (*models.PaymentHistory, error)

		Get_PaymentHistoryListByUserUuid(ctx context.Context, userUuid string, limit int) ([]models.PaymentHistory, error)
		Get_PaymentHistoryByYookassaPaymentID(ctx context.Context, paymentID string) (*models.PaymentHistory, error)
		Get_UserSubscriptionBillingByUserUuid(ctx context.Context, userUuid string) (*models.UserSubscriptionBilling, error)
		Get_PaymentActiveCount(ctx context.Context, userUuid string) (uint32, error)

		Update_PaymentHistoryStatus(ctx context.Context, yookassaPaymentID, status string, paidAt *time.Time) error
		Update_UserSubscriptionAutoRenew(ctx context.Context, userUuid string, enabled bool) error
		Update_YookassaPaymentMethod(ctx context.Context, userUuid, paymentMethodID, paymentMethodType, paymentMethodTitle string) error
		Update_ClearYookassaPaymentMethod(ctx context.Context, userUuid string) error
		Update_ActivateUserSubscriptionPlan(ctx context.Context, userUuid, planName string, durationDays int) error

		Delete_PaymentWithStatusPending(ctx context.Context) error
	}
	Subscriptions interface { //nolint
		Get_Subscriptions(ctx context.Context) ([]models.Subscription, error)
		Get_SubscriptionByPlanName(ctx context.Context, planName string) (*models.Subscription, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users:                &UserStore{db},
		TrackedSubscriptions: &TrackedSubscriptionStore{db},
		Payments:             &PaymentStore{db},
		Subscriptions:        &SubscriptionStore{db},
	}
}

func WithTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
