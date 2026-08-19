package maxbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	apiBaseURL         = "https://platform-api2.max.ru"
	longPollTimeoutSec = 30
)

type Client struct {
	token          string
	botName        string
	shortClient    *http.Client
	longPollClient *http.Client
}

type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return fmt.Sprintf("max API error %d", e.Code)
}

type Update struct {
	UpdateType string   `json:"update_type"`
	Timestamp  int64    `json:"timestamp"`
	ChatID     int64    `json:"chat_id"`
	Payload    string   `json:"payload"`
	User       *User    `json:"user"`
	Message    *Message `json:"message"`
}

type User struct {
	UserID   int64  `json:"user_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type Message struct {
	Sender *User        `json:"sender"`
	Body   *MessageBody `json:"body"`
}

type MessageBody struct {
	Text string `json:"text"`
}

type updatesResponse struct {
	Updates []Update `json:"updates"`
	Marker  *int64   `json:"marker"`
}

type subscriptionItem struct {
	URL string `json:"url"`
}

type subscriptionsResponse struct {
	Subscriptions []subscriptionItem `json:"subscriptions"`
}

func NewFromEnv() (*Client, error) {
	token := strings.TrimSpace(os.Getenv("MAX_BOT_TOKEN"))
	if token == "" {
		return nil, fmt.Errorf("MAX_BOT_TOKEN is not set")
	}

	botName := strings.TrimSpace(os.Getenv("MAX_BOT_USERNAME"))
	if botName == "" {
		return nil, fmt.Errorf("MAX_BOT_USERNAME is not set")
	}

	return &Client{
		token:   token,
		botName: strings.TrimPrefix(botName, "@"),
		shortClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		longPollClient: &http.Client{
			Timeout: time.Duration(longPollTimeoutSec+15) * time.Second,
		},
	}, nil
}

func (c *Client) BotURL(startToken string) string {
	return fmt.Sprintf("https://max.ru/%s?start=%s", c.botName, startToken)
}

func (c *Client) SendMessage(userID int64, text string) error {
	payload, err := json.Marshal(map[string]any{
		"text": text,
	})
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/messages?user_id=%d", apiBaseURL, userID)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	c.setAuth(req)
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
		return fmt.Errorf("max sendMessage failed: %s", strings.TrimSpace(string(responseBody)))
	}

	return nil
}

func (c *Client) GetUpdates(marker *int64) ([]Update, *int64, error) {
	query := url.Values{}
	query.Set("timeout", strconv.Itoa(longPollTimeoutSec))
	query.Set("types", "bot_started,message_created")
	if marker != nil {
		query.Set("marker", strconv.FormatInt(*marker, 10))
	}

	req, err := http.NewRequest(http.MethodGet, apiBaseURL+"/updates?"+query.Encode(), nil)
	if err != nil {
		return nil, nil, err
	}

	c.setAuth(req)

	resp, err := c.longPollClient.Do(req)
	if err != nil {
		return nil, nil, c.redactError(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, c.redactError(err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, nil, &APIError{Code: resp.StatusCode, Message: strings.TrimSpace(string(responseBody))}
	}

	var payload updatesResponse
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return nil, nil, c.redactError(err)
	}

	return payload.Updates, payload.Marker, nil
}

func (c *Client) SubscribeWebhook(webhookURL, secret string) error {
	payload, err := json.Marshal(map[string]any{
		"url":          webhookURL,
		"update_types": []string{"message_created", "bot_started"},
		"secret":       secret,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, apiBaseURL+"/subscriptions", bytes.NewReader(payload))
	if err != nil {
		return err
	}

	c.setAuth(req)
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
		return fmt.Errorf("max subscribe webhook failed: %s", strings.TrimSpace(string(responseBody)))
	}

	return nil
}

func (c *Client) DeleteSubscriptions() error {
	req, err := http.NewRequest(http.MethodGet, apiBaseURL+"/subscriptions", nil)
	if err != nil {
		return err
	}

	c.setAuth(req)

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
		return fmt.Errorf("max list subscriptions failed: %s", strings.TrimSpace(string(responseBody)))
	}

	var payload subscriptionsResponse
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return c.redactError(err)
	}

	for _, item := range payload.Subscriptions {
		if strings.TrimSpace(item.URL) == "" {
			continue
		}

		if err := c.deleteSubscription(item.URL); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) deleteSubscription(webhookURL string) error {
	endpoint := apiBaseURL + "/subscriptions?url=" + url.QueryEscape(webhookURL)
	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	c.setAuth(req)

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
		return fmt.Errorf("max delete subscription failed: %s", strings.TrimSpace(string(responseBody)))
	}

	return nil
}

func (c *Client) setAuth(req *http.Request) {
	req.Header.Set("Authorization", c.token)
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
