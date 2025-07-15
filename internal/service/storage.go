package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"
	"github.com/M2rk13/Otus-327619/internal/repository"
)

func StartStorageService(
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

				repository.AddRequest(req)
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

				repository.AddResponse(resp)
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

				repository.AddLog(convLog)
			case <-ctx.Done():
				fmt.Println("Log storage goroutine stopped by context.")

				return
			}
		}
	}()
}
