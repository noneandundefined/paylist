package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const longPollTimeoutSec = 25

type Client struct {
	token          string
	botName        string
	shortClient    *http.Client
	longPollClient *http.Client
}

type APIError struct {
	Code        int
	Description string
}

func (e *APIError) Error() string {
	if e.Description != "" {
		return e.Description
	}

	return fmt.Sprintf("telegram API error %d", e.Code)
}

func parseAPIError(statusCode int, responseBody []byte) error {
	var body struct {
		OK          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}

	if err := json.Unmarshal(responseBody, &body); err == nil && (body.ErrorCode != 0 || body.Description != "") {
		code := body.ErrorCode
		if code == 0 {
			code = statusCode
		}

		return &APIError{Code: code, Description: body.Description}
	}

	return fmt.Errorf("telegram request failed: %s", strings.TrimSpace(string(responseBody)))
}

func IsDeliveryRejected(err error) bool {
	if err == nil {
		return false
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		desc := strings.ToLower(apiErr.Description)
		if apiErr.Code == http.StatusForbidden {
			return true
		}

		if apiErr.Code == http.StatusBadRequest && (strings.Contains(desc, "chat not found") || strings.Contains(desc, "user not found") || strings.Contains(desc, "peer_id_invalid")) {
			return true
		}
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "bot was blocked") ||
		strings.Contains(msg, "bot was kicked") ||
		strings.Contains(msg, "user is deactivated") ||
		strings.Contains(msg, "chat not found")
}

type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	Text      string `json:"text"`
	Chat      Chat   `json:"chat"`
	From      *User  `json:"from"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

func NewFromEnv() (*Client, error) {
	token := strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is not set")
	}

	botName := strings.TrimSpace(os.Getenv("TELEGRAM_BOT_USERNAME"))
	if botName == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_USERNAME is not set")
	}

	return &Client{
		token:   token,
		botName: strings.TrimPrefix(botName, "@"),
		shortClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		longPollClient: &http.Client{
			// Must exceed long-poll timeout or getUpdates cancels early and may trigger 409 conflicts.
			Timeout: time.Duration(longPollTimeoutSec+15) * time.Second,
		},
	}, nil
}

func (c *Client) BotURL(startToken string) string {
	return fmt.Sprintf("https://t.me/%s?start=%s", c.botName, startToken)
}

func (c *Client) SendMessage(chatID int64, text string) error {
	return c.sendMessage(chatID, text, nil)
}

func (c *Client) SendMessageWithOpenApp(chatID int64, text, buttonText string) error {
	return c.SendMessageWithLink(chatID, text, buttonText, publicAppURL())
}

func (c *Client) SendMessageWithLink(chatID int64, text, buttonText, appURL string) error {
	return c.sendMessage(chatID, text, openAppMarkupURL(buttonText, appURL))
}

func (c *Client) sendMessage(chatID int64, text string, replyMarkup map[string]any) error {
	payload := map[string]any{
		"chat_id": chatID,
		"text":    text,
	}
	if replyMarkup != nil {
		payload["reply_markup"] = replyMarkup
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.token), bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.shortClient.Do(req)
	if err != nil {
		return c.redactError(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.redactError(err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		err := parseAPIError(resp.StatusCode, responseBody)
		if replyMarkup != nil && isInvalidButtonURL(err) {
			return c.sendMessage(chatID, text, nil)
		}

		return err
	}

	return nil
}

func openAppMarkup(buttonText string) map[string]any {
	return openAppMarkupURL(buttonText, publicAppURL())
}

func openAppMarkupURL(buttonText, appURL string) map[string]any {
	buttonText = strings.TrimSpace(buttonText)
	if buttonText == "" {
		buttonText = "Open Paylist"
	}

	appURL = strings.TrimSpace(appURL)
	if !strings.HasPrefix(strings.ToLower(appURL), "https://") {
		appURL = publicAppURL()
	}

	return map[string]any{
		"inline_keyboard": [][]map[string]string{
			{
				{"text": buttonText, "url": appURL},
			},
		},
	}
}

func publicAppURL() string {
	appURL := strings.TrimRight(strings.TrimSpace(os.Getenv("CLIENT_URL")), "/")
	if strings.HasPrefix(strings.ToLower(appURL), "https://") {
		return appURL
	}

	return "https://paylist.site"
}

func isInvalidButtonURL(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "button url") || strings.Contains(msg, "wrong http url")
}

func (c *Client) SetWebhook(webhookURL string) error {
	payload, err := json.Marshal(map[string]any{
		"url":                  webhookURL,
		"allowed_updates":      []string{"message"},
		"drop_pending_updates": true,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", c.token), bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.shortClient.Do(req)
	if err != nil {
		return c.redactError(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.redactError(err)
	}

	var body struct {
		OK          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}

	if err := json.Unmarshal(responseBody, &body); err != nil {
		return c.redactError(err)
	}

	if !body.OK {
		return &APIError{Code: body.ErrorCode, Description: body.Description}
	}

	return nil
}

func (c *Client) DeleteWebhook() error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/deleteWebhook?drop_pending_updates=true", c.token)

	resp, err := c.shortClient.Get(url)
	if err != nil {
		return c.redactError(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.redactError(err)
	}

	var payload struct {
		OK          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}

	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return c.redactError(err)
	}

	if !payload.OK {
		return &APIError{Code: payload.ErrorCode, Description: payload.Description}
	}

	return nil
}

func (c *Client) GetUpdates(offset int) ([]Update, error) {
	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/getUpdates?timeout=%d&offset=%d",
		c.token,
		longPollTimeoutSec,
		offset,
	)

	resp, err := c.longPollClient.Get(url)
	if err != nil {
		return nil, c.redactError(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, c.redactError(err)
	}

	var payload struct {
		OK          bool     `json:"ok"`
		Result      []Update `json:"result"`
		ErrorCode   int      `json:"error_code"`
		Description string   `json:"description"`
	}

	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return nil, c.redactError(err)
	}

	if !payload.OK {
		return nil, &APIError{Code: payload.ErrorCode, Description: payload.Description}
	}

	return payload.Result, nil
}

func (c *Client) redactError(err error) error {
	if err == nil {
		return nil
	}

	msg := err.Error()
	if c.token != "" {
		msg = strings.ReplaceAll(msg, c.token, "***")
	}

	return fmt.Errorf("%s", msg)
}

func ParseStartToken(text string) string {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "/start") {
		return ""
	}

	parts := strings.Fields(text)
	if len(parts) < 2 {
		return ""
	}

	return sanitizeStartToken(parts[1])
}

func sanitizeStartToken(token string) string {
	var builder strings.Builder

	for _, char := range strings.TrimSpace(token) {
		switch {
		case char >= '0' && char <= '9', char >= 'a' && char <= 'z', char >= 'A' && char <= 'Z', char == '_', char == '-':
			builder.WriteRune(char)
		}
	}

	return builder.String()
}
