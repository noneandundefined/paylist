package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"paylist.server/infra/logger"
	"paylist.server/pkg/currency"
)

const (
	cronPaymentsClearPending          = "payments.clear-pending"
	cronSubscriptionsAutoRenew        = "subscriptions.auto-renew"
	cronSubscriptionsTelegramNotify   = "subscriptions.telegram-notify"
	cronSubscriptionsMaxNotify        = "subscriptions.max-notify"
	cronSubscriptionsNotifyLogCleanup = "subscriptions.notify-log-cleanup"
	cronUsersResetExpiredPlans        = "users.reset-expired-plans"
	cronPremiumPlanPrice              = "plans.premium-price"
)

const (
	premiumPlanName = "Premium"
	premiumPriceUSD = 2.0
	premiumCurrency = "RUB"
)

func (s *httpServer) refreshPremiumPlanPrice(ctx context.Context) error {
	currency.InvalidateCBRCache()

	amount, err := currency.Convert(ctx, "USD", premiumCurrency, premiumPriceUSD)
	if err != nil {
		return err
	}

	amount = math.Round(amount*100) / 100
	if amount <= 0 {
		return fmt.Errorf("invalid CBR premium amount: %.2f", amount)
	}

	if err := s.store.Subscriptions.Update_PlanAmount(ctx, premiumPlanName, amount, premiumCurrency); err != nil {
		return err
	}

	logger.Info("[%s] Premium amount set to %.2f %s (from %.2f USD)", cronPremiumPlanPrice, amount, premiumCurrency, premiumPriceUSD)
	return nil
}

func cronContext(job string) context.Context {
	return context.WithValue(context.Background(), "XREQID", fmt.Sprintf("cron:%s", job))
}

func (s *httpServer) runCronJob(job string, fn func(context.Context) error) {
	ctx := cronContext(job)
	started := time.Now()

	logger.Info("[%s] started", job)

	if err := fn(ctx); err != nil {
		logger.Error("[%s] failed after %s: %s", job, time.Since(started), err.Error())
		return
	}

	logger.Info("[%s] finished in %s", job, time.Since(started))
}

func (s *httpServer) startCronJobs() {
	// Stale YooKassa pending payments (>10 min) — every 5 minutes.
	s.cron.AddFunc("@every 5m", func() {
		s.runCronJob(cronPaymentsClearPending, func(ctx context.Context) error {
			return s.store.Payments.Delete_PaymentWithStatusPending(ctx)
		})
	})

	// Billing day stays as-is; roll to next period at 00:00 the following day.
	// Aug 15 → Sep 15, job runs at 00:00 on Aug 16.
	s.cron.AddFunc("0 0 0 * * *", func() {
		s.runCronJob(cronSubscriptionsAutoRenew, func(ctx context.Context) error {
			return s.store.TrackedSubscriptions.Update_SubscriptionsMounth(ctx)
		})
	})

	// Messenger reminders 3 days before and on billing day — daily 09:00 UTC.
	if s.telegram != nil {
		s.cron.AddFunc("0 0 9 * * *", func() {
			s.runCronJob(cronSubscriptionsTelegramNotify, func(ctx context.Context) error {
				return s.telegram.SendDueReminders(ctx)
			})
		})
	}

	if s.maxbot != nil {
		s.cron.AddFunc("0 0 9 * * *", func() {
			s.runCronJob(cronSubscriptionsMaxNotify, func(ctx context.Context) error {
				return s.maxbot.SendDueReminders(ctx)
			})
		})
	}

	// Premium RUB price from $2 via CBR daily rate — 12:00 MSK (09:00 UTC).
	s.cron.AddFunc("0 0 9 * * *", func() {
		s.runCronJob(cronPremiumPlanPrice, s.refreshPremiumPlanPrice)
	})

	// Downgrade expired Premium SaaS plans to Free — daily 23:55 UTC.
	s.cron.AddFunc("0 55 23 * * *", func() {
		s.runCronJob(cronUsersResetExpiredPlans, func(ctx context.Context) error {
			return s.store.Users.Update_UserSubscriptionResetExpired(ctx)
		})
	})

	// Notification dedup log retention — weekly Sunday 03:30 UTC, keep 90 days.
	s.cron.AddFunc("0 30 3 * * 0", func() {
		s.runCronJob(cronSubscriptionsNotifyLogCleanup, func(ctx context.Context) error {
			return s.store.TrackedSubscriptions.Delete_OldSubscriptionNotificationLogs(ctx, 90)
		})
	})

	logger.Info("Cron scheduler registered: payments.clear-pending@5m, subscriptions.auto-renew@00:00, subscriptions.telegram-notify@09:00, subscriptions.max-notify@09:00, plans.premium-price@09:00, users.reset-expired-plans@23:55, subscriptions.notify-log-cleanup@Sun03:30")
}
