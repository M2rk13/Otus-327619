package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"
	"github.com/M2rk13/Otus-327619/internal/repository"
)

type StorageService struct {
	repo repository.Repository
}

func NewStorageService(repo repository.Repository) *StorageService {
	return &StorageService{repo: repo}
}

func (s *StorageService) StartStorageService(
	wg *sync.WaitGroup,
	ctx context.Context,
	requestChan <-chan *api.Request,
	responseChan <-chan *api.Response,
	logChan <-chan *log.ConversionLog,
) {
	wg.Add(3)

	go func() {
		defer wg.Done()

		for {
			select {
			case req, ok := <-requestChan:
				if !ok {
					fmt.Println("Request storage goroutine finished.")

					return
				}

				s.repo.CreateRequest(req)
			case <-ctx.Done():
				fmt.Println("Request storage goroutine stopped by context.")

				return
			}
		}
	}()

	go func() {
		defer wg.Done()

		for {
			select {
			case resp, ok := <-responseChan:
				if !ok {
					fmt.Println("Response storage goroutine finished.")

					return
				}

				s.repo.CreateResponse(resp)
			case <-ctx.Done():
				fmt.Println("Response storage goroutine stopped by context.")

				return
			}
		}
	}()

	go func() {
		defer wg.Done()

		for {
			select {
			case convLog, ok := <-logChan:
				if !ok {
					fmt.Println("Log storage goroutine finished.")

					return
				}

				s.repo.CreateConversionLog(convLog)
			case <-ctx.Done():
				fmt.Println("Log storage goroutine stopped by context.")

				return
			}
		}
	}()
}

func (s *StorageService) CreateRequest(req *api.Request) {
	s.repo.CreateRequest(req)
}

func (s *StorageService) GetRequestByID(id string) *api.Request {
	return s.repo.GetRequestByID(id)
}

func (s *StorageService) GetAllRequests() []*api.Request {
	return s.repo.GetAllRequests()
}

func (s *StorageService) UpdateRequest(req *api.Request) bool {
	return s.repo.UpdateRequest(req)
}

func (s *StorageService) DeleteRequest(id string) bool {
	return s.repo.DeleteRequest(id)
}

func (s *StorageService) CreateResponse(resp *api.Response) {
	s.repo.CreateResponse(resp)
}

func (s *StorageService) GetResponseByID(id string) *api.Response {
	return s.repo.GetResponseByID(id)
}

func (s *StorageService) GetAllResponses() []*api.Response {
	return s.repo.GetAllResponses()
}

func (s *StorageService) UpdateResponse(resp *api.Response) bool {
	return s.repo.UpdateResponse(resp)
}

func (s *StorageService) DeleteResponse(id string) bool {
	return s.repo.DeleteResponse(id)
}

func (s *StorageService) CreateConversionLog(log *log.ConversionLog) {
	s.repo.CreateConversionLog(log)
}

func (s *StorageService) GetConversionLogByID(id string) *log.ConversionLog {
	return s.repo.GetConversionLogByID(id)
}

func (s *StorageService) GetAllConversionLogs() []*log.ConversionLog {
	return s.repo.GetAllConversionLogs()
}

func (s *StorageService) UpdateConversionLog(log *log.ConversionLog) bool {
	return s.repo.UpdateConversionLog(log)
}

func (s *StorageService) DeleteConversionLog(id string) bool {
	return s.repo.DeleteConversionLog(id)
}
