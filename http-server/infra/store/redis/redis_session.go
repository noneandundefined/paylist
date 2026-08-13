package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"paylist.server/config"
	"paylist.server/infra/constants"
	"paylist.server/infra/encryption"
	"paylist.server/types"
)

func sessionKey(key string) string {
	return fmt.Sprintf("session:%s", key)
}

func deviceSessionsKey(userUuid string) string {
	return fmt.Sprintf("device_sessions:%s", userUuid)
}

func RedisDeviceSessionAdd(userUuid, sessionId string) error {
	if userUuid == "" || sessionId == "" {
		return nil
	}

	key := deviceSessionsKey(userUuid)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipe := client.Pipeline()

	pipe.SAdd(ctx, key, sessionId)
	pipe.Expire(ctx, key, constants.REDIS_SESSION_TTL)

	_, err := pipe.Exec(ctx)

	return err
}

func RedisDeviceSessionRemove(userUuid, sessionId string) error {
	if userUuid == "" || sessionId == "" {
		return nil
	}

	key := deviceSessionsKey(userUuid)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.SRem(ctx, key, sessionId).Err()
}

func RedisDeviceSessionIDs(userUuid string) ([]string, error) {
	key := deviceSessionsKey(userUuid)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.SMembers(ctx, key).Result()
}

func RedisDeviceSessionsDeleteAll(userUuid string) error {
	sessionIds, err := RedisDeviceSessionIDs(userUuid)
	if err != nil {
		return err
	}

	for _, sessionId := range sessionIds {
		_ = RedisSessionDelete(sessionId)
	}

	key := deviceSessionsKey(userUuid)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Del(ctx, key).Err()
}

func RedisSessionCreate(s *types.Session) (string, error) {
	/* Base data */
	s.CreatedAt = time.Now()

	sess, err := config.JSON.Marshal(s)
	if err != nil {
		return "", err
	}

	sessionId, err := encryption.Encrypt(string(sess))
	if err != nil {
		return "", err
	}

	key := sessionKey(sessionId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s.Platform == "device" && s.DeviceId != "" {
		existingSessionIds, err := RedisDeviceSessionIDs(s.UserUuid)
		if err != nil {
			return "", err
		}

		for _, existingSessionId := range existingSessionIds {
			existingSession, err := RedisSessionGet(existingSessionId)
			if err != nil {
				return "", err
			}

			if existingSession == nil {
				_ = RedisDeviceSessionRemove(s.UserUuid, existingSessionId)
				continue
			}

			if existingSession.Platform == "device" && existingSession.DeviceId == s.DeviceId {
				if err := RedisSessionDelete(existingSessionId); err != nil {
					return "", err
				}
			}
		}
	}

	exists, err := client.Exists(ctx, key).Result()
	if err != nil {
		return "", err
	}

	if exists == 1 {
		if err := client.Expire(ctx, key, constants.REDIS_SESSION_TTL).Err(); err != nil {
			return "", err
		}

		if err := RedisDeviceSessionAdd(s.UserUuid, sessionId); err != nil {
			return "", err
		}

		return sessionId, nil
	}

	if err := client.Set(ctx, key, sess, constants.REDIS_SESSION_TTL).Err(); err != nil {
		return "", err
	}

	if err := RedisDeviceSessionAdd(s.UserUuid, sessionId); err != nil {
		return "", err
	}

	return sessionId, nil
}

func RedisSessionGet(sessionId string) (*types.Session, error) {
	key := sessionKey(sessionId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	value, err := client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var session types.Session
	if err := config.JSON.Unmarshal([]byte(value), &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func RedisSessionRefresh(sessionId string) error {
	key := sessionKey(sessionId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ttl, err := client.TTL(ctx, key).Result()
	if err != nil {
		return err
	}

	if ttl > 12*time.Hour {
		return nil
	}

	_, err = client.Expire(ctx, key, constants.REDIS_SESSION_TTL).Result()
	return err
}

func RedisSessionDelete(sessionId string) error {
	session, _ := RedisSessionGet(sessionId)

	key := sessionKey(sessionId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Del(ctx, key).Err(); err != nil {
		return err
	}

	if session != nil {
		_ = RedisDeviceSessionRemove(session.UserUuid, sessionId)
	}

	return nil
}
