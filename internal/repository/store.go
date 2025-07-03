package store

import (
	"fmt"
	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"
	"sync"
	"time"
)

var (
	RequestChan  chan *api.Request
	ResponseChan chan *api.Response
	LogChan      chan *log.ConversionLog
)

var (
	RequestChanState  int
	ResponseChanState int
	LogChanState      int
)

var (
	conversionRequests  []*api.Request
	conversionResponses []*api.Response
	conversionLogs      []*log.ConversionLog
)

var (
	muRequests  sync.Mutex
	muResponses sync.Mutex
	muLogs      sync.Mutex
)

var (
	prevRequestsLen  int
	prevResponsesLen int
	prevLogsLen      int
)

func init() {
	RequestChan = make(chan *api.Request, 10)
	RequestChanState = 1
	prevRequestsLen = 0

	ResponseChan = make(chan *api.Response, 10)
	ResponseChanState = 1
	prevResponsesLen = 0

	LogChan = make(chan *log.ConversionLog, 10)
	LogChanState = 1
	prevLogsLen = 0
}

func GetNewConversionRequests() []*api.Request {
	muRequests.Lock()
	defer muRequests.Unlock()

	newItems := conversionRequests[prevRequestsLen:]
	prevRequestsLen = len(conversionRequests)

	return newItems
}

func GetNewConversionResponses() []*api.Response {
	muResponses.Lock()
	defer muResponses.Unlock()

	newItems := conversionResponses[prevResponsesLen:]
	prevResponsesLen = len(conversionResponses)

	return newItems
}

func GetNewConversionLogs() []*log.ConversionLog {
	muLogs.Lock()
	defer muLogs.Unlock()

	newItems := conversionLogs[prevLogsLen:]
	prevLogsLen = len(conversionLogs)

	return newItems
}

func StartStorageGoRoutines(wg *sync.WaitGroup) {
	wg.Add(3)

	go func() {
		defer wg.Done()

		for req := range RequestChan {
			muRequests.Lock()
			conversionRequests = append(conversionRequests, req)
			muRequests.Unlock()

			fmt.Printf("Added ConversionRequest: From=%s, To=%s, Amount=%.2f\n", req.From, req.To, req.Amount)
		}

		fmt.Println("Request storage goroutine finished.")
	}()

	go func() {
		defer wg.Done()

		for resp := range ResponseChan {
			muResponses.Lock()
			conversionResponses = append(conversionResponses, resp)
			muResponses.Unlock()

			fmt.Printf("Added ConversionResponse: Success=%t, Result=%.2f\n", resp.Success, resp.Result)
		}

		fmt.Println("Response storage goroutine finished.")
	}()

	go func() {
		defer wg.Done()

		for convLog := range LogChan {
			muLogs.Lock()
			conversionLogs = append(conversionLogs, convLog)
			muLogs.Unlock()

			fmt.Printf("Added ConversionLog: ID=%s, Timestamp=%s\n", convLog.ID(), convLog.Timestamp().Format(time.RFC3339))
		}

		fmt.Println("Log storage goroutine finished.")
	}()
}
