package rdb

import (
	"context"
	"errors"
	"roommates/logger"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var log = logger.RedisLoggger

// Redis key start for user sessions
const KSession = "session:"
const EUserSession = 48 * time.Hour

type UserSessionValue struct {
	UserID   string `redis:"user_id" json:"user_id"`
	Username string `redis:"username" json:"username"`
}

// can panic
func (r *RedisHandler) CreateUserSession(ctx context.Context, sessionValue UserSessionValue) uuid.UUID {
	key := uuid.New() // NewRandom generates V4
	rKey := KSession + key.String()
	rValue := Marshal(sessionValue)

	cmd := r.redis.SetNX(ctx, rKey, rValue, EUserSession)
	err := cmd.Err()
	if err != nil {
		log.Error().Err(err).Str("key", rKey).Caller().Msg("error during CreateUserSession")
		panic(err)
	}
	return key
}

// Key is just string form of UUID
//
// No need to add redis "topic"/"key start" (KSession) in front of `key`
func (r *RedisHandler) GetUserSession(ctx context.Context, key string) (*UserSessionValue, error) {
	rKey := KSession + key
	value := &UserSessionValue{}

	cmd := r.redis.Get(ctx, rKey)
	err := cmd.Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		log.Error().Err(err).Str("key", rKey).Caller().Msg("error during GetUserSession")
		return nil, err
	}

	Unmarshal([]byte(cmd.Val()), value)
	return value, nil
}

// Key is just string form of UUID
//
// No need to add redis "topic"/"key start" (KSession) in front of `key`
func (r *RedisHandler) DeleteUserSession(ctx context.Context, key string) error {
	rKey := KSession + key

	// .Result() is the same as getting .Val() and .Err()
	cmd := r.redis.Del(ctx, rKey)
	err := cmd.Err()
	if err != nil {
		log.Error().Err(err).Caller().Msg("error during DeleteUserSession")
		return err
	}
	return nil
}
