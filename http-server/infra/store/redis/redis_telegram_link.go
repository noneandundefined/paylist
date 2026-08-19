package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const telegramLinkTTL = 30 * time.Minute

func telegramLinkKey(token string) string {
	return fmt.Sprintf("telegram_link:%s", token)
}

func RedisTelegramLinkCreate(userUuid string) (string, error) {
	tokenBytes := make([]byte, 16)

	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	token := hex.EncodeToString(tokenBytes)
	key := telegramLinkKey(token)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Set(ctx, key, userUuid, telegramLinkTTL).Err(); err != nil {
		return "", err
	}

	return token, nil
}

func RedisTelegramLinkGet(token string) (string, error) {
	key := telegramLinkKey(token)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userUuid, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return userUuid, nil
}

func RedisTelegramLinkDelete(token string) error {
	key := telegramLinkKey(token)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Del(ctx, key).Err()
}
