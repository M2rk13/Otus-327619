package webserver

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"
	"github.com/M2rk13/Otus-327619/internal/repository"

	"github.com/gin-gonic/gin"
)

type crudService[T any] interface {
	Create(T)
	GetByID(string) T
	GetAll() []T
	Update(T) bool
	Delete(string) bool
}

func createHandler[T any](repo crudService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newItem T

		if err := c.ShouldBindJSON(&newItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})

			return
		}

		repo.Create(newItem)
		c.JSON(http.StatusCreated, newItem)
	}
}

func getByIDHandler[T any](repo crudService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		item := repo.GetByID(id)

		if reflect.ValueOf(item).IsNil() {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})

			return
		}

		c.JSON(http.StatusOK, item)
	}
}

func getAllHandler[T any](repo crudService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		items := repo.GetAll()
		c.JSON(http.StatusOK, items)
	}
}

func updateHandler[T any](repo crudService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var updatedItem T

		if err := c.ShouldBindJSON(&updatedItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})

			return
		}

		v := reflect.ValueOf(updatedItem).Elem().FieldByName("Id")

		if v.IsValid() && v.String() != id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Id in body must match Id in path"})

			return
		}

		if repo.Update(updatedItem) {
			c.JSON(http.StatusOK, updatedItem)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		}
	}
}

func deleteHandler[T any](repo crudService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if repo.Delete(id) {
			c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		}
	}
}

type RequestService struct{}

func (s *RequestService) Create(req *api.Request)        { repository.CreateRequest(req) }
func (s *RequestService) GetByID(id string) *api.Request { return repository.GetRequestByID(id) }
func (s *RequestService) GetAll() []*api.Request         { return repository.GetAllRequests() }
func (s *RequestService) Update(req *api.Request) bool   { return repository.UpdateRequest(req) }
func (s *RequestService) Delete(id string) bool          { return repository.DeleteRequest(id) }

type ResponseService struct{}

func (s *ResponseService) Create(resp *api.Response)       { repository.CreateResponse(resp) }
func (s *ResponseService) GetByID(id string) *api.Response { return repository.GetResponseByID(id) }
func (s *ResponseService) GetAll() []*api.Response         { return repository.GetAllResponses() }
func (s *ResponseService) Update(resp *api.Response) bool  { return repository.UpdateResponse(resp) }
func (s *ResponseService) Delete(id string) bool           { return repository.DeleteResponse(id) }

type LogService struct{}

func (s *LogService) Create(log *log.ConversionLog) { repository.CreateConversionLog(log) }
func (s *LogService) GetByID(id string) *log.ConversionLog {
	return repository.GetConversionLogByID(id)
}
func (s *LogService) GetAll() []*log.ConversionLog       { return repository.GetAllConversionLogs() }
func (s *LogService) Update(log *log.ConversionLog) bool { return repository.UpdateConversionLog(log) }
func (s *LogService) Delete(id string) bool              { return repository.DeleteConversionLog(id) }

var (
	requestService  = &RequestService{}
	responseService = &ResponseService{}
	logService      = &LogService{}
)
