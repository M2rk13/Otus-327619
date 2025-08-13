package config

import (
	"log"
	"os"
	"path/filepath"
)

type FileConfig struct {
	RequestsFilePath  string
	ResponsesFilePath string
	LogsFilePath      string
}

var (
	FileCfg FileConfig
)

func LoadFileConfig() FileConfig {
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

	log.Println("INFO: Application configuration loaded.")

	return cfg
}
