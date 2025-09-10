package rdb

import (
	"encoding/json"
	"roommates/utils"

	"github.com/redis/go-redis/v9"
)

type RedisHandler struct {
	redis *redis.Client
}

func New() *RedisHandler {
	redisAddress := utils.MustGetEnv("REDIS_ADDR")
	password := utils.MustGetEnv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: password,
		DB:       0,
		Protocol: 3,
	})

	return &RedisHandler{
		redis: client,
	}
}

// will marshal with json and panic on error
func Marshal(value any) []byte {
	redisValue, err := json.Marshal(value)
	if err != nil {
		log.Error().Err(err).Caller().
			Any("value", value).
			Msg("error during marshalling")
		panic(err)
	}
	return redisValue
}

// will unmarshal with json and panic on error
func Unmarshal(data []byte, value any) {
	err := json.Unmarshal(data, value)
	if err != nil {
		log.Error().Err(err).Caller().
			Any("data", data).
			Msg("error during unmarshalling")
		panic(err)
	}
}
