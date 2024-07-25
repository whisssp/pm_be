package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"pm/infrastructure/config"
	"strconv"
	"time"
)

type RedisCache struct {
	Client         *redis.Client
	Ctx            context.Context
	ExpirationTime time.Duration
}

func NewRedisClient(appConfig *config.AppConfig) (*RedisCache, error) {
	redisConfig := appConfig.RedisConfig
	dbInt, _ := strconv.Atoi(redisConfig.DBName)
	ctx := context.Background()
	exp, _ := time.ParseDuration(redisConfig.ExpirationTime)
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", appConfig.Server.Host, redisConfig.Port),
		DB:       dbInt,
		Password: redisConfig.Password,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Printf("Error connect Redis: %v", err)
		return &RedisCache{nil, ctx, exp}, err
	}
	fmt.Println("Connected to Redis successfully")
	return &RedisCache{client, ctx, exp}, nil
}