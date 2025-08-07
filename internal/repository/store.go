package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/M2rk13/Otus-327619/internal/config"
	"github.com/M2rk13/Otus-327619/internal/model/api"
	logmodel "github.com/M2rk13/Otus-327619/internal/model/log"

	"github.com/google/uuid"
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

type Identifiable interface {
	GetId() string
}

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

func setupPersistence[T any](repoItem *repositoryItem[T]) error {
	if err := os.MkdirAll(filepath.Dir(repoItem.filePath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	if err := loadDataFromFile(repoItem); err != nil {
		return fmt.Errorf("failed to load data from file %s: %w", repoItem.filePath, err)
	}

	file, err := os.OpenFile(repoItem.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return fmt.Errorf("failed to open file %s for appending: %w", repoItem.filePath, err)
	}

	repoItem.file = file
	fmt.Printf("Persistence setup complete for %s. Current items in memory: %d\n", repoItem.filePath, len(repoItem.data))

	return nil
}

func loadDataFromFile[T any](repoItem *repositoryItem[T]) error {
	file, err := os.OpenFile(repoItem.filePath, os.O_RDONLY, 0644)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("failed to open file %s for reading: %w", repoItem.filePath, err)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()
		var record T

		if err := json.Unmarshal(line, &record); err != nil {
			fmt.Printf("Error unmarshaling line from %s: %v, line: %s\n", repoItem.filePath, err, string(line))

			continue
		}

		repoItem.data = append(repoItem.data, record)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", repoItem.filePath, err)
	}

	repoItem.lastRead = len(repoItem.data)

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

func (ri *repositoryItem[T]) rewriteAllDataToFile() error {
	if ri.file != nil {
		if err := ri.file.Close(); err != nil {
			return fmt.Errorf("failed to close file before rewrite: %w", err)
		}
	}

	f, err := os.OpenFile(ri.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		return fmt.Errorf("failed to open file %s for rewriting: %w", ri.filePath, err)
	}

	defer func() {
		reopenFile, reErr := os.OpenFile(ri.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if reErr != nil {
			fmt.Printf("CRITICAL: Failed to reopen file %s for appending after rewrite: %v\n", ri.filePath, reErr)
			ri.file = nil
		} else {
			ri.file = reopenFile
		}
	}()

	writer := bufio.NewWriter(f)

	for _, data := range ri.data {
		jsonData, err := json.Marshal(data)

		if err != nil {
			return fmt.Errorf("failed to marshal data to JSON during rewrite: %w", err)
		}

		if _, err := writer.Write(jsonData); err != nil {
			return fmt.Errorf("failed to write JSON data during rewrite: %w", err)
		}

		if _, err := writer.WriteString("\n"); err != nil {
			return fmt.Errorf("failed to write newline during rewrite: %w", err)
		}
	}

	return writer.Flush()
}

func genericGetNew[T any](repoItem *repositoryItem[T]) []T {
	return repoItem.getNew()
}

func genericGetAll[T any](repoItem *repositoryItem[T]) []T {
	repoItem.mu.Lock()
	defer repoItem.mu.Unlock()

	return append([]T{}, repoItem.data...)
}

func genericAdd[T any](repoItem *repositoryItem[T], item T) {
	repoItem.add(item)
}

func genericGetByID[T Identifiable](repoItem *repositoryItem[T], id string) T {
	repoItem.mu.Lock()
	defer repoItem.mu.Unlock()

	for _, data := range repoItem.data {
		if data.GetId() == id {
			return data
		}
	}

	var zero T

	return zero
}

func genericUpdate[T Identifiable](repoItem *repositoryItem[T], updatedData T) bool {
	repoItem.mu.Lock()
	defer repoItem.mu.Unlock()

	for i, data := range repoItem.data {
		if data.GetId() == updatedData.GetId() {
			repoItem.data[i] = updatedData
			_ = repoItem.rewriteAllDataToFile()

			return true
		}
	}

	return false
}

func genericDelete[T Identifiable](repoItem *repositoryItem[T], id string) bool {
	repoItem.mu.Lock()
	defer repoItem.mu.Unlock()

	initialLen := len(repoItem.data)
	var newData []T

	for _, data := range repoItem.data {
		if data.GetId() != id {
			newData = append(newData, data)
		}
	}

	repoItem.data = newData

	if len(repoItem.data) < initialLen {
		_ = repoItem.rewriteAllDataToFile()

		return true
	}

	return false
}

func GetNewConversionRequests() []*api.Request {
	return genericGetNew(requestsItem)
}

func GetNewConversionResponses() []*api.Response {
	return genericGetNew(responsesItem)
}

func GetNewConversionLogs() []*logmodel.ConversionLog {
	return genericGetNew(logsItem)
}

func CreateRequest(req *api.Request) {
	req.Id = uuid.New().String()
	genericAdd(requestsItem, req)
}

func GetRequestByID(id string) *api.Request {
	return genericGetByID(requestsItem, id)
}

func GetAllRequests() []*api.Request {
	return genericGetAll(requestsItem)
}

func UpdateRequest(req *api.Request) bool {
	return genericUpdate(requestsItem, req)
}

func DeleteRequest(id string) bool {
	return genericDelete(requestsItem, id)
}

func CreateResponse(resp *api.Response) {
	resp.Id = uuid.New().String()
	genericAdd(responsesItem, resp)
}

func GetResponseByID(id string) *api.Response {
	return genericGetByID(responsesItem, id)
}

func GetAllResponses() []*api.Response {
	return genericGetAll(responsesItem)
}

func UpdateResponse(resp *api.Response) bool {
	return genericUpdate(responsesItem, resp)
}

func DeleteResponse(id string) bool {
	return genericDelete(responsesItem, id)
}

func CreateConversionLog(log *logmodel.ConversionLog) {
	log.Id = uuid.New().String()
	genericAdd(logsItem, log)
}

func GetConversionLogByID(id string) *logmodel.ConversionLog {
	return genericGetByID(logsItem, id)
}

func GetAllConversionLogs() []*logmodel.ConversionLog {
	return genericGetAll(logsItem)
}

func UpdateConversionLog(log *logmodel.ConversionLog) bool {
	return genericUpdate(logsItem, log)
}

func DeleteConversionLog(id string) bool {
	return genericDelete(logsItem, id)
}
