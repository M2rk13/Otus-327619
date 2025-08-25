package bootstrap

import (
	"log"

	"github.com/M2rk13/Otus-327619/internal/config"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	config.LoadAll()
}
