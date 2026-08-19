package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const maxLinkTTL = 30 * time.Minute

func maxLinkKey(token string) string {
	return fmt.Sprintf("max_link:%s", token)
}

func RedisMaxLinkCreate(userUuid, language string) (string, error) {
	tokenBytes := make([]byte, 16)

	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	token := hex.EncodeToString(tokenBytes)
	key := maxLinkKey(token)
	value := userUuid + "|" + language

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Set(ctx, key, value, maxLinkTTL).Err(); err != nil {
		return "", err
	}

	return token, nil
}

func RedisMaxLinkGet(token string) (string, string, error) {
	key := maxLinkKey(token)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	value, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", "", nil
	}

	if err != nil {
		return "", "", err
	}

	userUuid, language, _ := strings.Cut(value, "|")
	if language == "" {
		language = "en"
	}

	return userUuid, language, nil
}

func RedisMaxLinkDelete(token string) error {
	key := maxLinkKey(token)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Del(ctx, key).Err()
}
