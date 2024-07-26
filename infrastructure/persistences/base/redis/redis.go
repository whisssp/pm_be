package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"pm/infrastructure/config"
	"time"
)

type RedisCache struct {
	Client         *redis.Client
	Ctx            context.Context
	ExpirationTime time.Duration
}

func NewRedisClient(appConfig *config.AppConfig) (*RedisCache, error) {
	redisConfig := appConfig.RedisConfig
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", appConfig.Server.Host, redisConfig.Port),
		DB:       int(appConfig.RedisConfig.DBName),
		Password: redisConfig.Password,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Printf("Error connect Redis: %v", err)
		return &RedisCache{nil, ctx, appConfig.RedisConfig.ExpirationTime}, err
	}
	fmt.Println("Connected to Redis successfully")
	return &RedisCache{client, ctx, appConfig.RedisConfig.ExpirationTime}, nil
}