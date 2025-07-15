package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/M2rk13/Otus-327619/internal/config"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/M2rk13/Otus-327619/internal/model/api"
	logmodel "github.com/M2rk13/Otus-327619/internal/model/log"
)

type repositoryItem[T any] struct {
	data     []T
	mu       sync.Mutex
	lastRead int
	filePath string
	file     *os.File
}

var (
	requestsItem  *repositoryItem[*api.Request]
	responsesItem *repositoryItem[*api.Response]
	logsItem      *repositoryItem[*logmodel.ConversionLog]
)

func init() {
	fileConfig := config.LoadConfig()

	requestsItem = &repositoryItem[*api.Request]{filePath: fileConfig.RequestsFilePath}
	responsesItem = &repositoryItem[*api.Response]{filePath: fileConfig.ResponsesFilePath}
	logsItem = &repositoryItem[*logmodel.ConversionLog]{filePath: fileConfig.LogsFilePath}

	if err := setupPersistence(requestsItem); err != nil {
		log.Fatalf("Failed to setup persistence for requests: %v", err)
	}

	if err := setupPersistence(responsesItem); err != nil {
		log.Fatalf("Failed to setup persistence for responses: %v", err)
	}

	if err := setupPersistence(logsItem); err != nil {
		log.Fatalf("Failed to setup persistence for logs: %v", err)
	}
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

	jsonData, err := json.Marshal(item)

	if err != nil {
		fmt.Printf("Error marshaling to JSON for file %s: %v\n", ri.filePath, err)

		return
	}

	if _, err := ri.file.Write(jsonData); err != nil {
		fmt.Printf("Error writing data to file %s: %v\n", ri.filePath, err)

		return
	}

	if _, err := ri.file.WriteString("\n"); err != nil {
		fmt.Printf("Error writing newline to file %s: %v\n", ri.filePath, err)

		return
	}
}

func GetNewConversionRequests() []*api.Request {
	return requestsItem.getNew()
}

func GetNewConversionResponses() []*api.Response {
	return responsesItem.getNew()
}

func GetNewConversionLogs() []*logmodel.ConversionLog {
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

func AddLog(convLog *logmodel.ConversionLog) {
	logsItem.add(convLog)
	fmt.Printf("Added ConversionLog: GetId=%s, GetTimestamp=%s\n", convLog.Id, convLog.Timestamp.Format(time.RFC3339))
}

func setupPersistence[T any](item *repositoryItem[T]) error {
	if err := os.MkdirAll(filepath.Dir(item.filePath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	if err := loadDataFromFile(item); err != nil {
		return fmt.Errorf("failed to load data from file %s: %w", item.filePath, err)
	}

	file, err := os.OpenFile(item.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return fmt.Errorf("failed to open file %s for appending: %w", item.filePath, err)
	}

	item.file = file
	fmt.Printf("Persistence setup complete for %s. Current items in memory: %d\n", item.filePath, len(item.data))

	return nil
}

func loadDataFromFile[T any](item *repositoryItem[T]) error {
	file, err := os.OpenFile(item.filePath, os.O_RDONLY, 0644)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("failed to open file %s for reading: %w", item.filePath, err)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()
		var record T

		if err := json.Unmarshal(line, &record); err != nil {
			fmt.Printf("Error unmarshaling line from %s: %v, line: %s\n", item.filePath, err, string(line))

			continue
		}

		item.data = append(item.data, record)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", item.filePath, err)
	}

	item.lastRead = len(item.data)

	return nil
}

func ClosePersistence() {
	if requestsItem.file != nil {
		_ = requestsItem.file.Close()
		fmt.Println("Closed requests persistence file.")
	}

	if responsesItem.file != nil {
		_ = responsesItem.file.Close()
		fmt.Println("Closed responses persistence file.")
	}

	if logsItem.file != nil {
		_ = logsItem.file.Close()
		fmt.Println("Closed logs persistence file.")
	}
}
