package redisClient

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/scalable-websocket/config"
)

var ctx = context.Background()

func NewRedisClient(cfg *config.Config) *redis.Client {
	options := &redis.Options{
		Addr:     cfg.Redis.Addr,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
	}
	redisClient := redis.NewClient(options)
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis.")
	return redisClient
}
