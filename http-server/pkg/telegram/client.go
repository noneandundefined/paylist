package telegram

import (
	"bytes"
	"encoding/json"
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
	payload := map[string]any{
		"chat_id": chatID,
		"text":    text,
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
		return fmt.Errorf("telegram sendMessage failed: %s", strings.TrimSpace(string(responseBody)))
	}

	return nil
}

func (c *Client) DeleteWebhook() error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/deleteWebhook?drop_pending_updates=false", c.token)

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

	return strings.TrimSpace(parts[1])
}
