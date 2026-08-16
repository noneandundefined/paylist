package main

import (
	"context"
	"fmt"
	"time"

	"paylist.server/handler/v1/payment_handler_v1"
	"paylist.server/infra/logger"
	"paylist.server/pkg/yookassa"
)

const (
	cronPaymentsClearPending          = "payments.clear-pending"
	cronPaymentsAutoRenew             = "payments.auto-renew"
	cronSubscriptionsAutoRenew        = "subscriptions.auto-renew"
	cronSubscriptionsTelegramNotify   = "subscriptions.telegram-notify"
	cronSubscriptionsMaxNotify        = "subscriptions.max-notify"
	cronSubscriptionsNotifyLogCleanup = "subscriptions.notify-log-cleanup"
	cronUsersResetExpiredPlans        = "users.reset-expired-plans"
)

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

	// YooKassa Premium auto-renew — 21:00 MSK (18:00 UTC) on the billing day.
	s.cron.AddFunc("0 0 18 * * *", func() {
		s.runCronJob(cronPaymentsAutoRenew, func(ctx context.Context) error {
			client, err := yookassa.NewFromEnv()
			if err != nil {
				return nil
			}

			return payment_handler_v1.ChargeDueRenewals(ctx, s.store, client)
		})
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

	logger.Info("Cron scheduler registered: payments.clear-pending@5m, payments.auto-renew@18:00UTC/21:00MSK, subscriptions.auto-renew@00:00, subscriptions.telegram-notify@09:00, subscriptions.max-notify@09:00, users.reset-expired-plans@23:55, subscriptions.notify-log-cleanup@Sun03:30")
}
