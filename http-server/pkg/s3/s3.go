package s3

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	api        *s3.Client
	bucket     string
	publicBase string
}

func New() (*Client, error) {
	endpoint := endpointURL()
	bucket := strings.TrimSpace(os.Getenv("S3_BUCKET"))
	accessKey := strings.TrimSpace(os.Getenv("S3_ACCESS_KEY"))
	secretKey := strings.TrimSpace(os.Getenv("S3_SECRET_KEY"))
	region := strings.TrimSpace(os.Getenv("S3_REGION"))

	if endpoint == "" || bucket == "" || accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("s3 is not configured")
	}

	if region == "" {
		region = "us-east-1"
	}

	pathStyle := true
	if strings.TrimSpace(os.Getenv("S3_PATH_STYLE")) != "" {
		pathStyle = envBool("S3_PATH_STYLE")
	} else if envBool("S3_VIRTUAL_HOSTED_STYLE") {
		pathStyle = false
	}

	api := s3.NewFromConfig(aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
	}, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = pathStyle
	})

	return &Client{
		api:        api,
		bucket:     bucket,
		publicBase: strings.TrimRight(endpoint, "/") + "/" + bucket,
	}, nil
}

func (c *Client) Upload(ctx context.Context, key, contentType string, body io.Reader) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	_, err := c.api.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	return c.PublicURL(key), nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := c.api.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})

	return err
}

func (c *Client) PublicURL(key string) string {
	return strings.TrimRight(c.publicBase, "/") + "/" + strings.TrimLeft(key, "/")
}

func (c *Client) KeyFromURL(objectURL string) string {
	objectURL = strings.TrimSpace(objectURL)
	if objectURL == "" {
		return ""
	}

	parsed, err := url.Parse(objectURL)
	if err != nil {
		return ""
	}

	path := strings.TrimPrefix(parsed.Path, "/")
	if strings.HasPrefix(path, c.bucket+"/") {
		path = strings.TrimPrefix(path, c.bucket+"/")
	}

	return path
}

func endpointURL() string {
	endpoint := strings.TrimSpace(os.Getenv("S3_ENDPOINT"))
	endpoint = strings.TrimRight(endpoint, "/")
	if endpoint == "" {
		return ""
	}

	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return endpoint
	}

	scheme := "https"
	if strings.TrimSpace(os.Getenv("S3_USE_SSL")) != "" && !envBool("S3_USE_SSL") {
		scheme = "http"
	}

	return scheme + "://" + endpoint
}

func envBool(key string) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(key))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
