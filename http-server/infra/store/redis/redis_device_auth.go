package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"paylist.server/config"
	"paylist.server/infra/constants"
	"paylist.server/pkg"
	"paylist.server/types"
)

func deviceAuthKey(key string) string {
	return fmt.Sprintf("device_auth:%s", key)
}

func RedisDeviceAuthCreate(s *types.DeviceAuthSession) (string, error) {
	/* Base data */
	s.CreatedAt = time.Now()

	sess, err := config.JSON.Marshal(s)
	if err != nil {
		return "", err
	}

	sessionId, err := pkg.GenerateString(100)
	if err != nil {
		return "", err
	}

	key := deviceAuthKey(sessionId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := client.Exists(ctx, key).Result()
	if err != nil {
		return "", err
	}

	if exists == 1 {
		if err := client.Expire(ctx, key, constants.REDIS_DESKTOP_AUTH_TTL).Err(); err != nil {
			return "", err
		}

		return sessionId, nil
	}

	if err := client.Set(ctx, key, sess, constants.REDIS_DESKTOP_AUTH_TTL).Err(); err != nil {
		return "", err
	}

	return sessionId, nil
}

func RedisDeviceAuthUpdate(sessionId string, s *types.DeviceAuthSession) error {
	key := deviceAuthKey(sessionId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sess, err := config.JSON.Marshal(s)
	if err != nil {
		return err
	}

	ttl, err := client.TTL(ctx, key).Result()
	if err != nil {
		return err
	}

	if ttl <= 0 {
		return errors.New("desktop auth session expired")
	}

	return client.Set(ctx, key, sess, constants.REDIS_DESKTOP_AUTH_TTL).Err()
}

func RedisDeviceAuthGet(sessionId string) (*types.DeviceAuthSession, error) {
	key := deviceAuthKey(sessionId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	value, err := client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var session types.DeviceAuthSession
	if err := config.JSON.Unmarshal([]byte(value), &session); err != nil {
		return nil, err
	}

	return &session, nil
}
