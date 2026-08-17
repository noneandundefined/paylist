package profanity

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"paylist.server/infra/locale"
	"paylist.server/infra/logger"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/types"
)

const defaultTimeout = 5 * time.Second

var (
	ErrRejected    = errors.New("profanity: content rejected")
	ErrUnavailable = errors.New("profanity: service unavailable")
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type moderateRequest struct {
	Texts []string `json:"texts"`
}

type moderateResponse struct {
	Profane bool `json:"profane"`
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func NewFromEnv() *Client {
	return New(os.Getenv("PROFANITY_URL"))
}

func (c *Client) Configured() bool {
	return c != nil && c.baseURL != ""
}

func Pointers(values ...*string) []string {
	texts := make([]string, 0, len(values))
	for _, value := range values {
		if value == nil {
			continue
		}

		trimmed := strings.TrimSpace(*value)
		if trimmed == "" {
			continue
		}

		texts = append(texts, trimmed)
	}

	return texts
}

func CheckFromEnv(ctx context.Context, texts ...string) error {
	client := NewFromEnv()
	if !client.Configured() {
		if strings.EqualFold(strings.TrimSpace(os.Getenv("GO_ENV")), "DEV") {
			logger.Warning("profanity: PROFANITY_URL is empty, skipping text moderation")
			return nil
		}

		return ErrUnavailable
	}

	return client.Check(ctx, texts...)
}

func Reject(ctx context.Context, tr locale.Translator, reason string, texts ...string) error {
	err := CheckFromEnv(ctx, texts...)
	if err == nil {
		return nil
	}

	if errors.Is(err, ErrRejected) {
		logger.Moderation("user_uuid=%s reason=%s req=%s", userUUIDFromContext(ctx), reason, requestID(ctx))
		return httperr.BadRequest(tr.TErr("error.text-profanity"))
	}

	logger.Error("profanity check failed req={%s}: %s", requestID(ctx), err.Error())
	return httperr.ServiceUnavailable(tr.TErr("error.text-moderation-unavailable"))
}

func (c *Client) Check(ctx context.Context, texts ...string) error {
	if !c.Configured() {
		return ErrUnavailable
	}

	filtered := make([]string, 0, len(texts))
	for _, text := range texts {
		trimmed := strings.TrimSpace(text)
		if trimmed == "" {
			continue
		}

		filtered = append(filtered, trimmed)
	}

	if len(filtered) == 0 {
		return nil
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
		defer cancel()
	}

	body, err := json.Marshal(moderateRequest{Texts: filtered})
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnavailable, err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/moderate", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnavailable, err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnavailable, err.Error())
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnavailable, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return ErrUnavailable
	}

	var parsed moderateResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return fmt.Errorf("%w: invalid response", ErrUnavailable)
	}

	if parsed.Profane {
		return ErrRejected
	}

	return nil
}

func userUUIDFromContext(ctx context.Context) string {
	identity, ok := ctx.Value("identity").(*types.AuthToken)
	if !ok || identity == nil {
		return ""
	}

	return identity.User.UserUUID
}

func requestID(ctx context.Context) string {
	value, _ := ctx.Value("XREQID").(string)
	return value
}
