package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/M2rk13/Otus-327619/internal/repository"
)

type LoggerService struct {
	repo repository.Repository
}

func NewLoggerService(repo repository.Repository) *LoggerService {
	return &LoggerService{repo: repo}
}

func (l *LoggerService) StartSliceLogger(
	wg *sync.WaitGroup,
	ctx context.Context,
	requestChanState *int,
	responseChanState *int,
	logChanState *int,
) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		fmt.Println("Slice logger started.")

		for {
			select {
			case <-ticker.C:
				newRequests := l.repo.GetNewConversionRequests()

				if len(newRequests) > 0 {
					fmt.Println("--- New Conversion Requests ---")

					for _, req := range newRequests {
						fmt.Printf("Request: From=%s, To=%s, Amount=%.2f\n", req.From, req.To, req.Amount)
					}
				}

				newResponses := l.repo.GetNewConversionResponses()

				if len(newResponses) > 0 {
					fmt.Println("--- New Conversion Responses ---")

					for _, resp := range newResponses {
						fmt.Printf("Response: Success=%t, Result=%.2f\n", resp.Success, resp.Result)
					}
				}

				newLogs := l.repo.GetNewConversionLogs()

				if len(newLogs) > 0 {
					fmt.Println("--- New Conversion Logs ---")

					for _, logItem := range newLogs {
						fmt.Printf(
							"  Log: GetId=%s, GetTimestamp=%s, RequestFrom=%s, ResponseResult=%.2f\n",
							logItem.Id,
							logItem.Timestamp.Format(time.RFC3339),
							logItem.Request.From,
							logItem.Response.Result)
					}
				}

				if *requestChanState == 0 && *responseChanState == 0 && *logChanState == 0 {
					time.Sleep(250 * time.Millisecond)
					fmt.Println("All channels closed. Shutting down logger.")

					return
				}
			case <-ctx.Done():
				fmt.Println("Logger stopped due context cancel.")

				return
			}
		}
	}()
}
