package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"pm/domain/repository/caches"
	"time"
)

type RedisCacheRepository struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedisCacheRepository(rdb *redis.Client, ctx context.Context) caches.RedisCacheRepository {
	if rdb == nil {
		fmt.Println("implementations/cache: error nil redis.Client")
		return nil
	}
	return &RedisCacheRepository{rdb, ctx}
}

func (redisDriver *RedisCacheRepository) GetHash(key string, property string, src interface{}) {
	if redisDriver.rdb == nil {
		fmt.Println("RedisDriver not found")
		return
	}
	val, err := redisDriver.rdb.HGet(redisDriver.ctx, key, property).Result()
	if err != nil {
		return
	}
	errU := json.Unmarshal([]byte(val), &src)
	if errU != nil {
		fmt.Println("error!!!!")
	}
}

func (redisDriver *RedisCacheRepository) DeleteHash(key string, property string) error {
	if redisDriver.rdb == nil {
		return fmt.Errorf("RedisDriver not found")
	}
	_, err := redisDriver.rdb.HDel(redisDriver.ctx, key, property).Result()
	if err != nil {
		return err
	}
	return nil
}

func (redisDriver *RedisCacheRepository) SetHashObject(key string, property string, object interface{}) error {
	if redisDriver.rdb == nil {
		return fmt.Errorf("RedisDriver not found")
	}
	return redisDriver.rdb.HSet(redisDriver.ctx, key, property, object).Err()
}

func (redisDriver *RedisCacheRepository) SetExpireKey(key string, expiration time.Duration) error {
	if redisDriver.rdb == nil {
		return fmt.Errorf("RedisDriver not found")
	}
	return redisDriver.rdb.Expire(redisDriver.ctx, key, expiration).Err()
}