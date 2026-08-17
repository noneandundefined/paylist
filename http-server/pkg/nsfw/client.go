package nsfw

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"paylist.server/infra/logger"
)

const defaultTimeout = 5 * time.Second

var (
	ErrRejected     = errors.New("nsfw: content rejected")
	ErrUnavailable  = errors.New("nsfw: service unavailable")
	ErrInvalidImage = errors.New("nsfw: invalid image")
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type moderateResponse struct {
	NSFW   bool               `json:"nsfw"`
	Scores map[string]float64 `json:"scores"`
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
	return New(os.Getenv("NSFW_URL"))
}

func (c *Client) Configured() bool {
	return c != nil && c.baseURL != ""
}

func CheckFromEnv(ctx context.Context, payload []byte, contentType string) error {
	client := NewFromEnv()
	if !client.Configured() {
		if strings.EqualFold(strings.TrimSpace(os.Getenv("GO_ENV")), "DEV") {
			logger.Warning("nsfw: NSFW_URL is empty, skipping avatar moderation")
			return nil
		}

		return ErrUnavailable
	}

	return client.Check(ctx, payload, contentType)
}

func (c *Client) Check(ctx context.Context, payload []byte, contentType string) error {
	if !c.Configured() {
		return ErrUnavailable
	}

	if len(payload) == 0 {
		return ErrInvalidImage
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
		defer cancel()
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("image", "avatar"+extensionFor(contentType))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnavailable, err.Error())
	}

	if _, err := part.Write(payload); err != nil {
		return fmt.Errorf("%w: %s", ErrUnavailable, err.Error())
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("%w: %s", ErrUnavailable, err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/moderate", &body)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnavailable, err.Error())
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnavailable, err.Error())
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnavailable, err.Error())
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var parsed moderateResponse
		if err := json.Unmarshal(raw, &parsed); err != nil {
			return fmt.Errorf("%w: invalid response", ErrUnavailable)
		}

		if parsed.NSFW {
			return ErrRejected
		}

		return nil
	case http.StatusBadRequest, http.StatusRequestEntityTooLarge, http.StatusUnsupportedMediaType:
		return ErrInvalidImage
	default:
		return ErrUnavailable
	}
}

func extensionFor(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	default:
		return ".bin"
	}
}
