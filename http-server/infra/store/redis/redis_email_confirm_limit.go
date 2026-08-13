package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"paylist.server/infra/constants"
)

func emailConfirmSendLimitKey(email string) string {
	return fmt.Sprintf("email_confirm_send_limit:%s", email)
}

func emailConfirmPendingKey(email string) string {
	return fmt.Sprintf("email_confirm_pending:%s", email)
}

func RedisEmailConfirmCheckAndIncrement(email string, maxPerDay int) (bool, error) {
	if err := ensureClient(); err != nil {
		return false, err
	}

	key := emailConfirmSendLimitKey(email)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		if err := client.Expire(ctx, key, constants.EmailConfirmSendLimitTTL).Err(); err != nil {
			return false, err
		}
	}

	return count <= int64(maxPerDay), nil
}

func RedisEmailConfirmPendingSet(email string) error {
	if err := ensureClient(); err != nil {
		return err
	}

	key := emailConfirmPendingKey(email)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Set(ctx, key, "1", constants.EmailConfirmPendingTTL).Err()
}

func RedisEmailConfirmPendingCheck(email string) (bool, error) {
	if err := ensureClient(); err != nil {
		return false, err
	}

	key := emailConfirmPendingKey(email)

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

func RedisEmailConfirmPendingClear(email string) error {
	if err := ensureClient(); err != nil {
		return err
	}

	key := emailConfirmPendingKey(email)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Del(ctx, key).Err()
}
