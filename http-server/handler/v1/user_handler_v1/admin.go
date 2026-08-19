package user_handler_v1

import (
	"context"
	"html"
	"net/http"
	"strings"
	"unicode/utf8"

	"paylist.server/infra/locale"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/models"
	"paylist.server/middleware"
	"paylist.server/pkg"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/maxbot"
	"paylist.server/pkg/telegram"
	"paylist.server/types"
)

const adminMessageMaxRunes = 4000

type AdminSendMessagePayload struct {
	Channel  string  `json:"channel" validate:"required,oneof=email telegram max"`
	UserUUID *string `json:"user_uuid"`
	Text     string  `json:"text" validate:"required"`
}

type AdminSendMessageResponse struct {
	Sent    int `json:"sent"`
	Skipped int `json:"skipped"`
	Failed  int `json:"failed"`
}

func requireAdmin(authToken *types.AuthToken, tr locale.Translator) error {
	if authToken == nil || !authToken.IsAdmin {
		return httperr.Forbidden(tr.TErr("error.admin-required"))
	}

	return nil
}

func (h *Handler) AdminRecipientsHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	if err := requireAdmin(authToken, tr); err != nil {
		return err
	}

	users, err := h.Store.Users.List_AdminMessageRecipients(ctx)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, users)
	return nil
}

func (h *Handler) AdminSendMessageHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	if err := requireAdmin(authToken, tr); err != nil {
		return err
	}

	var payload AdminSendMessagePayload
	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	text := normalizeAdminMessage(payload.Text)
	if text == "" {
		return httperr.BadRequest(tr.TErr("error.admin-message-empty"))
	}

	if utf8.RuneCountInString(text) > adminMessageMaxRunes {
		return httperr.BadRequest(tr.TErr("error.admin-message-too-long"))
	}

	var userUUID *string
	if payload.UserUUID != nil {
		trimmed := strings.TrimSpace(*payload.UserUUID)
		if trimmed != "" && trimmed != "all" {
			userUUID = &trimmed
		}
	}

	targets, err := h.Store.Users.List_AdminMessageTargets(ctx, userUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if userUUID != nil && len(targets) == 0 {
		return httperr.NotFound(tr.TErr("error.user-not-found"))
	}

	result := h.sendAdminMessages(ctx, payload.Channel, text, targets, tr)
	if userUUID != nil && result.Sent == 0 && result.Skipped > 0 {
		return httperr.BadRequest(tr.TErr("error.admin-channel-not-connected"))
	}

	if result.Sent == 0 && result.Failed > 0 {
		return httperr.InternalServerError(tr.TErr("error.admin-message-failed"))
	}

	httpx.HttpResponse(w, r, http.StatusOK, result)
	return nil
}

func (h *Handler) sendAdminMessages(ctx context.Context, channel, text string, targets []models.AdminMessageTarget, tr locale.Translator) AdminSendMessageResponse {
	var result AdminSendMessageResponse
	htmlBody := adminMessageHTML(text)
	subject := adminMessageSubject(text)

	var telegramClient *telegram.Client
	var maxClient *maxbot.Client

	if channel == "telegram" {
		client, err := telegram.NewFromEnv()
		if err != nil {
			logger.Error("AdminSendMessage: telegram is not configured: %s", err.Error())
			result.Failed = len(targets)
			return result
		}

		telegramClient = client
	}

	if channel == "max" {
		client, err := maxbot.NewFromEnv()
		if err != nil {
			logger.Error("AdminSendMessage: max is not configured: %s", err.Error())
			result.Failed = len(targets)
			return result
		}

		maxClient = client
	}

	for _, target := range targets {
		switch channel {
		case "email":
			if err := pkg.SendEmail(target.Email, subject, htmlBody, tr); err != nil {
				logger.Error("AdminSendMessage email={%s}: %s", target.Email, err.Error())
				result.Failed++
				continue
			}

			result.Sent++

		case "telegram":
			if target.TelegramChatID == nil {
				result.Skipped++
				continue
			}

			if err := telegramClient.SendMessage(*target.TelegramChatID, text); err != nil {
				if telegram.IsDeliveryRejected(err) {
					logger.Warning("AdminSendMessage telegram user={%s}: bot was blocked, unlinking", target.UserUUID)
					if clearErr := h.Store.Users.Clear_UserTelegram(ctx, target.UserUUID); clearErr != nil {
						logger.Error("AdminSendMessage telegram unlink user={%s}: %s", target.UserUUID, clearErr.Error())
					}

					result.Skipped++
					continue
				}

				logger.Error("AdminSendMessage telegram user={%s}: %s", target.UserUUID, err.Error())
				result.Failed++
				continue
			}

			result.Sent++

		case "max":
			if target.MaxUserID == nil {
				result.Skipped++
				continue
			}

			if err := maxClient.SendMessage(*target.MaxUserID, text); err != nil {
				logger.Error("AdminSendMessage max user={%s}: %s", target.UserUUID, err.Error())
				result.Failed++
				continue
			}

			result.Sent++
		}
	}

	return result
}

func normalizeAdminMessage(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	if strings.TrimSpace(text) == "" {
		return ""
	}

	return strings.Trim(text, " \t")
}

func adminMessageHTML(text string) string {
	escaped := html.EscapeString(text)
	return strings.ReplaceAll(escaped, "\n", "<br>")
}

func adminMessageSubject(text string) string {
	firstLine := text
	if index := strings.IndexByte(text, '\n'); index >= 0 {
		firstLine = strings.TrimSpace(text[:index])
	}

	if firstLine == "" {
		return "Paylist"
	}

	runes := []rune(firstLine)
	if len(runes) > 80 {
		return string(runes[:80])
	}

	return firstLine
}
