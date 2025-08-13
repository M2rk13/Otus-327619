package config

import (
	"os"
)

type AdminConfig struct {
	Login    string
	Password string
	JwtKey   string
}

var (
	AdminCfg AdminConfig
)

func LoadAdminConfig() AdminConfig {
	cfg := AdminConfig{
		Login:    os.Getenv("LOGIN"),
		Password: os.Getenv("PASSWORD"),
		JwtKey:   os.Getenv("JWT_KEY"),
	}

	if cfg.Login == "" || cfg.Password == "" || cfg.JwtKey == "" {
		panic("Admin credentials is required")
	}

	return cfg
}
