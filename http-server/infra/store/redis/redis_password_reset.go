package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"paylist.server/infra/constants"
)

func passwordResetSendLimitKey(email string) string {
	return fmt.Sprintf("password_reset_send_limit:%s", email)
}

func passwordResetTokenKey(userUuid string) string {
	return fmt.Sprintf("password_reset_token:%s", userUuid)
}

func RedisPasswordResetCheckAndIncrement(email string, maxPerDay int) (bool, error) {
	if err := ensureClient(); err != nil {
		return false, err
	}

	key := passwordResetSendLimitKey(email)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		if err := client.Expire(ctx, key, constants.PasswordResetSendLimitTTL).Err(); err != nil {
			return false, err
		}
	}

	return count <= int64(maxPerDay), nil
}

func RedisPasswordResetTokenSet(userUuid, signature string) error {
	if err := ensureClient(); err != nil {
		return err
	}

	key := passwordResetTokenKey(userUuid)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Set(ctx, key, signature, constants.PasswordResetTokenTTL).Err()
}

func RedisPasswordResetTokenMatch(userUuid, signature string) (bool, error) {
	if err := ensureClient(); err != nil {
		return false, err
	}

	key := passwordResetTokenKey(userUuid)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	value, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return value == signature, nil
}

func RedisPasswordResetTokenClear(userUuid string) error {
	if err := ensureClient(); err != nil {
		return err
	}

	key := passwordResetTokenKey(userUuid)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Del(ctx, key).Err()
}
