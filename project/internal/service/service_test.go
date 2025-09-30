package service

import (
	"context"
	"errors"
	"testing"

	"currency-converter/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAPIClient struct {
	mock.Mock
}

func (m *MockAPIClient) GetConversionRate(from, to string, amount float64) (float64, float64, error) {
	args := m.Called(from, to, amount)
	return args.Get(0).(float64), args.Get(1).(float64), args.Error(2)
}

type MockHistorySaver struct {
	mock.Mock
}

func (m *MockHistorySaver) SaveConversion(ctx context.Context, history *model.ConversionHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) LogConversion(ctx context.Context, from, to string, amount, result float64) error {
	args := m.Called(ctx, from, to, amount, result)
	return args.Error(0)
}

func TestConversionService_PerformConversion_Success(t *testing.T) {
	mockAPI := new(MockAPIClient)
	mockSaver := new(MockHistorySaver)
	mockLogger := new(MockLogger)
	service := NewConversionService(mockAPI, mockSaver, mockLogger)

	req := &model.ConversionAPIRequest{From: "USD", To: "EUR", Amount: 100}
	expectedResult := 92.5
	expectedRate := 0.925

	mockAPI.On("GetConversionRate", "USD", "EUR", 100.0).Return(expectedResult, expectedRate, nil)
	mockSaver.On("SaveConversion", mock.Anything, mock.AnythingOfType("*model.ConversionHistory")).Return(nil)
	mockLogger.On("LogConversion", mock.Anything, "USD", "EUR", 100.0, expectedResult).Return(nil)

	result, err := service.PerformConversion(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResult, result.Result)
	assert.Equal(t, expectedRate, result.Rate)

	mockAPI.AssertExpectations(t)
	mockSaver.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestConversionService_PerformConversion_APIClientError(t *testing.T) {
	mockAPI := new(MockAPIClient)
	mockSaver := new(MockHistorySaver)
	mockLogger := new(MockLogger)
	service := NewConversionService(mockAPI, mockSaver, mockLogger)

	req := &model.ConversionAPIRequest{From: "USD", To: "EUR", Amount: 100}
	apiError := errors.New("API is down")

	mockAPI.On("GetConversionRate", "USD", "EUR", 100.0).Return(0.0, 0.0, apiError)
	result, err := service.PerformConversion(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "API is down")

	mockAPI.AssertExpectations(t)
	mockSaver.AssertNotCalled(t, "SaveConversion", mock.Anything, mock.Anything)
	mockLogger.AssertNotCalled(t, "LogConversion", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestConversionService_PerformConversion_HistorySaverError(t *testing.T) {
	mockAPI := new(MockAPIClient)
	mockSaver := new(MockHistorySaver)
	mockLogger := new(MockLogger)
	service := NewConversionService(mockAPI, mockSaver, mockLogger)

	req := &model.ConversionAPIRequest{From: "USD", To: "EUR", Amount: 100}
	expectedResult := 92.5
	expectedRate := 0.925
	dbError := errors.New("connection to postgres failed")

	mockAPI.On("GetConversionRate", "USD", "EUR", 100.0).Return(expectedResult, expectedRate, nil)
	mockSaver.On("SaveConversion", mock.Anything, mock.AnythingOfType("*model.ConversionHistory")).Return(dbError)
	mockLogger.On("LogConversion", mock.Anything, "USD", "EUR", 100.0, expectedResult).Return(nil)

	result, err := service.PerformConversion(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResult, result.Result)

	mockAPI.AssertExpectations(t)
	mockSaver.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestConversionService_PerformConversion_LoggerError(t *testing.T) {
	mockAPI := new(MockAPIClient)
	mockSaver := new(MockHistorySaver)
	mockLogger := new(MockLogger)
	service := NewConversionService(mockAPI, mockSaver, mockLogger)

	req := &model.ConversionAPIRequest{From: "USD", To: "EUR", Amount: 100}
	expectedResult := 92.5
	expectedRate := 0.925
	redisError := errors.New("connection to redis failed")

	mockAPI.On("GetConversionRate", "USD", "EUR", 100.0).Return(expectedResult, expectedRate, nil)
	mockSaver.On("SaveConversion", mock.Anything, mock.AnythingOfType("*model.ConversionHistory")).Return(nil)
	mockLogger.On("LogConversion", mock.Anything, "USD", "EUR", 100.0, expectedResult).Return(redisError)

	result, err := service.PerformConversion(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResult, result.Result)

	mockAPI.AssertExpectations(t)
	mockSaver.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestConversionService_PerformConversion_SaverAndLoggerErrors(t *testing.T) {
	mockAPI := new(MockAPIClient)
	mockSaver := new(MockHistorySaver)
	mockLogger := new(MockLogger)
	service := NewConversionService(mockAPI, mockSaver, mockLogger)

	req := &model.ConversionAPIRequest{From: "USD", To: "EUR", Amount: 100}
	expectedResult := 92.5
	expectedRate := 0.925
	dbError := errors.New("db is down")
	redisError := errors.New("redis is down")

	mockAPI.On("GetConversionRate", "USD", "EUR", 100.0).Return(expectedResult, expectedRate, nil)
	mockSaver.On("SaveConversion", mock.Anything, mock.AnythingOfType("*model.ConversionHistory")).Return(dbError)
	mockLogger.On("LogConversion", mock.Anything, "USD", "EUR", 100.0, expectedResult).Return(redisError)

	result, err := service.PerformConversion(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResult, result.Result)

	mockAPI.AssertExpectations(t)
	mockSaver.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
