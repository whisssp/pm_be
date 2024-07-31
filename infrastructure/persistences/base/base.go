package base

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	storage_go "github.com/supabase-community/storage-go"
	"gorm.io/gorm"
	"pm/infrastructure/config"
	"pm/infrastructure/implementations/loggers"
	db2 "pm/infrastructure/persistences/base/db"
	"pm/infrastructure/persistences/base/logger"
	rd2 "pm/infrastructure/persistences/base/redis"
	"pm/infrastructure/persistences/base/supabase"
	"time"
)

type Persistence struct {
	Ctx             context.Context
	GormDB          *gorm.DB
	Redis           RedisPersistence
	SupabaseStorage *storage_go.Client
	Logger          *loggers.LoggerRepo
}

type RedisPersistence struct {
	RedisDB           *redis.Client
	KeyExpirationTime time.Duration
}

func InitPersistence(appConfig *config.AppConfig) *Persistence {
	persistence := &Persistence{
		GormDB: nil,
		Redis: RedisPersistence{
			RedisDB:           nil,
			KeyExpirationTime: appConfig.RedisConfig.ExpirationTime,
		},
		SupabaseStorage: nil,
		Ctx:             context.Background(),
	}
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
	persistence.Redis.RedisDB = redisClient.Client

	supabaseStorage := supabase.NewSupabaseStorage(appConfig)
	if supabaseStorage == nil {
		fmt.Println("error connecting to redis", err)
	}
	persistence.SupabaseStorage = supabaseStorage.StorageClient

	logger, err := logger.NewLogger()
	if err != nil {
		fmt.Println("error initializing logger", err)
	}
	persistence.Logger = logger
	return persistence
}