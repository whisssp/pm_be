package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

//type AppConfig struct {
//	DatabaseConfig struct {
//		Port     string `yaml:"port"`
//		DBName   string `yaml:"name"`
//		Domain   string `yaml:"domain"`
//		Password string `yaml:"password"`
//		Username string `yaml:"username"`
//		Url      string `yaml:"url"`
//	} `yaml:"database"`
//
//	RedisConfig struct {
//		Port           string `yaml:"port"`
//		Password       string `yaml:"password"`
//		DBName         string `yaml:"name"`
//		ExpirationTime string `yaml:"expirationTime"`
//	} `yaml:"redis"`
//
//	Server struct {
//		Port string `yaml:"port"`
//	} `yaml:"server"`
//
//	Host string `yaml:"host"`
//}

type Server struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Port     string
	DBName   string
	Domain   string
	Password string
	Username string
	Url      string
}

type RedisConfig struct {
	Port           string
	Password       string
	DBName         string
	ExpirationTime string
}

type SupabaseStorage struct {
	Url    string
	Key    string
	Header string
}

type AppConfig struct {
	DatabaseConfig  DatabaseConfig
	RedisConfig     RedisConfig
	Server          Server
	SupabaseStorage SupabaseStorage
}

func LoadConfig() (*AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	config := &AppConfig{
		DatabaseConfig: DatabaseConfig{
			Port:     getEnv("DB_PORT", "5432"),
			DBName:   getEnv("DB_NAME", "postgres"),
			Domain:   getEnv("DB_HOST", "vital-fawn-7347.6xw.aws-ap-southeast-1.cockroachlabs.cloud"),
			Password: getEnv("DB_PASSWORD", "_oFbv9Yh03Ads4QzUYtiZg"),
			Username: getEnv("DB_USER", "whisper"),
			Url:      getEnv("DB_URL", "postgresql://whisper:_oFbv9Yh03Ads4QzUYtiZg@vital-fawn-7347.6xw.aws-ap-southeast-1.cockroachlabs.cloud:26257/postgres?sslmode=verify-full"),
		},
		RedisConfig: RedisConfig{
			Port:           getEnv("REDIS_PORT", "6379"),
			Password:       getEnv("REDIS_PASSWORD", "6379"),
			DBName:         getEnv("REDIS_DB", "0"),
			ExpirationTime: getEnv("REDIS_KEY_EXPIRATION", "600000000000"),
		},
		Server: Server{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "localhost"),
		},
		SupabaseStorage: SupabaseStorage{
			Url:    getEnv("SUPABASE_STORAGE_URL", "https://drqbnazyxxjvqapzcmkz.supabase.co/storage/v1"),
			Key:    getEnv("SUPABASE_STORAGE_KEY", "Mat key"),
			Header: "",
		},
	}

	//file, err := os.Open("./infrastructure/config/application.yml")
	//if err != nil {
	//	return nil, err
	//}
	//defer file.Close()
	//
	//d := yaml.NewDecoder(file)
	//
	//if err := d.Decode(&config); err != nil {
	//	return nil, err
	//}
	//return config, nil

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}