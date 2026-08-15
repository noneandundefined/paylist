package maxbot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"paylist.server/infra/locale"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/store"
	"paylist.server/infra/store/redis"
)

const notifyChannel = "max"

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

	candidates, err := n.store.TrackedSubscriptions.Get_SubscriptionsForMaxNotify(ctx)
	if err != nil {
		return err
	}

	if candidates == nil || len(*candidates) == 0 {
		return nil
	}

	today := time.Now().UTC()

	for _, item := range *candidates {
		alreadySent, err := n.store.TrackedSubscriptions.Has_SubscriptionNotificationLog(ctx, item.ID, item.MemberUserUUID, notifyChannel, today)
		if err != nil {
			logger.Error("SendDueReminders: failed to check max notification log for subscription=%d: %s", item.ID, err.Error())
			continue
		}

		if alreadySent {
			continue
		}

		lang := item.MaxLanguage
		if lang == "" {
			lang = "en"
		}

		tr := locale.NewTranslator(lang)

		localeKey := "telegram.subscription-reminder"
		if item.NotifyKind == "today" {
			localeKey = "telegram.subscription-reminder-today"
		}

		text := tr.T(localeKey)
		text = fmt.Sprintf(text, item.Name, formatPrice(item.SharePrice, item.Currency), item.DatePay.Format("02.01.2006"))

		if err := n.client.SendMessage(item.MaxUserID, text); err != nil {
			logger.Error("SendDueReminders: failed to send max message subscription=%d user=%d: %s", item.ID, item.MaxUserID, err.Error())
			continue
		}

		if err := n.store.TrackedSubscriptions.Create_SubscriptionNotificationLog(ctx, item.ID, item.MemberUserUUID, notifyChannel, today); err != nil {
			logger.Error("SendDueReminders: failed to log max notification subscription=%d: %s", item.ID, err.Error())
		}
	}

	return nil
}

func (n *Notifier) HandleUpdate(ctx context.Context, update Update) error {
	userID, username, token := linkPayload(update)
	if userID == 0 || token == "" {
		return nil
	}

	userUuid, language, err := redis.RedisMaxLinkConsume(token)
	if err != nil {
		return err
	}

	if userUuid == "" {
		return n.client.SendMessage(userID, locale.NewTranslator("en").T("max.link-expired"))
	}

	existingUserUuid, err := n.store.Users.Get_UserUuidByMaxUserID(ctx, userID)
	if err != nil {
		return err
	}

	if existingUserUuid != "" && existingUserUuid != userUuid {
		return n.client.SendMessage(userID, locale.NewTranslator(language).T("max.already-linked-other-account"))
	}

	if language == "" {
		language = "en"
	}

	if len(language) > 2 {
		language = language[:2]
	}

	if err := n.store.Users.Upsert_UserMax(ctx, userUuid, userID, username, language); err != nil {
		return err
	}

	return n.client.SendMessage(userID, locale.NewTranslator(language).T("max.link-success"))
}

func linkPayload(update Update) (int64, string, string) {
	switch update.UpdateType {
	case "bot_started":
		if update.User == nil {
			return 0, "", ""
		}

		return update.User.UserID, update.User.Username, strings.TrimSpace(update.Payload)
	case "message_created":
		if update.Message == nil || update.Message.Sender == nil || update.Message.Body == nil {
			return 0, "", ""
		}

		return update.Message.Sender.UserID, update.Message.Sender.Username, ParseStartToken(update.Message.Body.Text)
	default:
		return 0, "", ""
	}
}

func formatPrice(price float64, currency string) string {
	return fmt.Sprintf("%.2f %s", price, currency)
}
