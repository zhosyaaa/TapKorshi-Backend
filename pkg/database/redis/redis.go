package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/zhosyaaa/RoommateTap/internal/config"
	"strconv"
)

func NewRedisClient(cfg config.RedisConfig) *redis.Client {
	db, err := strconv.Atoi(cfg.DB)
	if err != nil {
		return nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       db,
	})
	return rdb
}
