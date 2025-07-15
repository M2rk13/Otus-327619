package repository

import (
	"fmt"
	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"
	"sync"
	"time"
)

type repositoryItem[T any] struct {
	data     []T
	mu       sync.Mutex
	lastRead int
}

var (
	requestsItem  *repositoryItem[*api.Request]
	responsesItem *repositoryItem[*api.Response]
	logsItem      *repositoryItem[*log.ConversionLog]
)

func init() {
	requestsItem = &repositoryItem[*api.Request]{}
	responsesItem = &repositoryItem[*api.Response]{}
	logsItem = &repositoryItem[*log.ConversionLog]{}
}

func (ri *repositoryItem[T]) getNew() []T {
	ri.mu.Lock()
	defer ri.mu.Unlock()

	newItems := ri.data[ri.lastRead:]
	ri.lastRead = len(ri.data)

	return newItems
}

func (ri *repositoryItem[T]) add(item T) {
	ri.mu.Lock()
	defer ri.mu.Unlock()
	ri.data = append(ri.data, item)
}

func GetNewConversionRequests() []*api.Request {
	return requestsItem.getNew()
}

func GetNewConversionResponses() []*api.Response {
	return responsesItem.getNew()
}

func GetNewConversionLogs() []*log.ConversionLog {
	return logsItem.getNew()
}

func AddRequest(req *api.Request) {
	requestsItem.add(req)
	fmt.Printf("Added ConversionRequest: From=%s, To=%s, Amount=%.2f\n", req.From, req.To, req.Amount)
}

func AddResponse(resp *api.Response) {
	responsesItem.add(resp)
	fmt.Printf("Added ConversionResponse: Success=%t, Result=%.2f\n", resp.Success, resp.Result)
}

func AddLog(convLog *log.ConversionLog) {
	logsItem.add(convLog)
	fmt.Printf("Added ConversionLog: ID=%s, Timestamp=%s\n", convLog.ID(), convLog.Timestamp().Format(time.RFC3339))
}
