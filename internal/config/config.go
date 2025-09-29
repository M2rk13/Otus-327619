package config

import (
	_ "github.com/jackc/pgx/v4/stdlib"

	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/M2rk13/Otus-327619/internal/enum"
)

type AppConfig struct {
	StorageType string
}

type AdminConfig struct {
	Login    string
	Password string
	JwtKey   string
}

type FileConfig struct {
	RequestsFilePath  string
	ResponsesFilePath string
	LogsFilePath      string
}

type MongoConfig struct {
	MogoUri string
	MongoDb string
}

type RedisConfig struct {
	Addr     string
	Password string
	TTL      time.Duration
}

type PostgresConfig struct {
	PostgresUri string
}

var (
	AppCfg      AppConfig
	AdminCfg    AdminConfig
	FileCfg     FileConfig
	MongoCfg    MongoConfig
	RedisCfg    RedisConfig
	PostgresCfg PostgresConfig
)

func LoadAll() {
	AppCfg = loadAppConfig()
	AdminCfg = loadAdminConfig()
	FileCfg = loadFileConfig()
	MongoCfg = loadMongoConfig()
	RedisCfg = loadRedisConfig()
	PostgresCfg = loadPostgresConfig()

	if AdminCfg.Login == "" || AdminCfg.Password == "" || AdminCfg.JwtKey == "" {
		panic("Admin credentials are required")
	}

	log.Println("INFO: Application configuration loaded.")
}

func loadAppConfig() AppConfig {
	storageType := os.Getenv("STORAGE_TYPE")

	if storageType == "" {
		storageType = enum.File
	}

	return AppConfig{
		StorageType: storageType,
	}
}

func loadAdminConfig() AdminConfig {
	return AdminConfig{
		Login:    os.Getenv("LOGIN"),
		Password: os.Getenv("PASSWORD"),
		JwtKey:   os.Getenv("JWT_KEY"),
	}
}

func loadFileConfig() FileConfig {
	cfg := FileConfig{
		RequestsFilePath:  os.Getenv("REQUESTS_FILE_PATH"),
		ResponsesFilePath: os.Getenv("RESPONSES_FILE_PATH"),
		LogsFilePath:      os.Getenv("LOGS_FILE_PATH"),
	}

	if cfg.RequestsFilePath == "" {
		cfg.RequestsFilePath = filepath.Join("data", "requests.json")
	}

	if cfg.ResponsesFilePath == "" {
		cfg.ResponsesFilePath = filepath.Join("data", "responses.json")
	}

	if cfg.LogsFilePath == "" {
		cfg.LogsFilePath = filepath.Join("data", "logs.json")
	}

	return cfg
}

func loadMongoConfig() MongoConfig {
	return MongoConfig{
		MogoUri: os.Getenv("MONGO_URI"),
		MongoDb: os.Getenv("MONGO_DATABASE"),
	}
}

func loadRedisConfig() RedisConfig {
	ttlStr := os.Getenv("REDIS_TTL_SECONDS")
	ttl, err := strconv.Atoi(ttlStr)

	if err != nil || ttl <= 0 {
		ttl = 3600
	}

	return RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		TTL:      time.Duration(ttl) * time.Second,
	}
}

func loadPostgresConfig() PostgresConfig {
	return PostgresConfig{
		PostgresUri: os.Getenv("POSTGRES_DSN"),
	}
}
