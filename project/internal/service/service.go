package service

import (
	"context"
	"fmt"

	"currency-converter/internal/model"
)

type APIClient interface {
	GetConversionRate(from, to string, amount float64) (float64, float64, error)
}

type HistorySaver interface {
	SaveConversion(ctx context.Context, history *model.ConversionHistory) error
}

type Logger interface {
	LogConversion(ctx context.Context, from, to string, amount, result float64) error
}

type ConversionService struct {
	apiClient    APIClient
	historySaver HistorySaver
	logger       Logger
}

func NewConversionService(apiClient APIClient, historySaver HistorySaver, logger Logger) *ConversionService {
	return &ConversionService{
		apiClient:    apiClient,
		historySaver: historySaver,
		logger:       logger,
	}
}

func (s *ConversionService) PerformConversion(ctx context.Context, req *model.ConversionAPIRequest) (*model.ConversionAPIResponse, error) {
	result, rate, err := s.apiClient.GetConversionRate(req.From, req.To, req.Amount)

	if err != nil {
		return nil, fmt.Errorf("failed to get data from external api: %w", err)
	}

	history := &model.ConversionHistory{
		From:   req.From,
		To:     req.To,
		Amount: req.Amount,
		Result: result,
		Rate:   rate,
	}
	if err := s.historySaver.SaveConversion(ctx, history); err != nil {
		fmt.Printf("warning: failed to save conversion history: %v\n", err)
	}

	if err := s.logger.LogConversion(ctx, req.From, req.To, req.Amount, result); err != nil {
		fmt.Printf("warning: failed to log conversion: %v\n", err)
	}

	response := &model.ConversionAPIResponse{
		From:   req.From,
		To:     req.To,
		Amount: req.Amount,
		Result: result,
		Rate:   rate,
	}

	return response, nil
}
