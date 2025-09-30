package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort         string
	PostgresURL        string
	RedisAddr          string
	RedisPass          string
	RedisDB            int
	ExchangeRateAPIKey string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	pgURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	redisAddr := fmt.Sprintf("%s:%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)

	return &Config{
		ServerPort:         os.Getenv("SERVER_PORT"),
		PostgresURL:        pgURL,
		RedisAddr:          redisAddr,
		RedisPass:          os.Getenv("REDIS_PASSWORD"),
		RedisDB:            0,
		ExchangeRateAPIKey: os.Getenv("EXCHANGERATE_API_KEY"),
	}, nil
}
