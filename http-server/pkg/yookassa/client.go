package yookassa

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const apiBaseURL = "https://api.yookassa.ru/v3"

type Client struct {
	shopID     string
	secretKey  string
	returnURL  string
	httpClient *http.Client
}

func NewFromEnv() (*Client, error) {
	shopID := strings.TrimSpace(os.Getenv("YOOKASSA_SHOP_ID"))
	secretKey := strings.TrimSpace(os.Getenv("YOOKASSA_SECRET_KEY"))
	if shopID == "" || secretKey == "" {
		return nil, fmt.Errorf("yookassa: YOOKASSA_SHOP_ID and YOOKASSA_SECRET_KEY are required")
	}

	returnURL := strings.TrimSpace(os.Getenv("YOOKASSA_RETURN_URL"))
	if returnURL == "" {
		clientURL := strings.TrimRight(strings.TrimSpace(os.Getenv("CLIENT_URL")), "/")
		if clientURL != "" {
			returnURL = clientURL + "/pay"
		}
	}

	return &Client{
		shopID:    shopID,
		secretKey: secretKey,
		returnURL: returnURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (c *Client) ReturnURL() string {
	return c.returnURL
}

func (c *Client) Configured() bool {
	return c != nil && c.shopID != "" && c.secretKey != ""
}

func (c *Client) CreatePayment(ctx context.Context, idempotenceKey string, req CreatePaymentRequest) (*Payment, error) {
	var payment Payment
	if err := c.doJSON(ctx, http.MethodPost, "/payments", idempotenceKey, req, &payment); err != nil {
		return nil, err
	}

	return &payment, nil
}

func (c *Client) GetPayment(ctx context.Context, paymentID string) (*Payment, error) {
	var payment Payment
	if err := c.doJSON(ctx, http.MethodGet, "/payments/"+paymentID, "", nil, &payment); err != nil {
		return nil, err
	}

	return &payment, nil
}

func (c *Client) GetPaymentMethod(ctx context.Context, paymentMethodID string) (*PaymentMethod, error) {
	var method PaymentMethod
	if err := c.doJSON(ctx, http.MethodGet, "/payment_methods/"+paymentMethodID, "", nil, &method); err != nil {
		return nil, err
	}

	return &method, nil
}

func (c *Client) doJSON(ctx context.Context, method, path, idempotenceKey string, body any, out any) error {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}

		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, apiBaseURL+path, reader)
	if err != nil {
		return err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(c.shopID + ":" + c.secretKey))

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	if idempotenceKey != "" {
		req.Header.Set("Idempotence-Key", idempotenceKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("yookassa: %s %s", resp.Status, strings.TrimSpace(string(respBody)))
	}

	if out == nil {
		return nil
	}

	if err := json.Unmarshal(respBody, out); err != nil {
		return err
	}

	return nil
}
