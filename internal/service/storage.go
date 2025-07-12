package service

import (
	"fmt"
	"sync"

	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"
	"github.com/M2rk13/Otus-327619/internal/repository"
)

func StartStorageService(
	wg *sync.WaitGroup,
	requestChan <-chan *api.Request,
	responseChan <-chan *api.Response,
	logChan <-chan *log.ConversionLog,
) {
	wg.Add(3)

	go func() {
		defer wg.Done()

		for req := range requestChan {
			repository.AddRequest(req)
		}

		fmt.Println("Request storage goroutine finished.")
	}()

	go func() {
		defer wg.Done()

		for resp := range responseChan {
			repository.AddResponse(resp)
		}

		fmt.Println("Response storage goroutine finished.")
	}()

	go func() {
		defer wg.Done()

		for convLog := range logChan {
			repository.AddLog(convLog)
		}

		fmt.Println("Log storage goroutine finished.")
	}()
}
