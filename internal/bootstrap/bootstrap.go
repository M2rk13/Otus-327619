package bootstrap

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/M2rk13/Otus-327619/internal/config"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	config.AdminCfg = config.LoadAdminConfig()
	config.FileCfg = config.LoadFileConfig()

	if config.AdminCfg.Login == "" || config.AdminCfg.Password == "" || config.AdminCfg.JwtKey == "" {
		panic("Admin credentials is required")
	}
}
