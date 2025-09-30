package main

import (
	"context"
	"currency-converter/internal/api"
	"currency-converter/internal/config"
	"currency-converter/internal/handler"
	"currency-converter/internal/repository"
	"currency-converter/internal/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	if cfg.ExchangeRateAPIKey == "" {
		log.Fatalf("EXCHANGERATE_API_KEY is not set in .env file")
	}

	ctx := context.Background()
	pgRepo, err := repository.NewPostgresRepository(ctx, cfg.PostgresURL)

	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL")

	redisRepo, err := repository.NewRedisRepository(ctx, cfg)

	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	log.Println("Successfully connected to Redis")

	apiClient := api.NewClient(cfg.ExchangeRateAPIKey)
	convService := service.NewConversionService(apiClient, pgRepo, redisRepo)
	h := handler.NewHandler(convService)

	router := gin.Default()
	router.POST("/convert", h.Convert)

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Starting server on %s", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
