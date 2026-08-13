package telegram

import (
	"context"
	"fmt"
	"time"

	"paylist.server/infra/locale"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/store"
	"paylist.server/infra/store/redis"
)

const notifyChannel = "telegram"

type Notifier struct {
	store  store.Storage
	client *Client
}

func NewNotifier(storage store.Storage, client *Client) *Notifier {
	return &Notifier{
		store:  storage,
		client: client,
	}
}

func (n *Notifier) SendDueReminders(ctx context.Context) error {
	if n.client == nil {
		return nil
	}

	candidates, err := n.store.TrackedSubscriptions.Get_SubscriptionsForTelegramNotify(ctx)
	if err != nil {
		return err
	}

	if candidates == nil || len(*candidates) == 0 {
		return nil
	}

	today := time.Now().UTC()

	for _, item := range *candidates {
		alreadySent, err := n.store.TrackedSubscriptions.Has_SubscriptionNotificationLog(ctx, item.ID, notifyChannel, today)
		if err != nil {
			logger.Error("SendDueReminders: failed to check notification log for subscription=%d: %s", item.ID, err.Error())
			continue
		}

		if alreadySent {
			continue
		}

		lang := item.TelegramLanguage
		if lang == "" {
			lang = "en"
		}

		tr := locale.NewTranslator(lang)

		localeKey := "telegram.subscription-reminder"
		if item.NotifyKind == "today" {
			localeKey = "telegram.subscription-reminder-today"
		}

		text := tr.T(localeKey)
		text = fmt.Sprintf(text, item.Name, formatPrice(item.Price, item.Currency), item.DatePay.Format("02.01.2006"))

		if err := n.client.SendMessage(item.TelegramChatID, text); err != nil {
			logger.Error("SendDueReminders: failed to send telegram message subscription=%d chat=%d: %s", item.ID, item.TelegramChatID, err.Error())
			continue
		}

		if err := n.store.TrackedSubscriptions.Create_SubscriptionNotificationLog(ctx, item.ID, notifyChannel, today); err != nil {
			logger.Error("SendDueReminders: failed to log notification subscription=%d: %s", item.ID, err.Error())
		}
	}

	return nil
}

func (n *Notifier) HandleUpdate(ctx context.Context, update Update) error {
	if update.Message == nil || update.Message.Text == "" {
		return nil
	}

	token := ParseStartToken(update.Message.Text)
	if token == "" {
		return nil
	}

	userUuid, err := redis.RedisTelegramLinkConsume(token)
	if err != nil {
		return err
	}

	if userUuid == "" {
		return n.client.SendMessage(update.Message.Chat.ID, locale.NewTranslator("en").T("telegram.link-expired"))
	}

	existingUserUuid, err := n.store.Users.Get_UserUuidByTelegramChatID(ctx, update.Message.Chat.ID)
	if err != nil {
		return err
	}

	if existingUserUuid != "" && existingUserUuid != userUuid {
		return n.client.SendMessage(update.Message.Chat.ID, locale.NewTranslator("en").T("telegram.already-linked-other-account"))
	}

	username := ""
	language := "en"

	if update.Message.From != nil {
		username = update.Message.From.Username
		if len(update.Message.From.LanguageCode) >= 2 {
			language = update.Message.From.LanguageCode[:2]
		}
	}

	if err := n.store.Users.Upsert_UserTelegram(ctx, userUuid, update.Message.Chat.ID, username, language); err != nil {
		return err
	}

	return n.client.SendMessage(update.Message.Chat.ID, locale.NewTranslator(language).T("telegram.link-success"))
}

func formatPrice(price float64, currency string) string {
	return fmt.Sprintf("%.2f %s", price, currency)
}
