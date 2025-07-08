package repository

import (
	"fmt"
	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"
	"sync"
	"time"
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
	prevRequestsLen = 0
	prevResponsesLen = 0
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

func AddRequest(req *api.Request) {
	muRequests.Lock()
	defer muRequests.Unlock()

	conversionRequests = append(conversionRequests, req)
	fmt.Printf("Added ConversionRequest: From=%s, To=%s, Amount=%.2f\n", req.From, req.To, req.Amount)
}

func AddResponse(resp *api.Response) {
	muResponses.Lock()
	defer muResponses.Unlock()

	conversionResponses = append(conversionResponses, resp)
	fmt.Printf("Added ConversionResponse: Success=%t, Result=%.2f\n", resp.Success, resp.Result)
}

func AddLog(convLog *log.ConversionLog) {
	muLogs.Lock()
	defer muLogs.Unlock()

	conversionLogs = append(conversionLogs, convLog)
	fmt.Printf("Added ConversionLog: ID=%s, Timestamp=%s\n", convLog.ID(), convLog.Timestamp().Format(time.RFC3339))
}
