package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func emailConfirmedKey(key string) string {
	return fmt.Sprintf("email_confirmed:%s", key)
}

func RedisEmailConfirmedMark(email string) error {
	key := emailConfirmedKey(email)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Set(ctx, key, "1", 30*time.Minute).Err()
}

func RedisEmailConfirmedCheck(email string) (bool, error) {
	key := emailConfirmedKey(email)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	value, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return value != "", nil
}
