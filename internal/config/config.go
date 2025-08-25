package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
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
	URI      string
	Database string
}

type RedisConfig struct {
	Addr     string
	Password string
	TTL      time.Duration
}

var (
	AppCfg   AppConfig
	AdminCfg AdminConfig
	FileCfg  FileConfig
	MongoCfg MongoConfig
	RedisCfg RedisConfig
)

func LoadAll() {
	AppCfg = loadAppConfig()
	AdminCfg = loadAdminConfig()
	FileCfg = loadFileConfig()
	MongoCfg = loadMongoConfig()
	RedisCfg = loadRedisConfig()

	if AdminCfg.Login == "" || AdminCfg.Password == "" || AdminCfg.JwtKey == "" {
		panic("Admin credentials are required")
	}

	log.Println("INFO: Application configuration loaded.")
}

func loadAppConfig() AppConfig {
	storageType := os.Getenv("STORAGE_TYPE")

	if storageType == "" {
		storageType = "file"
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
		URI:      os.Getenv("MONGO_URI"),
		Database: os.Getenv("MONGO_DATABASE"),
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
