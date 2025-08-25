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

type MockRepository struct {
	mu        sync.Mutex
	requests  map[string]*api.Request
	responses map[string]*api.Response
	logs      map[string]*log.ConversionLog
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		requests:  make(map[string]*api.Request),
		responses: make(map[string]*api.Response),
		logs:      make(map[string]*log.ConversionLog),
	}
}

func (m *MockRepository) CreateRequest(req *api.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if req.Id == "" {
		req.Id = uuid.New().String()
	}

	m.requests[req.Id] = req
}

func (m *MockRepository) GetRequestByID(id string) *api.Request {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.requests[id]
}

func (m *MockRepository) GetAllRequests() []*api.Request {
	m.mu.Lock()
	defer m.mu.Unlock()

	var all []*api.Request

	for _, r := range m.requests {
		all = append(all, r)
	}

	return all
}
func (m *MockRepository) UpdateRequest(req *api.Request) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.requests[req.Id]; ok {
		m.requests[req.Id] = req

		return true
	}

	return false
}

func (m *MockRepository) DeleteRequest(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.requests[id]; ok {
		delete(m.requests, id)

		return true
	}

	return false
}

func (m *MockRepository) CreateResponse(resp *api.Response) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if resp.Id == "" {
		resp.Id = uuid.New().String()
	}

	m.responses[resp.Id] = resp
}

func (m *MockRepository) GetResponseByID(id string) *api.Response {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.responses[id]
}

func (m *MockRepository) GetAllResponses() []*api.Response {
	m.mu.Lock()
	defer m.mu.Unlock()

	var all []*api.Response

	for _, r := range m.responses {
		all = append(all, r)
	}

	return all
}

func (m *MockRepository) UpdateResponse(resp *api.Response) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.responses[resp.Id]; ok {
		m.responses[resp.Id] = resp

		return true
	}

	return false
}

func (m *MockRepository) DeleteResponse(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.responses[id]; ok {
		delete(m.responses, id)

		return true
	}

	return false
}

func (m *MockRepository) CreateConversionLog(item *log.ConversionLog) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if item.Id == "" {
		item.Id = uuid.New().String()
	}

	m.logs[item.Id] = item
}

func (m *MockRepository) GetConversionLogByID(id string) *log.ConversionLog {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.logs[id]
}

func (m *MockRepository) GetAllConversionLogs() []*log.ConversionLog {
	m.mu.Lock()
	defer m.mu.Unlock()

	var all []*log.ConversionLog

	for _, r := range m.logs {
		all = append(all, r)
	}

	return all
}

func (m *MockRepository) UpdateConversionLog(item *log.ConversionLog) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.logs[item.Id]; ok {
		m.logs[item.Id] = item

		return true
	}

	return false
}

func (m *MockRepository) DeleteConversionLog(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.logs[id]; ok {
		delete(m.logs, id)

		return true
	}

	return false
}

func (m *MockRepository) GetNewConversionRequests() []*api.Request   { return nil }
func (m *MockRepository) GetNewConversionResponses() []*api.Response { return nil }
func (m *MockRepository) GetNewConversionLogs() []*log.ConversionLog { return nil }

var _ repository.Repository = (*MockRepository)(nil)

func eventually(t *testing.T, ok func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)

	for time.Now().Before(deadline) {
		if ok() {
			return
		}

		time.Sleep(5 * time.Millisecond)
	}

	t.Fatal("condition not met before deadline")
}

func TestStorageService_CRUD_Wrappers(t *testing.T) {
	m := NewMockRepository()
	s := NewStorageService(m)

	req := &api.Request{From: "USD", To: "EUR", Amount: 100}
	s.CreateRequest(req)

	if req.Id == "" {
		t.Fatal("CreateRequest must set Id")
	}

	gotReq := s.GetRequestByID(req.Id)

	if gotReq == nil || gotReq.Amount != 100 {
		t.Fatalf("GetRequestByID failed: got=%v", gotReq)
	}

	if len(s.GetAllRequests()) != 1 {
		t.Fatal("GetAllRequests should return 1")
	}

	resp := &api.Response{Success: true, Result: 123.45}
	s.CreateResponse(resp)

	if resp.Id == "" {
		t.Fatal("CreateResponse must set Id")
	}

	gotResp := s.GetResponseByID(resp.Id)

	if gotResp == nil || !gotResp.Success {
		t.Fatalf("GetResponseByID failed: got=%v", gotResp)
	}

	if len(s.GetAllResponses()) != 1 {
		t.Fatal("GetAllResponses should return 1")
	}

	cl := &log.ConversionLog{}
	s.CreateConversionLog(cl)

	if cl.Id == "" {
		t.Fatal("CreateConversionLog must set Id")
	}

	if len(s.GetAllConversionLogs()) != 1 {
		t.Fatal("GetAllConversionLogs should return 1")
	}

	reqUpd := &api.Request{Id: req.Id, From: "USD", To: "RUB", Amount: 200}

	if !s.UpdateRequest(reqUpd) {
		t.Fatal("UpdateRequest should return true")
	}

	if s.GetRequestByID(req.Id).Amount != 200 {
		t.Fatal("UpdateRequest did not apply changes")
	}

	respUpd := &api.Response{Id: resp.Id, Success: false, Result: 0}

	if !s.UpdateResponse(respUpd) {
		t.Fatal("UpdateResponse should return true")
	}

	if s.GetResponseByID(resp.Id).Success {
		t.Fatal("UpdateResponse did not apply changes")
	}

	logUpd := &log.ConversionLog{Id: cl.Id}

	if !s.UpdateConversionLog(logUpd) {
		t.Fatal("UpdateConversionLog should return true")
	}

	if s.UpdateRequest(&api.Request{Id: "nope"}) {
		t.Fatal("UpdateRequest should return false for unknown id")
	}

	if s.DeleteRequest("nope") {
		t.Fatal("DeleteRequest should return false for unknown id")
	}

	if s.UpdateResponse(&api.Response{Id: "nope"}) {
		t.Fatal("UpdateResponse should return false for unknown id")
	}

	if s.DeleteResponse("nope") {
		t.Fatal("DeleteResponse should return false for unknown id")
	}

	if s.UpdateConversionLog(&log.ConversionLog{Id: "nope"}) {
		t.Fatal("UpdateConversionLog should return false for unknown id")
	}

	if s.DeleteConversionLog("nope") {
		t.Fatal("DeleteConversionLog should return false for unknown id")
	}

	if !s.DeleteRequest(req.Id) || s.GetRequestByID(req.Id) != nil {
		t.Fatal("DeleteRequest failed")
	}

	if !s.DeleteResponse(resp.Id) || s.GetResponseByID(resp.Id) != nil {
		t.Fatal("DeleteResponse failed")
	}

	if !s.DeleteConversionLog(cl.Id) || s.GetConversionLogByID(cl.Id) != nil {
		t.Fatal("DeleteConversionLog failed")
	}
}

func TestStartStorageService_CloseChannels(t *testing.T) {
	mockRepo := NewMockRepository()
	svc := NewStorageService(mockRepo)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reqCh := make(chan *api.Request, 2)
	respCh := make(chan *api.Response, 2)
	logCh := make(chan *log.ConversionLog, 2)

	svc.StartStorageService(&wg, ctx, reqCh, respCh, logCh)

	reqCh <- &api.Request{From: "USD", To: "EUR", Amount: 1}
	respCh <- &api.Response{Success: true, Result: 42}
	logCh <- &log.ConversionLog{}

	close(reqCh)
	close(respCh)
	close(logCh)

	wg.Wait()

	if got := len(mockRepo.requests); got != 1 {
		t.Fatalf("requests saved = %d, want 1", got)
	}

	if got := len(mockRepo.responses); got != 1 {
		t.Fatalf("responses saved = %d, want 1", got)
	}

	if got := len(mockRepo.logs); got != 1 {
		t.Fatalf("logs saved = %d, want 1", got)
	}
}

func TestStartStorageService_ContextCancel(t *testing.T) {
	mockRepo := NewMockRepository()
	svc := NewStorageService(mockRepo)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	reqCh := make(chan *api.Request, 1)
	respCh := make(chan *api.Response, 1)
	logCh := make(chan *log.ConversionLog, 1)

	svc.StartStorageService(&wg, ctx, reqCh, respCh, logCh)

	reqCh <- &api.Request{From: "GBP", To: "USD", Amount: 5}
	respCh <- &api.Response{Success: true, Result: 1.23}
	logCh <- &log.ConversionLog{}

	eventually(t, func() bool {
		return len(mockRepo.requests) == 1 &&
			len(mockRepo.responses) == 1 &&
			len(mockRepo.logs) == 1
	})

	cancel()
	wg.Wait()
}

func TestStartStorageService_ConcurrentPush(t *testing.T) {
	mockRepo := NewMockRepository()
	svc := NewStorageService(mockRepo)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reqCh := make(chan *api.Request, 1000)
	respCh := make(chan *api.Response, 1000)
	logCh := make(chan *log.ConversionLog, 1000)

	svc.StartStorageService(&wg, ctx, reqCh, respCh, logCh)

	var producers sync.WaitGroup
	N := 200

	producers.Add(3)

	go func() {
		defer producers.Done()

		for i := 0; i < N; i++ {
			reqCh <- &api.Request{From: "A", To: "B", Amount: float64(i)}
		}
	}()

	go func() {
		defer producers.Done()

		for i := 0; i < N; i++ {
			respCh <- &api.Response{Success: i%2 == 0, Result: float64(i)}
		}
	}()

	go func() {
		defer producers.Done()

		for i := 0; i < N; i++ {
			logCh <- &log.ConversionLog{}
		}
	}()

	producers.Wait()

	eventually(t, func() bool {
		return len(mockRepo.requests) == N &&
			len(mockRepo.responses) == N &&
			len(mockRepo.logs) == N
	})

	cancel()
	wg.Wait()
}

func TestNewStorageService_Implements(t *testing.T) {
	var _ repository.Repository = NewMockRepository()
	_ = NewStorageService(NewMockRepository())
}
