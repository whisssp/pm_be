package base

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	storage_go "github.com/supabase-community/storage-go"
	"gorm.io/gorm"
	"pm/infrastructure/config"
	db2 "pm/infrastructure/persistences/base/db"
	rd2 "pm/infrastructure/persistences/base/redis"
	"pm/infrastructure/persistences/base/supabase"
	"time"
)

type Persistence struct {
	GormDB              *gorm.DB
	RedisDB             *redis.Client
	SupabaseStorage     *storage_go.Client
	Ctx                 context.Context
	RedisExpirationTime time.Duration
}

func InitPersistence(appConfig *config.AppConfig) *Persistence {
	persistence := &Persistence{
		GormDB:              nil,
		RedisDB:             nil,
		SupabaseStorage:     nil,
		Ctx:                 context.Background(),
		RedisExpirationTime: 10 * time.Minute,
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
	persistence.RedisDB = redisClient.Client

	supabaseStorage := supabase.NewSupabaseStorage(appConfig)
	if supabaseStorage == nil {
		fmt.Println("error connecting to redis", err)
	}
	persistence.SupabaseStorage = supabaseStorage.StorageClient

	return persistence
}