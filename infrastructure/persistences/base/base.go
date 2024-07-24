package base

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pm/infrastructure/config"
	db2 "pm/infrastructure/persistences/base/db"
	rd2 "pm/infrastructure/persistences/base/redis"
	"time"
)

type Persistence struct {
	GormDB              *gorm.DB
	RedisDB             *redis.Client
	Ctx                 context.Context
	RedisExpirationTime time.Duration
}

func InitPersistence(appConfig *config.AppConfig) *Persistence {
	var persistence Persistence
	gormDBConfig := appConfig.DatabaseConfig
	gormDB, err := db2.SetupDatabase(db2.GetDSN(gormDBConfig.Username, gormDBConfig.Password, gormDBConfig.Domain, gormDBConfig.Port, gormDBConfig.DBName))
	if err != nil {
		fmt.Println("error connecting to database", err)
	}
	persistence.GormDB = gormDB

	redisClient, err := rd2.NewRedisClient(appConfig)
	if err != nil {
		fmt.Println("error connecting to redis", err)
	}
	persistence.RedisDB = nil

	return &Persistence{
		GormDB:              gormDB,
		RedisDB:             redisClient.Client,
		Ctx:                 redisClient.Ctx,
		RedisExpirationTime: redisClient.ExpirationTime,
	}
}