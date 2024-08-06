package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
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
	DBName         int64
	ExpirationTime time.Duration
}

type SupabaseStorage struct {
	Url    string
	Key    string
	Header string
}

type JwtConfig struct {
	SecretKey       string
	TokenExpiration time.Duration
}

type MailConfig struct {
	Username string
	Password string
}

type AppConfig struct {
	DatabaseConfig        DatabaseConfig
	RedisConfig           RedisConfig
	Server                Server
	SupabaseStorageConfig SupabaseStorage
	JwtConfig             JwtConfig
	MailConfig            MailConfig
}

var Configs, _ = LoadConfig()

func LoadConfig() (*AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	config := &AppConfig{
		Server: Server{
			Port: GetEnv("PORT", "8080"),
			Host: GetEnv("HOST", "localhost"),
		},
		DatabaseConfig: DatabaseConfig{
			Port:     GetEnv("DB_PORT", "5432"),
			DBName:   GetEnv("DB_NAME", "postgres"),
			Domain:   GetEnv("DB_HOST", "vital-fawn-7347.6xw.aws-ap-southeast-1.cockroachlabs.cloud"),
			Password: GetEnv("DB_PASSWORD", "_oFbv9Yh03Ads4QzUYtiZg"),
			Username: GetEnv("DB_USER", "whisper"),
			Url:      GetEnv("DB_URL", "postgresql://whisper:_oFbv9Yh03Ads4QzUYtiZg@vital-fawn-7347.6xw.aws-ap-southeast-1.cockroachlabs.cloud:26257/postgres?sslmode=verify-full"),
		},
		RedisConfig: RedisConfig{
			Port:           GetEnv("REDIS_PORT", "6379"),
			Password:       GetEnv("REDIS_PASSWORD", "6379"),
			DBName:         GetEnvAsInt("REDIS_DB", 0),
			ExpirationTime: GetEnvAsDuration("REDIS_KEY_EXPIRATION", 10*time.Minute),
		},
		SupabaseStorageConfig: SupabaseStorage{
			Url:    GetEnv("SUPABASE_STORAGE_URL", "https://drqbnazyxxjvqapzcmkz.supabase.co/storage/v1"),
			Key:    GetEnv("SUPABASE_STORAGE_KEY", "Mat key"),
			Header: "",
		},
		JwtConfig: JwtConfig{
			SecretKey:       GetEnv("JWT_SECRET_KEY", "https://drqbnazyxxjvqapzcmkz.supabase.co/storage/v1"),
			TokenExpiration: GetEnvAsDuration("JWT_TOKEN_EXPIRATION", 1*time.Minute),
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

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}

func GetEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		i, err := time.ParseDuration(value)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}