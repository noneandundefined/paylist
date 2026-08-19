package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"paylist.server/infra/locale"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/store"
	"paylist.server/infra/store/redis"
	"paylist.server/pkg/referral"
)

var errDeliveryRejected = errors.New("telegram delivery rejected")

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
		alreadySent, err := n.store.TrackedSubscriptions.Has_SubscriptionNotificationLog(ctx, item.ID, item.MemberUserUUID, notifyChannel, today)
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
		text = fmt.Sprintf(text, item.Name, formatPrice(item.SharePrice, item.Currency), item.DatePay.Format("02.01.2006"))

		if err := n.client.SendMessage(item.TelegramChatID, text); err != nil {
			if IsDeliveryRejected(err) {
				n.disconnectBlockedChat(ctx, item.TelegramChatID)
				continue
			}

			logger.Error("SendDueReminders: failed to send telegram message subscription=%d chat=%d: %s", item.ID, item.TelegramChatID, err.Error())
			continue
		}

		if err := n.store.TrackedSubscriptions.Create_SubscriptionNotificationLog(ctx, item.ID, item.MemberUserUUID, notifyChannel, today); err != nil {
			logger.Error("SendDueReminders: failed to log notification subscription=%d: %s", item.ID, err.Error())
		}
	}

	return nil
}

func (n *Notifier) HandleUpdate(ctx context.Context, update Update) error {
	if update.Message == nil {
		return nil
	}

	text := strings.TrimSpace(update.Message.Text)
	if !isStartCommand(text) {
		return nil
	}

	chatID := update.Message.Chat.ID
	username, language := fromTelegramUser(update.Message.From)
	tr := locale.NewTranslator(language)
	token := ParseStartToken(text)

	if code := referral.CodeFromStartToken(token); code != "" {
		return n.handleReferralStart(ctx, chatID, code, tr)
	}

	if token != "" {
		return finishUpdate(n.linkFromStartToken(ctx, chatID, token, username, language, tr))
	}

	existingUserUuid, err := n.store.Users.Get_UserUuidByTelegramChatID(ctx, chatID)
	if err != nil {
		return err
	}

	if existingUserUuid != "" {
		return finishUpdate(n.sendOrDetach(ctx, chatID, n.client.SendMessageWithOpenApp(chatID, tr.T("telegram.account-created"), tr.T("telegram.open-app"))))
	}

	if err := n.sendOrDetach(ctx, chatID, n.client.SendMessage(chatID, tr.T("telegram.welcome"))); err != nil {
		return finishUpdate(err)
	}

	return finishUpdate(n.sendOrDetach(ctx, chatID, n.client.SendMessageWithOpenApp(chatID, tr.T("telegram.sign-in-hint"), tr.T("telegram.open-app"))))
}

func (n *Notifier) handleReferralStart(ctx context.Context, chatID int64, code string, tr locale.Translator) error {
	signupURL := referral.SiteURL(publicAppURL(), code)

	existingUserUuid, err := n.store.Users.Get_UserUuidByTelegramChatID(ctx, chatID)
	if err != nil {
		return err
	}

	if existingUserUuid != "" {
		return finishUpdate(n.sendOrDetach(ctx, chatID, n.client.SendMessageWithOpenApp(chatID, tr.T("telegram.account-created"), tr.T("telegram.open-app"))))
	}

	if err := n.sendOrDetach(ctx, chatID, n.client.SendMessage(chatID, tr.T("telegram.welcome"))); err != nil {
		return finishUpdate(err)
	}

	return finishUpdate(n.sendOrDetach(ctx, chatID, n.client.SendMessageWithLink(chatID, tr.T("telegram.sign-in-hint"), tr.T("telegram.open-app"), signupURL)))
}

func (n *Notifier) linkFromStartToken(ctx context.Context, chatID int64, token, username, language string, tr locale.Translator) error {
	userUuid, err := redis.RedisTelegramLinkGet(token)
	if err != nil {
		return err
	}

	existingUserUuid, err := n.store.Users.Get_UserUuidByTelegramChatID(ctx, chatID)
	if err != nil {
		return err
	}

	if userUuid == "" {
		if existingUserUuid != "" {
			return n.sendOrDetach(ctx, chatID, n.client.SendMessageWithOpenApp(chatID, tr.T("telegram.link-success"), tr.T("telegram.open-app")))
		}

		return n.sendOrDetach(ctx, chatID, n.client.SendMessage(chatID, tr.T("telegram.link-expired")))
	}

	if existingUserUuid != "" && existingUserUuid != userUuid {
		return n.sendOrDetach(ctx, chatID, n.client.SendMessage(chatID, tr.T("telegram.already-linked-other-account")))
	}

	if err := n.store.Users.Upsert_UserTelegram(ctx, userUuid, chatID, username, language); err != nil {
		return err
	}

	_ = redis.RedisTelegramLinkDelete(token)

	return n.sendOrDetach(ctx, chatID, n.client.SendMessageWithOpenApp(chatID, tr.T("telegram.link-success"), tr.T("telegram.open-app")))
}

func (n *Notifier) sendOrDetach(ctx context.Context, chatID int64, err error) error {
	if err == nil {
		return nil
	}

	if IsDeliveryRejected(err) {
		n.disconnectBlockedChat(ctx, chatID)
		return errDeliveryRejected
	}

	return err
}

func (n *Notifier) disconnectBlockedChat(ctx context.Context, chatID int64) {
	userUuid, err := n.store.Users.Get_UserUuidByTelegramChatID(ctx, chatID)
	if err != nil {
		logger.Error("Telegram unlink failed chat=%d: %s", chatID, err.Error())
		return
	}

	if userUuid != "" {
		if err := n.store.Users.Clear_UserTelegram(ctx, userUuid); err != nil {
			logger.Error("Telegram unlink failed user=%s chat=%d: %s", userUuid, chatID, err.Error())
			return
		}
	}

	logger.Warning("Telegram bot was blocked chat=%d user=%s", chatID, userUuid)
}

func finishUpdate(err error) error {
	if errors.Is(err, errDeliveryRejected) {
		return nil
	}

	return err
}

func isStartCommand(text string) bool {
	return strings.HasPrefix(text, "/start")
}

func fromTelegramUser(from *User) (string, string) {
	username := ""
	language := "ru"

	if from == nil {
		return username, language
	}

	username = strings.TrimSpace(from.Username)
	if len(from.LanguageCode) >= 2 {
		language = strings.ToLower(from.LanguageCode[:2])
	}

	return username, language
}

func formatPrice(price float64, currency string) string {
	return fmt.Sprintf("%.2f %s", price, currency)
}
