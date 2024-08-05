package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"pm/domain/repository/caches"
	"pm/infrastructure/implementations/cache"
	"pm/infrastructure/persistences/base"
	"time"
)

var rd2Driver *redis.Client = nil
var rd2Repo caches.RedisCacheRepository = nil
var ctx context.Context = nil
var expirationTime time.Duration = time.Minute

func InitCacheHelper(p *base.Persistence) {
	if p.Redis.RedisDB == nil {
		fmt.Println("error cannot get redisClient from persistence")
		return
	}
	rd2Repo = cache.NewRedisCacheRepository(p.Redis.RedisDB, p.Ctx)
	rd2Driver = p.Redis.RedisDB
	ctx = p.Ctx
	expirationTime = p.Redis.KeyExpirationTime
}

func GetAllHashGeneric[M any](key string, src *map[string]M) {
	if rd2Driver == nil {
		return
	}

	val, err := rd2Driver.HGetAll(ctx, key).Result()
	if err != nil {
		return
	}

	*src = make(map[string]M)
	for k, v := range val {
		var item M
		if err := json.Unmarshal([]byte(v), &item); err != nil {
			fmt.Printf("Error unmarshaling value for key %s: %v\n", k, err)
			continue
		}
		(*src)[k] = item
	}

	return
}

func RedisSetHashGenericKeySlice[M any](path string, objects []M, getID func(M) int64, exp time.Duration) error {
	if rd2Driver == nil {
		return fmt.Errorf("RedisDriver not found")
	}
	mapErrs := make(map[interface{}]string)
	for _, obj := range objects {
		key := getID(obj)
		marshalObject, _ := sonic.Marshal(obj)
		errR := rd2Repo.SetHashObject(path, fmt.Sprintf("%v", key), marshalObject)
		if errR != nil {
			errSet := fmt.Errorf("RedisSetGenericHashSlice", "Error to save key:", key, "Object", string(marshalObject))
			mapErrs[key] = errSet.Error()
			continue
		}
	}

	if len(mapErrs) > 0 {
		return fmt.Errorf("error set: %v", mapErrs)
	}

	// Set expiration if needed (only once)
	if exp > 0 {
		if err := rd2Repo.SetExpireKey(path, exp); err != nil {
			return fmt.Errorf("RedisSetGenericHashSlice", "Error setting expiration:", path, "err", err)
		}
	}
	return nil
}

func RedisSetHashGenericKey[M any](path string, key string, object M, exp time.Duration) error {
	if rd2Driver == nil {
		return fmt.Errorf("RedisDriver not found")
	}
	marshalObject, _ := sonic.Marshal(object)
	errR := rd2Repo.SetHashObject(path, key, marshalObject)
	if errR != nil {
		return fmt.Errorf("RedisSet Generic", "Error to save key:", key, "Object", string(marshalObject))
	}
	// Set expiration if needed
	if exp > 0 {
		if err := rd2Repo.SetExpireKey(path, exp); err != nil {
			return fmt.Errorf("RedisSetGenericHash", "Error setting expiration:", key, "err", err.Error())
		}
	}
	return nil
}

func RedisRemoveHashGenericKey(path string, key string) error {
	if rd2Driver == nil {
		return fmt.Errorf("RedisDriver not found")
	}

	if err := rd2Repo.DeleteHash(path, key); err != nil {
		return fmt.Errorf("RedisRemoveGenericHash", "Error to delete key:", key, "err", err.Error())
	}
	return nil
}

func RedisGetHashGenericKey[M any](path string, key string, object *M) {
	if rd2Driver == nil {
		return
	}
	rd2Repo.GetHash(path, key, object)
}