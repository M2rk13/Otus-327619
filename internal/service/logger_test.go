package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"
	"github.com/M2rk13/Otus-327619/internal/repository"

	"github.com/google/uuid"
)

type loggerMockRepo struct {
	mu       sync.Mutex
	onceReq  []*api.Request
	onceResp []*api.Response
	onceLogs []*log.ConversionLog
}

func (m *loggerMockRepo) CreateRequest(*api.Request)                     {}
func (m *loggerMockRepo) GetRequestByID(string) *api.Request             { return nil }
func (m *loggerMockRepo) GetAllRequests() []*api.Request                 { return nil }
func (m *loggerMockRepo) UpdateRequest(*api.Request) bool                { return false }
func (m *loggerMockRepo) DeleteRequest(string) bool                      { return false }
func (m *loggerMockRepo) CreateResponse(*api.Response)                   {}
func (m *loggerMockRepo) GetResponseByID(string) *api.Response           { return nil }
func (m *loggerMockRepo) GetAllResponses() []*api.Response               { return nil }
func (m *loggerMockRepo) UpdateResponse(*api.Response) bool              { return false }
func (m *loggerMockRepo) DeleteResponse(string) bool                     { return false }
func (m *loggerMockRepo) CreateConversionLog(*log.ConversionLog)         {}
func (m *loggerMockRepo) GetConversionLogByID(string) *log.ConversionLog { return nil }
func (m *loggerMockRepo) GetAllConversionLogs() []*log.ConversionLog     { return nil }
func (m *loggerMockRepo) UpdateConversionLog(*log.ConversionLog) bool    { return false }
func (m *loggerMockRepo) DeleteConversionLog(string) bool                { return false }

func (m *loggerMockRepo) GetNewConversionRequests() []*api.Request {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := m.onceReq
	m.onceReq = nil

	return out
}

func (m *loggerMockRepo) GetNewConversionResponses() []*api.Response {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := m.onceResp
	m.onceResp = nil

	return out
}

func (m *loggerMockRepo) GetNewConversionLogs() []*log.ConversionLog {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := m.onceLogs
	m.onceLogs = nil

	return out
}

var _ repository.Repository = (*loggerMockRepo)(nil)

func TestNewLoggerService(t *testing.T) {
	l := NewLoggerService(&loggerMockRepo{})

	if l == nil {
		t.Fatal("NewLoggerService returned nil")
	}
}

func TestLoggerStartSliceLoggerTickerAndShutdown(t *testing.T) {
	r := &loggerMockRepo{}
	r.onceReq = []*api.Request{{Id: uuid.New().String(), From: "GBP", To: "USD", Amount: 5}}
	r.onceResp = []*api.Response{{Id: uuid.New().String(), Success: true, Result: 1.23}}
	r.onceLogs = []*log.ConversionLog{log.NewConversionLog(uuid.New().String(),
		api.Request{Id: uuid.New().String(), From: "A", To: "B", Amount: 1},
		api.Response{Id: uuid.New().String(), Success: true, Result: 2.0},
	)}

	l := NewLoggerService(r)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reqState, respState, logState := 0, 0, 0
	l.StartSliceLogger(&wg, ctx, &reqState, &respState, &logState)

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()

	select {
	case <-done:
		// ок
	case <-time.After(3 * time.Second):
		t.Fatal("logger did not stop on closed states")
	}
}

func TestLoggerStartSliceLoggerContextCancel(t *testing.T) {
	r := &loggerMockRepo{}
	l := NewLoggerService(r)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	reqState, respState, logState := 1, 1, 1

	l.StartSliceLogger(&wg, ctx, &reqState, &respState, &logState)

	cancel()

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()

	select {
	case <-done:
		// ок
	case <-time.After(2 * time.Second):
		t.Fatal("logger did not stop on context cancel")
	}
}
